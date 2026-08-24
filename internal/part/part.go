package part

import "time"

// Prune returns the partitions of a table that overlap [from, to].
func Prune(parts []*Partition, from, to time.Time) []*Partition {
	kept := make([]*Partition, 0, len(parts))
	for _, p := range parts {
		if p.End.Before(from) || !p.Start.Before(to) {
			continue
		}
		kept = append(kept, p)
	}
	return kept
}
