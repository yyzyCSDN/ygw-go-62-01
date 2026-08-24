package compress

import (
	"encoding/binary"

	"github.com/cespare/xxhash/v2"
)

// Fingerprint hashes dictionary ids plus payload to detect corruption.
func Fingerprint(dictID uint64, payload []byte) uint64 {
	h := xxhash.New()
	var id [8]byte
	binary.LittleEndian.PutUint64(id[:], dictID)
	_, _ = h.Write(id[:])
	_, _ = h.Write(payload)
	return h.Sum64()
}
