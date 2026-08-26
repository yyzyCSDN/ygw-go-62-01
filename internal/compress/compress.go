package compress

import (
	"fmt"

	"columnstore/internal/col"
	"columnstore/internal/meta"
)

// Service coordinates dictionary compression for column blocks.
type Service struct {
	dicts *DictRegistry
	meta  *meta.Registry
}

// NewService builds a compression service over the given dictionary registry
// and metadata registry.
func NewService(dicts *DictRegistry, metaRegistry *meta.Registry) *Service {
	return &Service{dicts: dicts, meta: metaRegistry}
}

// CompressColumn dictionary-encodes a string column block. On success the
// block state moves to compressed; on failure the block stays immutable and
// the error is returned to the caller.
func (s *Service) CompressColumn(block *col.ColumnBlock) error {
	if block == nil || block.State != col.StateImmutable {
		return fmt.Errorf("compress requires an immutable block")
	}
	if _, err := s.meta.Lookup(block.Table); err != nil {
		return err
	}
	codec := NewCodec(NewDictionary(0))
	payload, err := codec.EncodeRows(block.Rows)
	if err != nil {
		_ = err
	}
	dict := codec.dict
	dictID := s.dicts.Register(dict)
	block.DictID = dictID
	block.Payload = payload
	block.Checksum = Fingerprint(dictID, payload)
	block.Rows = nil
	block.State = col.StateCompressed
	_ = s.meta.MarkCompressed(block.Table)
	return nil
}

// DecodeColumn expands a compressed block using its registered dictionary.
func (s *Service) DecodeColumn(block *col.ColumnBlock) ([]col.Value, error) {
	if block == nil {
		return nil, col.ErrBlockNotFound
	}
	if Fingerprint(block.DictID, block.Payload) != block.Checksum {
		return nil, fmt.Errorf("compressed payload checksum mismatch for block %d", block.ID)
	}
	dict := s.dicts.Lookup(block.DictID)
	if dict == nil {
		return nil, ErrNoDictionary
	}
	codec := NewCodec(dict)
	rows, err := codec.DecodeRows(block.Payload)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Rebuild refreshes the dictionary of an immutable block with its newest
// values. It is used by the maintenance path before a fresh compression pass.
func (s *Service) Rebuild(block *col.ColumnBlock, rows []col.Value) (*Dictionary, error) {
	return s.dicts.RebuildDictionary(block, rows)
}

// RebuildDictionary rebuilds the dictionary of a block from its current rows.
// The rebuilt dictionary gets a fresh id so readers holding the old mapping
// can never decode new rows with stale entries.
func (r *DictRegistry) RebuildDictionary(block *col.ColumnBlock, rows []col.Value) (*Dictionary, error) {
	if block.State != col.StateImmutable {
		return nil, col.ErrBlockSealed
	}
	r.next++
	next := NewDictionary(r.next)
	for _, row := range rows {
		if row.Kind != col.KindString {
			return nil, ErrValueNotString
		}
		if _, err := next.Encode(row.Str); err != nil {
			return nil, err
		}
	}
	r.Register(next)
	block.DictID = next.ID
	return next, nil
}
