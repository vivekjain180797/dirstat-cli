package dirstat

import (
	"os"
	"path/filepath"
)

// DiskItem holds directory and file information.
type DiskItem struct {
	Name     string
	Size     uint64
	Children []*DiskItem
}

// AnalyzeDir creates a DiskItem from the directory structure.
func AnalyzeDir(path string, maxDepth int, minPercent float64) (*DiskItem, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	root := &DiskItem{
		Name: filepath.Base(path),
		Size: 0,
	}

	// If it's a directory, analyze contents
	if info.IsDir() {
		root.Children, err = analyzeDirectoryContents(path, maxDepth)
		if err != nil {
			return nil, err
		}

		// Compute total size
		for _, child := range root.Children {
			root.Size += child.Size
		}
	} else {
		// It's a file, just get the size
		root.Size = uint64(info.Size())
	}

	return root, nil
}

// analyzeDirectoryContents recursively analyzes the directory and its contents.
func analyzeDirectoryContents(path string, depth int) ([]*DiskItem, error) {
	if depth == 0 {
		return nil, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var items []*DiskItem
	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		info, err := os.Stat(fullPath)
		if err != nil {
			return nil, err
		}

		item := &DiskItem{
			Name: entry.Name(),
			Size: uint64(info.Size()),
		}

		if info.IsDir() {
			item.Children, err = analyzeDirectoryContents(fullPath, depth-1)
			if err != nil {
				return nil, err
			}

			// Add directory size (sum of children's sizes)
			for _, child := range item.Children {
				item.Size += child.Size
			}
		}

		items = append(items, item)
	}
	return items, nil
}
