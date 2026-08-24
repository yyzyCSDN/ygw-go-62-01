package part_test

import (
	"testing"
	"time"

	"columnstore/internal/part"
)

func TestPartitionPruneKeepsBoundaryData(t *testing.T) {
	base := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	inside := part.NewPartition("events", base, base.Add(time.Hour), 1)
	boundary := part.NewPartition("events", base.Add(time.Hour), base.Add(2*time.Hour), 2)
	kept := part.Prune([]*part.Partition{inside, boundary}, base, base.Add(time.Hour))
	if len(kept) != 2 {
		t.Fatalf("boundary partition must be kept, got %d partitions", len(kept))
	}
	found := false
	for _, p := range kept {
		if p.ID == 2 {
			found = true
		}
	}
	if !found {
		t.Fatal("partition starting exactly at the query end was pruned")
	}
}
