package compress

import (
	"errors"
	"sort"
)

var (
	ErrDictionaryFull = errors.New("dictionary exceeds capacity")
	ErrValueNotString = errors.New("dictionary encoding requires string values")
	ErrNoDictionary   = errors.New("block has no active dictionary")
)

const DefaultDictCapacity = 4096

// Dictionary maps string values to dense ids for a low-cardinality column.
type Dictionary struct {
	ID      uint64
	lookup  map[string]uint32
	reverse []string
}

// NewDictionary builds an empty dictionary with the given id.
func NewDictionary(id uint64) *Dictionary {
	return &Dictionary{
		ID:     id,
		lookup: make(map[string]uint32),
	}
}

// Encode returns the id for a value, adding it when unknown.
func (d *Dictionary) Encode(value string) (uint32, error) {
	if id, ok := d.lookup[value]; ok {
		return id, nil
	}
	if len(d.reverse) >= DefaultDictCapacity {
		return 0, ErrDictionaryFull
	}
	id := uint32(len(d.reverse))
	d.lookup[value] = id
	d.reverse = append(d.reverse, value)
	return id, nil
}

// Decode returns the original value for an id.
func (d *Dictionary) Decode(id uint32) (string, error) {
	if int(id) >= len(d.reverse) {
		return "", ErrNoDictionary
	}
	return d.reverse[id], nil
}

// Size returns the number of entries.
func (d *Dictionary) Size() int {
	return len(d.reverse)
}

// Values returns the reverse mapping sorted by id.
func (d *Dictionary) Values() []string {
	out := append([]string(nil), d.reverse...)
	sort.Strings(out)
	return out
}

// DictRegistry owns every dictionary by id.
type DictRegistry struct {
	byID map[uint64]*Dictionary
	next uint64
}

// NewDictRegistry creates an empty dictionary registry.
func NewDictRegistry() *DictRegistry {
	return &DictRegistry{byID: make(map[uint64]*Dictionary)}
}

// Register adds a dictionary and returns its id.
func (r *DictRegistry) Register(d *Dictionary) uint64 {
	if d.ID == 0 {
		r.next++
		d.ID = r.next
	}
	r.byID[d.ID] = d
	return d.ID
}

// Lookup resolves a dictionary by id.
func (r *DictRegistry) Lookup(id uint64) *Dictionary {
	return r.byID[id]
}
