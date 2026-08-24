package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/gc"
	"columnstore/internal/load"
	"columnstore/internal/part"
	"columnstore/internal/query"
	"columnstore/internal/scan"
)

// loadRequest is the body accepted by POST /api/load.
type loadRequest struct {
	Table   string      `json:"table"`
	Batches []batchBody `json:"batches"`
}

type batchBody struct {
	Column string        `json:"column"`
	Values []json.Number `json:"values"`
	Start  string        `json:"start"`
	End    string        `json:"end"`
	Seal   bool          `json:"seal"`
}

func (a *App) handleLoad(w http.ResponseWriter, r *http.Request) {
	var req loadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Table == "" {
		http.Error(w, "table is required", http.StatusBadRequest)
		return
	}
	batches := make([]load.Batch, 0, len(req.Batches))
	for _, body := range req.Batches {
		start, err := time.Parse(time.RFC3339, body.Start)
		if err != nil {
			http.Error(w, "start must be RFC3339", http.StatusBadRequest)
			return
		}
		end, err := time.Parse(time.RFC3339, body.End)
		if err != nil {
			http.Error(w, "end must be RFC3339", http.StatusBadRequest)
			return
		}
		rows := make([]col.Value, 0, len(body.Values))
		for _, raw := range body.Values {
			num, _ := raw.Float64()
			rows = append(rows, col.Float(num))
		}
		batches = append(batches, load.Batch{
			Table:  req.Table,
			Column: body.Column,
			Rows:   rows,
			Start:  start,
			End:    end,
			Seal:   body.Seal,
		})
	}
	written, err := a.loader.Load(req.Table, load.NewSliceSource(batches))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	writeJSON(w, map[string]any{
		"written": written,
		"blocks":  a.store.BlockCount(),
		"table":   a.loader.Manifest().TableBlocks(req.Table),
	})
}

func (a *App) handleQuery(w http.ResponseWriter, r *http.Request) {
	table := r.URL.Query().Get("table")
	fromText := r.URL.Query().Get("from")
	toText := r.URL.Query().Get("to")
	column := r.URL.Query().Get("column")
	if table == "" {
		table = "events"
	}
	if column == "" {
		column = "latency_ms"
	}
	from, to, err := query.ParseWindow(fromText, toText)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	minText := r.URL.Query().Get("min")
	var pred *scan.Predicate
	if minText != "" {
		pred, err = query.BuildPredicate(column, "ge", minText)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	plan, err := a.engine.Build(table, from, to, []string{column}, pred)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	result, err := a.engine.Execute(plan)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{
		"columns": result.Columns,
		"text":    result.Text(),
		"scanned": result.Scanned,
		"rows":    result.Rows,
	})
}

func (a *App) handleCompress(w http.ResponseWriter, _ *http.Request) {
	compressed := 0
	for _, block := range a.store.ListBlocks() {
		if block.Column != "region" {
			continue
		}
		if _, err := a.compress.Rebuild(block, block.Rows); err != nil {
			http.Error(w, fmt.Sprintf("dict rebuild failed: %v", err), http.StatusInternalServerError)
			return
		}
		if err := a.compress.CompressColumn(block); err != nil {
			http.Error(w, fmt.Sprintf("compress failed: %v", err), http.StatusInternalServerError)
			return
		}
		compressed++
	}
	writeJSON(w, map[string]any{"compressed": compressed})
}

func (a *App) handleGC(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	report, err := a.gc.Run(now)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sealed, archived := part.SweepArchived(a.parts, "events", now)
	diskBytes, diskErr := gc.DiskUsage(a.dir)
	if diskErr != nil {
		diskBytes = 0
	}
	writeJSON(w, map[string]any{
		"scanned":    report.Scanned,
		"reclaimed":  report.Reclaimed,
		"fresh":      report.Fresh,
		"sealed":     sealed,
		"archived":   archived,
		"disk_bytes": diskBytes,
	})
}

func (a *App) handleVerify(w http.ResponseWriter, _ *http.Request) {
	verified, err := a.store.VerifyAllSegments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"verified": verified})
}
