package part

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/cespare/xxhash/v2"
)

// Registry owns every partition of every table.
type Registry struct {
	mu       sync.RWMutex
	byID     map[uint64]*Partition
	byTable  map[string][]*Partition
	nextPart uint64
}

// NewRegistry builds an empty partition registry.
func NewRegistry() *Registry {
	return &Registry{
		byID:    make(map[uint64]*Partition),
		byTable: make(map[string][]*Partition),
	}
}

// Register inserts a partition and assigns it a registry-wide id.
func (r *Registry) Register(p *Partition) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p == nil {
		return 0, fmt.Errorf("cannot register a nil partition")
	}
	hash := partitionHash(p.Table, p.Start, p.End)
	if _, exists := r.byID[hash]; !exists {
		p.ID = hash
	} else {
		r.nextPart++
		p.ID = r.nextPart
	}
	r.byID[p.ID] = p
	r.byTable[p.Table] = append(r.byTable[p.Table], p)
	return p.ID, nil
}

// partitionHash derives a stable id from a partition's table and interval.
func partitionHash(table string, start, end time.Time) uint64 {
	h := xxhash.New()
	_, _ = h.WriteString(table)
	_, _ = h.Write([]byte{0})
	_, _ = h.WriteString(start.Format(time.RFC3339Nano))
	_, _ = h.Write([]byte{0})
	_, _ = h.WriteString(end.Format(time.RFC3339Nano))
	return h.Sum64()
}


// Remove drops a partition and its table index entry.
func (r *Registry) Remove(id uint64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	p, ok := r.byID[id]
	if !ok {
		return fmt.Errorf("partition %d not found", id)
	}
	delete(r.byID, id)
	rows := r.byTable[p.Table]
	for i, item := range rows {
		if item.ID == id {
			r.byTable[p.Table] = append(rows[:i], rows[i+1:]...)
			break
		}
	}
	return nil
}

// Table returns all partitions of a table ordered by start time.
func (r *Registry) Table(table string) []*Partition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := append([]*Partition(nil), r.byTable[table]...)
	sort.Slice(out, func(i, j int) bool { return out[i].Start.Before(out[j].Start) })
	return out
}

// Count returns the number of registered partitions.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.byID)
}
