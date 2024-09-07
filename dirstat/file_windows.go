//go:build windows
// +build windows

package dirstat

import (
	"syscall"
)

func getDiskSize(path string) (int64, error) {
	// Convert the path to a UTF-16 pointer
	widePath, err := syscall.UTF16FromString(path)
	if err != nil {
		return 0, err
	}
	widePath = append(widePath, 0) // Append null terminator

	handle, err := syscall.CreateFile(
		&widePath[0],                  // Pointer to the UTF-16 string
		syscall.GENERIC_READ,          // Desired access
		syscall.FILE_SHARE_READ,       // Share mode
		nil,                           // Security attributes
		syscall.OPEN_EXISTING,         // Creation disposition
		syscall.FILE_ATTRIBUTE_NORMAL, // Flags and attributes
		0,                             // Template file
	)
	if err != nil {
		return 0, err
	}
	defer syscall.CloseHandle(handle)

	var fileInformation syscall.ByHandleFileInformation
	err = syscall.GetFileInformationByHandle(handle, &fileInformation)
	if err != nil {
		return 0, err
	}

	// Combine high and low parts of the file size
	fileSize := int64(fileInformation.FileSizeHigh)<<32 | int64(fileInformation.FileSizeLow)
	return fileSize, nil
}
