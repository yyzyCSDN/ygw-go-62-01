package main

import (
	"encoding/json"
	"net/http"

	"columnstore/internal/col"
	"columnstore/internal/compress"
	"columnstore/internal/gc"
	"columnstore/internal/load"
	"columnstore/internal/meta"
	"columnstore/internal/part"
	"columnstore/internal/query"
	"columnstore/internal/scan"
)

// App bundles every component behind one HTTP surface.
type App struct {
	store    *col.Store
	parts    *part.Registry
	meta     *meta.Registry
	loader   *load.Loader
	engine   *query.Engine
	gc       *gc.GC
	compress compressService
	dir      string
}

type compressService interface {
	CompressColumn(block *col.ColumnBlock) error
	Rebuild(block *col.ColumnBlock, rows []col.Value) (*compress.Dictionary, error)
}

// NewApp wires the server components together.
func NewApp(store *col.Store, parts *part.Registry, meta *meta.Registry,
	loader *load.Loader, engine *query.Engine, reclaimer *gc.GC, compressor compressService, dir string) *App {
	return &App{
		store:    store,
		parts:    parts,
		meta:     meta,
		loader:   loader,
		engine:   engine,
		gc:       reclaimer,
		compress: compressor,
		dir:      dir,
	}
}

// healthPayload is the /api/health response body.
type healthPayload struct {
	Status     string         `json:"status"`
	Blocks     int            `json:"blocks"`
	Partitions int            `json:"partitions"`
	Handles    int            `json:"handles"`
	Tables     int            `json:"tables"`
	PerTable   map[string]int `json:"per_table"`
}

func (a *App) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, healthPayload{
		Status:     "ok",
		Blocks:     a.store.BlockCount(),
		Partitions: a.parts.Count(),
		Handles:    a.store.OpenHandles(),
		Tables:     len(a.store.Tables()),
		PerTable:   a.loader.Manifest().Snapshot(),
	})
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(value)
}

func (a *App) partitions(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"partitions": part.DescribeAll(a.parts.Table("events"))})
}

func (a *App) blocks(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"blocks": a.store.ListBlocks()})
}

func (a *App) scanStats(w http.ResponseWriter, _ *http.Request) {
	out := make([]scan.BlockStats, 0)
	for _, block := range a.store.ListBlocks() {
		out = append(out, scan.Stats(block))
	}
	writeJSON(w, map[string]any{"stats": out})
}
