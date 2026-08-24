package compress

import (
	"encoding/binary"
	"fmt"

	"columnstore/internal/col"
)

// Codec turns a column of string values into dictionary ids and back.
type Codec struct {
	dict *Dictionary
}

// NewCodec creates a codec over one dictionary.
func NewCodec(dict *Dictionary) *Codec {
	return &Codec{dict: dict}
}

// EncodeRows converts string values into packed dictionary ids.
func (c *Codec) EncodeRows(rows []col.Value) ([]byte, error) {
	if c.dict == nil {
		return nil, ErrNoDictionary
	}
	out := make([]byte, 0, len(rows)*4)
	buf := make([]byte, 4)
	for i, row := range rows {
		if row.Kind != col.KindString {
			return nil, fmt.Errorf("%w at row %d", ErrValueNotString, i)
		}
		id, err := c.dict.Encode(row.Str)
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint32(buf, id)
		out = append(out, buf...)
	}
	return out, nil
}

// DecodeRows expands packed ids back into string values.
func (c *Codec) DecodeRows(payload []byte) ([]col.Value, error) {
	if c.dict == nil {
		return nil, ErrNoDictionary
	}
	if len(payload)%4 != 0 {
		return nil, fmt.Errorf("malformed dictionary payload length %d", len(payload))
	}
	rows := make([]col.Value, 0, len(payload)/4)
	for pos := 0; pos < len(payload); pos += 4 {
		id := binary.LittleEndian.Uint32(payload[pos : pos+4])
		value, err := c.dict.Decode(id)
		if err != nil {
			return nil, err
		}
		rows = append(rows, col.Str(value))
	}
	return rows, nil
}
