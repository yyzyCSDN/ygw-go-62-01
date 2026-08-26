package query

import (
	"fmt"
	"strings"

	"columnstore/internal/col"
)

// Row is one projected output row.
type Row struct {
	Values []col.Value
}

// Result is the materialized output of a query.
type Result struct {
	Columns []string
	Rows    []Row
	Scanned int
}

// Empty returns an empty result with the given column names.
func Empty(columns []string) *Result {
	return &Result{Columns: append([]string(nil), columns...)}
}

// Append adds a row of values.
func (r *Result) Append(values []col.Value) {
	r.Rows = append(r.Rows, Row{Values: values})
}

// Text renders rows as aligned columns for the console.
func (r *Result) Text() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", strings.Join(r.Columns, "\t"))
	for _, row := range r.Rows {
		parts := make([]string, 0, len(row.Values))
		for _, v := range row.Values {
			if v.Kind == col.KindString {
				parts = append(parts, v.Str)
			} else {
				parts = append(parts, fmt.Sprintf("%v", v.Num))
			}
		}
		fmt.Fprintf(&b, "%s\n", strings.Join(parts, "\t"))
	}
	return b.String()
}
