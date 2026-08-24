package col

import (
	"fmt"
	"sort"
	"sync"
)

// Store is the in-memory registry of column blocks plus the segment manager.
type Store struct {
	mu       sync.RWMutex
	blocks   map[uint64]*ColumnBlock
	index    *BlockIndex
	segments *SegmentManager
	nextID   uint64
}

// NewStore builds an empty store over the given segment directory.
func NewStore(dir string) *Store {
	return &Store{
		blocks:   make(map[uint64]*ColumnBlock),
		index:    NewBlockIndex(),
		segments: NewSegmentManager(dir),
		nextID:   1,
	}
}

// AddBlock registers a new block and assigns it a store-wide id.
func (s *Store) AddBlock(block *ColumnBlock) (uint64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if block == nil {
		return 0, ErrEmptyColumn
	}
	id := s.nextID
	s.nextID++
	block.ID = id
	block.Segment = fmt.Sprintf("seg-%d", id)
	s.blocks[id] = block
	s.index.Add(block)
	return id, nil
}

// GetBlock returns a block by id.
func (s *Store) GetBlock(id uint64) (*ColumnBlock, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.blocks[id]
	if !ok {
		return nil, ErrBlockNotFound
	}
	return b, nil
}

// RemoveBlock drops a block from the store.
func (s *Store) RemoveBlock(id uint64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	block, ok := s.blocks[id]
	if !ok {
		return ErrBlockNotFound
	}
	delete(s.blocks, id)
	s.index.Remove(block)
	return nil
}

// Tables returns the table names present in the block index.
func (s *Store) Tables() []string {
	return s.index.Tables()
}

// ListBlocks returns all blocks ordered by id.
func (s *Store) ListBlocks() []*ColumnBlock {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*ColumnBlock, 0, len(s.blocks))
	for _, b := range s.blocks {
		out = append(out, b)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// OpenHandles reports open segment handles across the store.
func (s *Store) OpenHandles() int {
	return s.segments.OpenHandles()
}

// BlockCount returns the number of live blocks.
func (s *Store) BlockCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.blocks)
}
