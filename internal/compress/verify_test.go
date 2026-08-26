package compress_test

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/compress"
	"columnstore/internal/meta"
)

func TestCompressErrorNotSwallowed(t *testing.T) {
	dicts := compress.NewDictRegistry()
	metaRegistry := meta.NewRegistry()
	now := time.Date(2026, 8, 24, 0, 0, 0, 0, time.UTC)
	if _, err := metaRegistry.Update(meta.BuildSchema("t", []string{"c"}, nil), now); err != nil {
		t.Fatal(err)
	}
	service := compress.NewService(dicts, metaRegistry)
	block := col.NewBlock("t", "c", []col.Value{col.Str("a"), col.Int(1)}, now)
	if err := block.Seal("seg-x"); err != nil {
		t.Fatal(err)
	}
	err := service.CompressColumn(block)
	if err == nil {
		t.Fatal("compression failure must be surfaced")
	}
	if block.State == col.StateCompressed {
		t.Fatal("block must stay immutable after a failed compression")
	}
}
