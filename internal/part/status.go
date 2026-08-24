package part

// Status is a compact text summary of a partition for the console API.
type Status struct {
	ID       uint64
	Table    string
	Start    string
	End      string
	State    string
	BlockIDs []uint64
}

// Describe renders a partition as a status record.
func Describe(p *Partition) Status {
	return Status{
		ID:       p.ID,
		Table:    p.Table,
		Start:    p.Start.Format("2006-01-02T15:04:05"),
		End:      p.End.Format("2006-01-02T15:04:05"),
		State:    p.State.String(),
		BlockIDs: append([]uint64(nil), p.BlockIDs...),
	}
}

// DescribeAll renders every partition of a table.
func DescribeAll(parts []*Partition) []Status {
	out := make([]Status, 0, len(parts))
	for _, p := range parts {
		out = append(out, Describe(p))
	}
	return out
}
