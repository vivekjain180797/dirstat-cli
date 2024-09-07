//go:build linux
// +build linux

package dirstat

import (
	"syscall"
)

func getDiskSize(path string) (int64, error) {
	// Implement Linux-specific logic to get file size
	var stat syscall.Stat_t

	err := syscall.Stat(path, &stat)
	if err != nil {
		return 0, err
	}

	return stat.Size, nil
}
