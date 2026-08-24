package gc

import (
	"os"
)

// DiskUsage reports the total bytes of segment files in the data directory.
func DiskUsage(dir string) (int64, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	var total int64
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		info, err := entry.Info()
		if err != nil {
			return 0, err
		}
		total += info.Size()
	}
	return total, nil
}
