package compress

import (
	"testing"
	"time"

	"columnstore/internal/col"
	"columnstore/internal/meta"
)

func TestDictionaryRoundTrip(t *testing.T) {
	dict := NewDictionary(1)
	id, err := dict.Encode("alpha")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := dict.Encode("beta"); err != nil {
		t.Fatal(err)
	}
	got, err := dict.Decode(id)
	if err != nil {
		t.Fatal(err)
	}
	if got != "alpha" {
		t.Fatalf("expected alpha, got %q", got)
	}
}

func TestCodecRoundTrip(t *testing.T) {
	dict := NewDictionary(1)
	codec := NewCodec(dict)
	rows := []col.Value{col.Str("a"), col.Str("b"), col.Str("a")}
	payload, err := codec.EncodeRows(rows)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := codec.DecodeRows(payload)
	if err != nil {
		t.Fatal(err)
	}
	if len(decoded) != 3 || decoded[2].Str != "a" {
		t.Fatalf("unexpected decoded rows: %+v", decoded)
	}
}

func TestCompressThenDecode(t *testing.T) {
	dicts := NewDictRegistry()
	metaRegistry := meta.NewRegistry()
	_, err := metaRegistry.Update(meta.BuildSchema("t", []string{"region"}, nil), time.Now())
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(dicts, metaRegistry)
	block := col.NewBlock("t", "region", []col.Value{col.Str("east"), col.Str("west"), col.Str("east")}, time.Now())
	if err := block.Seal("seg-1"); err != nil {
		t.Fatal(err)
	}
	if err := service.CompressColumn(block); err != nil {
		t.Fatal(err)
	}
	if block.State != col.StateCompressed {
		t.Fatal("block must be compressed after success")
	}
	rows, err := service.DecodeColumn(block)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 3 || rows[1].Str != "west" {
		t.Fatalf("unexpected decoded values: %+v", rows)
	}
}

func TestFingerprintStable(t *testing.T) {
	a := Fingerprint(7, []byte{1, 2, 3})
	b := Fingerprint(7, []byte{1, 2, 3})
	if a != b {
		t.Fatalf("fingerprint must be stable")
	}
}
