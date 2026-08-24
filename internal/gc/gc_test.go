package gc

import (
	"testing"
	"time"

	"columnstore/internal/col"
)

func TestRunReclaimsOldBlocksOnly(t *testing.T) {
	store := col.NewStore(t.TempDir())
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	oldBlock := col.NewBlock("t", "c", []col.Value{col.Int(1)}, now.Add(-48*time.Hour))
	freshBlock := col.NewBlock("t", "c", []col.Value{col.Int(2)}, now)
	if _, err := store.AddBlock(oldBlock); err != nil {
		t.Fatal(err)
	}
	if _, err := store.AddBlock(freshBlock); err != nil {
		t.Fatal(err)
	}
	reclaimer := New(store, DefaultPolicy().WithTTL(24*time.Hour))
	report, err := reclaimer.Run(now)
	if err != nil {
		t.Fatal(err)
	}
	if report.Reclaimed != 1 || report.Fresh != 1 {
		t.Fatalf("expected one reclaimed and one fresh, got %+v", report)
	}
	if store.BlockCount() != 1 {
		t.Fatal("fresh block must survive reclamation")
	}
}

func TestSchedulerRespectsInterval(t *testing.T) {
	store := col.NewStore(t.TempDir())
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	reclaimer := New(store, DefaultPolicy())
	scheduler := NewScheduler(10*time.Minute, reclaimer.Run)
	first, err := scheduler.Tick(now)
	if err != nil {
		t.Fatal(err)
	}
	second, err := scheduler.Tick(now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if second != first {
		t.Fatal("tick inside the interval must reuse the previous report")
	}
}

func TestPolicyOverride(t *testing.T) {
	policy := DefaultPolicy().WithTTL(time.Hour)
	if policy.TTL != time.Hour {
		t.Fatal("policy override failed")
	}
}
