package meta

import (
	"sort"
	"sync"
	"time"
)

// Registry stores the latest metadata snapshot per table.
type Registry struct {
	mu    sync.RWMutex
	byTbl map[string]*ColumnMeta
}

// NewRegistry builds an empty metadata registry.
func NewRegistry() *Registry {
	return &Registry{byTbl: make(map[string]*ColumnMeta)}
}

// Update registers a new schema revision and bumps the table version.
func (r *Registry) Update(schema Schema, now time.Time) (*ColumnMeta, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing := r.byTbl[schema.Table]
	version := Version{Table: schema.Table, Number: 1}
	if existing != nil {
		version = existing.SchemaVersion.Next()
	}
	meta := &ColumnMeta{
		Table:         schema.Table,
		Schema:        schema,
		SchemaVersion: version,
		UpdatedAt:     now,
	}
	r.byTbl[schema.Table] = meta
	return meta, nil
}

// Lookup returns the latest metadata for a table.
func (r *Registry) Lookup(table string) (*ColumnMeta, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	meta, ok := r.byTbl[table]
	if !ok {
		return nil, ErrUnknownTable
	}
	return meta, nil
}

// LookupAtVersion returns the metadata if it still matches a version number.
func (r *Registry) LookupAtVersion(table string, number uint64) (*ColumnMeta, error) {
	meta, err := r.Lookup(table)
	if err != nil {
		return nil, err
	}
	if !meta.AtVersion(number) {
		return nil, ErrUnknownTable
	}
	return meta, nil
}

// MarkCompressedAtVersion records that a table's blocks are dictionary
// compressed, but only when the metadata still matches the expected version.
func (r *Registry) MarkCompressedAtVersion(table string, number uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	meta, ok := r.byTbl[table]
	if !ok {
		return ErrUnknownTable
	}
	if !meta.AtVersion(number) {
		return ErrUnknownTable
	}
	meta.Compressed = true
	return nil
}


// Tables returns all registered table names.
func (r *Registry) Tables() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.byTbl))
	for name := range r.byTbl {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
