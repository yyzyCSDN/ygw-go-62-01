package meta

import "columnstore/internal/col"

// ColumnDef describes one column inside a table schema.
type ColumnDef struct {
	Name string
	Kind col.Kind
	Pos  int
}

// Schema is the ordered column layout of a table.
type Schema struct {
	Table   string
	Columns []ColumnDef
}

// Position returns the zero-based position of a column, or -1.
func (s Schema) Position(name string) int {
	for _, def := range s.Columns {
		if def.Name == name {
			return def.Pos
		}
	}
	return -1
}

// BuildSchema creates a schema from column names and kinds.
func BuildSchema(table string, names []string, kinds []col.Kind) Schema {
	columns := make([]ColumnDef, 0, len(names))
	for i, name := range names {
		kind := col.KindInt
		if i < len(kinds) {
			kind = kinds[i]
		}
		columns = append(columns, ColumnDef{Name: name, Kind: kind, Pos: i})
	}
	return Schema{Table: table, Columns: columns}
}
