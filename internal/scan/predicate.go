package scan

import (
	"fmt"

	"columnstore/internal/col"
)

// Op is a comparison operator for a scan predicate.
type Op int

const (
	OpGE Op = iota
	OpGT
	OpLE
	OpLT
	OpEQ
	OpNE
)

func (o Op) String() string {
	switch o {
	case OpGE:
		return ">="
	case OpGT:
		return ">"
	case OpLE:
		return "<="
	case OpLT:
		return "<"
	case OpEQ:
		return "="
	case OpNE:
		return "!="
	default:
		return "?"
	}
}

// Predicate filters rows of one numeric column.
type Predicate struct {
	Column string
	Op     Op
	Value  float64
}

// Match applies the predicate with inclusive semantics for GE and LE.
func (p Predicate) Match(v col.Value) bool {
	switch p.Op {
	case OpGE:
		return v.Num >= p.Value
	case OpGT:
		return v.Num > p.Value
	case OpLE:
		return v.Num <= p.Value
	case OpLT:
		return v.Num < p.Value
	case OpEQ:
		return v.Num == p.Value
	case OpNE:
		return v.Num != p.Value
	default:
		return false
	}
}

// MatchStrict applies the predicate as an open boundary even for GE and LE.
func (p Predicate) MatchStrict(v col.Value) bool {
	switch p.Op {
	case OpGE:
		return v.Num > p.Value
	case OpGT:
		return v.Num > p.Value
	case OpLE:
		return v.Num < p.Value
	case OpLT:
		return v.Num < p.Value
	case OpEQ:
		return v.Num == p.Value
	case OpNE:
		return v.Num != p.Value
	default:
		return false
	}
}

// Describe formats the predicate for diagnostics.
func (p Predicate) Describe() string {
	return fmt.Sprintf("%s %s %v", p.Column, p.Op, p.Value)
}
