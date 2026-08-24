package part

import "time"

// Overlaps reports whether a partition intersects the inclusive query window
// [from, to]. A partition is outside the window only when its entire interval
// lies before from or strictly after to.
func Overlaps(p *Partition, from, to time.Time) bool {
	if p == nil {
		return false
	}
	return !p.End.Before(from) && !p.Start.After(to)
}

// ClampRange narrows a query window to the part of the partition that is in
// range. The returned bounds are half-open [lo, hi).
func ClampRange(p *Partition, from, to time.Time) (time.Time, time.Time) {
	lo := from
	if p.Start.After(lo) {
		lo = p.Start
	}
	hi := to.Add(time.Nanosecond)
	if p.End.Before(hi) {
		hi = p.End
	}
	return lo, hi
}
