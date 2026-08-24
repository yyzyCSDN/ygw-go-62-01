package col

import (
	"io"
	"os"
	"path/filepath"
	"sync/atomic"
)

// SegmentManager owns the on-disk segment files that back immutable blocks.
type SegmentManager struct {
	dir         string
	openHandles int32
}

// NewSegmentManager creates a manager rooted at dir.
func NewSegmentManager(dir string) *SegmentManager {
	return &SegmentManager{dir: dir}
}

// OpenHandles returns how many segment file handles are currently open.
func (m *SegmentManager) OpenHandles() int {
	return int(atomic.LoadInt32(&m.openHandles))
}

// Open opens a segment file for reading and tracks the handle.
func (m *SegmentManager) Open(name string) (*os.File, error) {
	path := filepath.Join(m.dir, name)
	f, err := os.Open(path)
	if err != nil {
		return nil, ErrSegmentOpen
	}
	atomic.AddInt32(&m.openHandles, 1)
	return f, nil
}

// Close releases a previously opened handle.
func (m *SegmentManager) Close(f *os.File) {
	if f == nil {
		return
	}
	_ = f.Close()
	atomic.AddInt32(&m.openHandles, -1)
}

// ReadSegmentBytes reads the full content of a segment and guarantees closure.
func (m *SegmentManager) ReadSegmentBytes(name string) ([]byte, error) {
	f, err := m.Open(name)
	if err != nil {
		return nil, err
	}
	defer m.Close(f)
	return io.ReadAll(f)
}

// WriteSegment writes payload to a segment file and closes it immediately.
func (m *SegmentManager) WriteSegment(name string, payload []byte) error {
	path := filepath.Join(m.dir, name)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	if _, err = f.Write(payload); err != nil {
		_ = f.Close()
		return err
	}
	return f.Close()
}

// RemoveSegment deletes a segment file from disk.
func (m *SegmentManager) RemoveSegment(name string) error {
	return os.Remove(filepath.Join(m.dir, name))
}
