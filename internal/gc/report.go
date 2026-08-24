package gc

import "time"

// Report summarizes one reclamation pass.
type Report struct {
	StartedAt  time.Time
	FinishedAt time.Time
	Scanned    int
	Reclaimed  int
	Fresh      int
}

// Finish stamps the end time and returns the report.
func (r *Report) Finish(now time.Time) *Report {
	r.FinishedAt = now
	return r
}

// Add merges a scan result into the report counters.
func (r *Report) Add(scanned, reclaimed, fresh int) {
	r.Scanned += scanned
	r.Reclaimed += reclaimed
	r.Fresh += fresh
}
