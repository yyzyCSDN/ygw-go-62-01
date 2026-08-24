package load

import (
	"fmt"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/meta"
	"columnstore/internal/part"
)

// Loader persists batches into blocks and partitions atomically.
type Loader struct {
	store    *col.Store
	parts    *part.Registry
	meta     *meta.Registry
	manifest *Manifest
	now      func() time.Time
}

// NewLoader wires a loader over the shared stores.
func NewLoader(store *col.Store, parts *part.Registry, meta *meta.Registry) *Loader {
	return &Loader{
		store:    store,
		parts:    parts,
		meta:     meta,
		manifest: NewManifest(),
		now:      time.Now,
	}
}

// Manifest returns the load manifest of the loader.
func (l *Loader) Manifest() *Manifest {
	return l.manifest
}

// Load writes every batch and registers matching partitions. When any batch
// fails validation the whole operation rolls back.
func (l *Loader) Load(table string, src Source) (int, error) {
	batches := src.Batches()
	plan := BuildPlan(table, src)
	if _, err := l.meta.Update(plan.Schema, l.now()); err != nil {
		return 0, err
	}
	written := make([]uint64, 0, len(batches))
	partitions := make([]uint64, 0, len(batches))
	flush := func(b Batch) error {
		if err := Validate(b); err != nil {
			return err
		}
		block := col.NewBlock(b.Table, b.Column, b.Rows, l.now())
		id, err := l.store.AddBlock(block)
		if err != nil {
			return err
		}
		written = append(written, id)
		if err := block.Seal(block.Segment); err != nil {
			return err
		}
		if err := l.store.WriteBlock(block); err != nil {
			return err
		}
		if !b.Seal {
			return nil
		}
		p := PartitionFor(b.Table, b, 0)
		pid, err := l.parts.Register(p)
		if err != nil {
			return err
		}
		p.AddBlock(id)
		partitions = append(partitions, pid)
		return nil
	}
	for _, b := range batches {
		if err := flush(b); err != nil {
			l.manifest.Rollback(table, len(written))
			_ = rollback(l.store, l.parts, written, partitions)
			return 0, fmt.Errorf("load %s failed: %w", b.Describe(), err)
		}
	}
	l.manifest.Add(table, len(batches))
	return len(batches), nil
}

// rollback undoes every block and partition written by a failed load.
// In-flight blocks (written but not yet sealed into a partition) and already
// sealed partitions must both be removed so no residue remains visible.
// Every block id in written is removed from the store, which also drops the
// backing segment file, covering both the in-flight batches that never reached
// a partition and the sealed batches whose partition is removed below.
func rollback(store *col.Store, parts *part.Registry, written []uint64, partitions []uint64) error {
	for _, pid := range partitions {
		_ = parts.Remove(pid)
	}
	for _, id := range written {
		_ = store.RemoveBlock(id)
	}
	return nil
}
