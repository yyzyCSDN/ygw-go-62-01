package meta

import (
	"errors"
	"time"
)

var (
	ErrUnknownTable  = errors.New("table has no registered metadata")
	ErrUnknownColumn = errors.New("column is missing from table schema")
)

// ColumnMeta carries the schema version of one table.
type ColumnMeta struct {
	Table         string
	Schema        Schema
	SchemaVersion Version
	Compressed    bool
	UpdatedAt     time.Time
}

// Column returns the definition of a column by name.
func (m *ColumnMeta) Column(name string) (ColumnDef, error) {
	pos := m.Schema.Position(name)
	if pos < 0 {
		return ColumnDef{}, ErrUnknownColumn
	}
	return m.Schema.Columns[pos], nil
}

// AtVersion reports whether the meta matches the given version number.
func (m *ColumnMeta) AtVersion(number uint64) bool {
	return m.SchemaVersion.Number == number
}
