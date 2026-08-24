package col

import "sync"

// BlockIndex maps table and column names to the blocks that carry them.
type BlockIndex struct {
	mu      sync.RWMutex
	byTable map[string]map[uint64]*ColumnBlock
	byCol   map[string]map[uint64]*ColumnBlock
}

// NewBlockIndex creates an empty block index.
func NewBlockIndex() *BlockIndex {
	return &BlockIndex{
		byTable: make(map[string]map[uint64]*ColumnBlock),
		byCol:   make(map[string]map[uint64]*ColumnBlock),
	}
}

// Add registers a block under its table and column keys.
func (x *BlockIndex) Add(block *ColumnBlock) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if x.byTable[block.Table] == nil {
		x.byTable[block.Table] = make(map[uint64]*ColumnBlock)
	}
	x.byTable[block.Table][block.ID] = block
	key := block.Table + "." + block.Column
	if x.byCol[key] == nil {
		x.byCol[key] = make(map[uint64]*ColumnBlock)
	}
	x.byCol[key][block.ID] = block
}

// Remove unregisters a block from the index.
func (x *BlockIndex) Remove(block *ColumnBlock) {
	x.mu.Lock()
	defer x.mu.Unlock()
	if blocks := x.byTable[block.Table]; blocks != nil {
		delete(blocks, block.ID)
		if len(blocks) == 0 {
			delete(x.byTable, block.Table)
		}
	}
	key := block.Table + "." + block.Column
	if blocks := x.byCol[key]; blocks != nil {
		delete(blocks, block.ID)
		if len(blocks) == 0 {
			delete(x.byCol, key)
		}
	}
}

// TableBlocks returns every block of a table ordered by id.
func (x *BlockIndex) TableBlocks(table string) []*ColumnBlock {
	x.mu.RLock()
	defer x.mu.RUnlock()
	blocks := x.byTable[table]
	out := make([]*ColumnBlock, 0, len(blocks))
	for _, block := range blocks {
		out = append(out, block)
	}
	sortBlocks(out)
	return out
}

// Tables returns the distinct table names present in the index.
func (x *BlockIndex) Tables() []string {
	x.mu.RLock()
	defer x.mu.RUnlock()
	out := make([]string, 0, len(x.byTable))
	for name := range x.byTable {
		out = append(out, name)
	}
	sortStrings(out)
	return out
}

func sortBlocks(blocks []*ColumnBlock) {
	for i := 1; i < len(blocks); i++ {
		for j := i; j > 0 && blocks[j-1].ID > blocks[j].ID; j-- {
			blocks[j-1], blocks[j] = blocks[j], blocks[j-1]
		}
	}
}

func sortStrings(values []string) {
	for i := 1; i < len(values); i++ {
		for j := i; j > 0 && values[j-1] > values[j]; j-- {
			values[j-1], values[j] = values[j], values[j-1]
		}
	}
}
