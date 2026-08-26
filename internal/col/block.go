package col

import (
	"time"
)

// Kind describes the physical type of a stored value.
type Kind int

const (
	KindInt Kind = iota
	KindFloat
	KindString
)

// Value is one cell inside a column block.
type Value struct {
	Kind Kind
	Num  float64
	Str  string
}

// Int builds an integer value.
func Int(v int64) Value {
	return Value{Kind: KindInt, Num: float64(v)}
}

// Float builds a floating point value.
func Float(v float64) Value {
	return Value{Kind: KindFloat, Num: v}
}

// Str builds a string value.
func Str(v string) Value {
	return Value{Kind: KindString, Str: v}
}

// State is the lifecycle state of a column block.
type State int

const (
	StateMutable State = iota
	StateImmutable
	StateCompressed
	StateReclaimed
)

func (s State) String() string {
	switch s {
	case StateMutable:
		return "mutable"
	case StateImmutable:
		return "immutable"
	case StateCompressed:
		return "compressed"
	case StateReclaimed:
		return "reclaimed"
	default:
		return "unknown"
	}
}

// ColumnBlock stores one immutable batch of values for a single column.
type ColumnBlock struct {
	ID        uint64
	Table     string
	Column    string
	Rows      []Value
	Payload   []byte
	Checksum  uint64
	State     State
	DictID    uint64
	CreatedAt time.Time
	Segment   string
}

// NewBlock creates a mutable block carrying a creation timestamp.
func NewBlock(table, column string, rows []Value, now time.Time) *ColumnBlock {
	return &ColumnBlock{
		Table:     table,
		Column:    column,
		Rows:      append([]Value(nil), rows...),
		State:     StateMutable,
		CreatedAt: now,
	}
}

// Seal moves the block into the immutable state and records the segment name.
func (b *ColumnBlock) Seal(segment string) error {
	if b.State != StateMutable {
		return ErrBlockSealed
	}
	b.State = StateImmutable
	b.Segment = segment
	return nil
}

// Len returns the number of rows in the block.
func (b *ColumnBlock) Len() int {
	if b == nil {
		return 0
	}
	return len(b.Rows)
}

// At returns the value at row index i.
func (b *ColumnBlock) At(i int) Value {
	return b.Rows[i]
}
