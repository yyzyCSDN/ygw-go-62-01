package load

import "sync"

// Manifest tracks how many blocks each table has received from loads.
type Manifest struct {
	mu      sync.RWMutex
	byTable map[string]int
}

// NewManifest builds an empty load manifest.
func NewManifest() *Manifest {
	return &Manifest{byTable: make(map[string]int)}
}

// Add records one more batch for a table.
func (m *Manifest) Add(table string, blocks int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byTable[table] += blocks
}

// Rollback subtracts blocks that a failed load rolled back.
func (m *Manifest) Rollback(table string, blocks int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byTable[table] -= blocks
	if m.byTable[table] < 0 {
		m.byTable[table] = 0
	}
}

// TableBlocks returns how many blocks a table currently has.
func (m *Manifest) TableBlocks(table string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.byTable[table]
}

// Snapshot copies the whole manifest.
func (m *Manifest) Snapshot() map[string]int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[string]int, len(m.byTable))
	for table, count := range m.byTable {
		out[table] = count
	}
	return out
}
