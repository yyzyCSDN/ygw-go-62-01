package main

import (
	"fmt"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/load"
)

// seedDemo loads a small analytics table so the console has data to query.
func seedDemo(loader *load.Loader) error {
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	batches := []load.Batch{
		{
			Table:  "events",
			Column: "latency_ms",
			Rows: []col.Value{
				col.Int(12), col.Int(45), col.Int(31), col.Int(28),
			},
			Start: now,
			End:   now.Add(time.Hour),
			Seal:  true,
		},
		{
			Table:  "events",
			Column: "latency_ms",
			Rows: []col.Value{
				col.Int(9), col.Int(77), col.Int(63), col.Int(51), col.Int(40),
			},
			Start: now.Add(time.Hour),
			End:   now.Add(2 * time.Hour),
			Seal:  true,
		},
	}
	_, err := loader.Load("events", load.NewSliceSource(batches))
	if err != nil {
		return fmt.Errorf("seed demo data: %w", err)
	}
	return nil
}
