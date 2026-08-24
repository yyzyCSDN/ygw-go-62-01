package meta

// Catalog wraps the registry with read-through access used by the query layer.
type Catalog struct {
	registry *Registry
}

// NewCatalog builds a query-facing catalog over a registry.
func NewCatalog(registry *Registry) *Catalog {
	return &Catalog{registry: registry}
}

// SchemaFor returns the latest schema for a table.
func (c *Catalog) SchemaFor(table string) (Schema, error) {
	meta, err := c.registry.Lookup(table)
	if err != nil {
		return Schema{}, err
	}
	return meta.Schema, nil
}

// MetaFor returns the latest metadata for a table.
func (c *Catalog) MetaFor(table string) (*ColumnMeta, error) {
	return c.registry.Lookup(table)
}
