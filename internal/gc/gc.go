package gc

import (
	"time"

	"columnstore/internal/col"
)

// GC reclaims column blocks that outlived the retention policy.
type GC struct {
	store  *col.Store
	policy Policy
}

// New creates a reclaimer over the block store.
func New(store *col.Store, policy Policy) *GC {
	return &GC{store: store, policy: policy}
}

// Run performs one reclamation pass at the given instant.
func (g *GC) Run(now time.Time) (*Report, error) {
	report := &Report{StartedAt: now}
	scanned := 0
	reclaimed := 0
	fresh := 0
	for _, block := range g.store.ListBlocks() {
		scanned++
		if now.Sub(block.CreatedAt) >= g.policy.TTL {
			if err := g.store.ReclaimBlock(block); err != nil {
				return nil, err
			}
			reclaimed++
			continue
		}
		fresh++
	}
	report.Add(scanned, reclaimed, fresh)
	return report.Finish(now), nil
}
