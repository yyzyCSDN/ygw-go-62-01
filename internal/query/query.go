package query

import (
	"fmt"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/compress"
	"columnstore/internal/meta"
	"columnstore/internal/part"
	"columnstore/internal/scan"
)

// Engine runs analytical queries against the column store.
type Engine struct {
	store   *col.Store
	parts   *part.Registry
	meta    *meta.Registry
	catalog *meta.Catalog
	dicts   *compress.DictRegistry
	chunk   int
}

// NewEngine wires the query engine over the shared stores.
func NewEngine(store *col.Store, parts *part.Registry, metaRegistry *meta.Registry, dicts *compress.DictRegistry) *Engine {
	return &Engine{
		store:   store,
		parts:   parts,
		meta:    metaRegistry,
		catalog: meta.NewCatalog(metaRegistry),
		dicts:   dicts,
		chunk:   scan.ChunkSize,
	}
}

// SetChunk overrides the default scan chunk size.
func (e *Engine) SetChunk(n int) {
	if n > 0 {
		e.chunk = n
	}
}

// Build prepares a plan from the latest metadata snapshot.
func (e *Engine) Build(table string, from, to time.Time, columns []string, pred *scan.Predicate) (*Plan, error) {
	metaValue, err := e.catalog.MetaFor(table)
	if err != nil {
		return nil, err
	}
	for _, column := range columns {
		if _, err := metaValue.Column(column); err != nil {
			return nil, err
		}
	}
	return &Plan{
		Table:   table,
		From:    from,
		To:      to,
		Meta:    metaValue,
		Version: metaValue.SchemaVersion.Number,
		Columns: append([]string(nil), columns...),
		Pred:    pred,
	}, nil
}

// Execute runs a prepared plan and returns the projected result.
func (e *Engine) Execute(plan *Plan) (*Result, error) {
	if plan == nil {
		return nil, fmt.Errorf("nil query plan")
	}
	if err := e.meta.VerifyVersion(plan.Table, plan.Version); err != nil {
		return nil, err
	}
	parts := part.Prune(e.parts.Table(plan.Table), plan.From, plan.To)
	for _, p := range parts {
		if !part.Overlaps(p, plan.From, plan.To) {
			return nil, fmt.Errorf("%s: pruned partition %d falls outside the query window", plan.RangeText(), p.ID)
		}
	}
	exec := newExecutor(e.chunk)
	result := Empty(plan.Columns)
	aggregate := Aggregate{Column: first(plan.Columns)}
	scanned := 0
	for _, p := range parts {
		for _, blockID := range e.blocksFor(plan, p) {
			block, err := e.store.GetBlock(blockID)
			if err != nil {
				return nil, err
			}
			consume := func(batch []col.Value) {
				acc := &scan.Accumulator{}
				acc.AddRows(batch)
				aggregate.FoldAccumulator(acc)
			}
			count := 0
			if plan.Pred != nil {
				count = exec.scanBlockFiltered(block, *plan.Pred, plan.Strict, consume)
			} else if block.State == col.StateCompressed {
				rows, err := e.materialize(block)
				if err != nil {
					return nil, err
				}
				count = exec.scanRows(rows, consume)
			} else {
				count = scan.ScanBlock(block, e.chunk, consume)
			}
			scanned += count
		}
	}
	result.Scanned = scanned
	if aggregate.Count > 0 {
		result.Append([]col.Value{col.Int(aggregate.Count), col.Float(aggregate.Sum)})
	}
	return result, nil
}

// blocksFor returns the partition blocks that carry the requested column.
func (e *Engine) blocksFor(plan *Plan, p *part.Partition) []uint64 {
	target := first(plan.Columns)
	out := make([]uint64, 0, 1)
	for _, id := range p.BlockIDs {
		block, err := e.store.GetBlock(id)
		if err != nil {
			continue
		}
		if block.Column == target {
			out = append(out, id)
		}
	}
	return out
}

// materialize loads block rows, decoding compressed dictionaries.
func (e *Engine) materialize(block *col.ColumnBlock) ([]col.Value, error) {
	if block.State == col.StateCompressed {
		reader := scan.OpenReader(e.dicts, block)
		return reader.Rows()
	}
	return col.ReadColumn(block), nil
}

func first(names []string) string {
	if len(names) == 0 {
		return ""
	}
	return names[0]
}
