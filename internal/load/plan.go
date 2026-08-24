package load

import (
	"columnstore/internal/col"
	"columnstore/internal/meta"
	"columnstore/internal/part"
)

// Plan is the prepared state of one load operation.
type Plan struct {
	Table   string
	Batches int
	Schema  meta.Schema
}

// BuildPlan derives the schema and batch count for a source.
func BuildPlan(table string, src Source) Plan {
	batch := src.Batches()
	count := len(batch)
	columns := make([]string, 0, count)
	kinds := make([]col.Kind, 0, count)
	seen := map[string]bool{}
	for _, b := range batch {
		if !seen[b.Column] {
			seen[b.Column] = true
			columns = append(columns, b.Column)
			kinds = append(kinds, ColumnKind(b))
		}
	}
	return Plan{
		Table:   table,
		Batches: count,
		Schema:  meta.BuildSchema(table, columns, kinds),
	}
}

// PartitionFor builds a partition record for a batch.
func PartitionFor(table string, b Batch, id uint64) *part.Partition {
	return part.NewPartition(table, b.Start, b.End, id)
}
