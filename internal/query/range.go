package query

import (
	"fmt"
	"strconv"
	"time"

	"columnstore/internal/scan"
)

// ParseWindow parses an inclusive query window from RFC3339 text.
func ParseWindow(fromText, toText string) (time.Time, time.Time, error) {
	from, err := time.Parse(time.RFC3339, fromText)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("from must be RFC3339: %w", err)
	}
	to, err := time.Parse(time.RFC3339, toText)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("to must be RFC3339: %w", err)
	}
	if to.Before(from) {
		return time.Time{}, time.Time{}, fmt.Errorf("query window is reversed")
	}
	return from, to, nil
}

// BuildPredicate converts operator text and a threshold into a scan predicate.
// Supported operators are ge, gt, le, lt, eq and ne.
func BuildPredicate(column, opText, valueText string) (*scan.Predicate, error) {
	value, err := strconv.ParseFloat(valueText, 64)
	if err != nil {
		return nil, fmt.Errorf("threshold must be numeric: %w", err)
	}
	var op scan.Op
	switch opText {
	case "ge":
		op = scan.OpGE
	case "gt":
		op = scan.OpGT
	case "le":
		op = scan.OpLE
	case "lt":
		op = scan.OpLT
	case "eq":
		op = scan.OpEQ
	case "ne":
		op = scan.OpNE
	default:
		return nil, fmt.Errorf("unsupported operator %q", opText)
	}
	return &scan.Predicate{Column: column, Op: op, Value: value}, nil
}
