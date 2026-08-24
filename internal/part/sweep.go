package part

import "time"

// SweepArchived seals active partitions whose window ended before cutoff and
// archives sealed partitions that have not been touched since the cutoff.
func SweepArchived(registry *Registry, table string, cutoff time.Time) (sealed, archived int) {
	for _, p := range registry.Table(table) {
		switch p.State {
		case StateActive:
			if p.End.Before(cutoff) {
				p.Seal()
				sealed++
			}
		case StateSealed:
			if p.End.Before(cutoff.Add(-24 * time.Hour)) {
				p.Archive()
				archived++
			}
		}
	}
	return sealed, archived
}
