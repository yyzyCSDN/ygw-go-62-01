package part

import "time"

// keepPartition reports whether p overlaps the inclusive query window
// [from, to]. A partition is outside the window only when its entire interval
// lies before from or strictly after to; any partition touching the boundary
// is retained so boundary data is never lost.
func keepPartition(p *Partition, from, to time.Time) bool {
	if p == nil {
		return false
	}
	if p.End.Before(from) {
		return false
	}
	if p.Start.After(to) {
		return false
	}
	return true
}

// Prune returns the partitions of a table that overlap [from, to].
func Prune(parts []*Partition, from, to time.Time) []*Partition {
	kept := make([]*Partition, 0, len(parts))
	for _, p := range parts {
		if !keepPartition(p, from, to) {
			continue
		}
		kept = append(kept, p)
	}
	return kept
}
