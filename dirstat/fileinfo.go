package dirstat

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileInfo represents information about a file or directory.
type FileInfo struct {
	Name     string
	IsDir    bool
	DiskSize int64 // Size in bytes (for files) or sum of sizes (for directories)
	Children []*FileInfo
}

// FromPath takes a path and analyzes whether it's a file or directory.
// It returns a FileInfo struct populated with the appropriate information.
func FromPath(path string) (*FileInfo, error) {
	// Get file or directory info
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("unable to access path %s: %v", path, err)
	}

	// If it's a directory, recurse through contents and gather information
	if fileInfo.IsDir() {
		return analyzeDirectory(path)
	} else {
		// Otherwise, just return the file's information
		size, err := getDiskSize(path)
		if err != nil {
			return nil, fmt.Errorf("unable to get file size for %s: %v", path, err)
		}
		return &FileInfo{
			Name:     filepath.Base(path),
			IsDir:    false,
			DiskSize: size,
			Children: nil,
		}, nil
	}
}

// analyzeDirectory recursively gathers information about a directory and its contents.
func analyzeDirectory(dirPath string) (*FileInfo, error) {
	// Create the root directory's FileInfo
	root := &FileInfo{
		Name:     filepath.Base(dirPath),
		IsDir:    true,
		DiskSize: 0,
		Children: []*FileInfo{},
	}

	// Walk the directory to gather size and file information
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root itself
		if path == dirPath {
			return nil
		}

		// Create a new FileInfo for each file or directory encountered
		child := &FileInfo{
			Name:     filepath.Base(path),
			IsDir:    info.IsDir(),
			DiskSize: info.Size(),
			Children: nil,
		}

		// If it's a directory, calculate its total size recursively
		if info.IsDir() {
			dirInfo, err := analyzeDirectory(path)
			if err != nil {
				return err
			}
			child.Children = dirInfo.Children
			child.DiskSize = dirInfo.DiskSize
		} else {
			size, err := getDiskSize(path)
			if err != nil {
				return err
			}
			child.DiskSize = size
		}

		// Add the child info to the root's Children and accumulate the size
		root.Children = append(root.Children, child)
		root.DiskSize += child.DiskSize
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error while walking through directory %s: %v", dirPath, err)
	}

	return root, nil
}
