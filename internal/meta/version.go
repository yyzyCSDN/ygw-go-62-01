package meta

// Version identifies a schema revision for a table.
type Version struct {
	Table   string
	Number  uint64
	Changed bool
}

// Next returns the version following this one.
func (v Version) Next() Version {
	return Version{Table: v.Table, Number: v.Number + 1, Changed: true}
}

// Equal reports whether two versions describe the same revision.
func (v Version) Equal(other Version) bool {
	return v.Table == other.Table && v.Number == other.Number
}
