package part

import "time"

// PartState is the lifecycle state of a partition.
type PartState int

const (
	StateActive PartState = iota
	StateSealed
	StateArchived
)

func (s PartState) String() string {
	switch s {
	case StateActive:
		return "active"
	case StateSealed:
		return "sealed"
	case StateArchived:
		return "archived"
	default:
		return "unknown"
	}
}

// Partition groups the block ids that fall inside one time interval.
type Partition struct {
	ID       uint64
	Table    string
	Start    time.Time
	End      time.Time
	State    PartState
	BlockIDs []uint64
}

// NewPartition creates an active partition over [start, end).
func NewPartition(table string, start, end time.Time, id uint64) *Partition {
	return &Partition{
		ID:    id,
		Table: table,
		Start: start,
		End:   end,
		State: StateActive,
	}
}

// Seal freezes the partition so no further blocks may join it.
func (p *Partition) Seal() {
	p.State = StateSealed
}

// Archive marks the partition as no longer queryable.
func (p *Partition) Archive() {
	p.State = StateArchived
}

// AddBlock appends a block id when the partition is still active.
func (p *Partition) AddBlock(id uint64) bool {
	if p.State != StateActive {
		return false
	}
	p.BlockIDs = append(p.BlockIDs, id)
	return true
}

// Contains reports whether the partition covers the given instant.
func (p *Partition) Contains(instant time.Time) bool {
	return !instant.Before(p.Start) && instant.Before(p.End)
}
