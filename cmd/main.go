package main

import (
	"fmt"
	"os"
	"path/filepath"
	"project/dirstat"
	"project/display"
)

func main() {
	// Default config values
	maxDepth := 2
	minPercent := 0.1
	targetDir := "."

	// Allow directory path as command-line argument
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	// Get absolute path
	absDir, err := filepath.Abs(targetDir)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}

	// Analyze directory and create disk item tree
	rootDiskItem, err := dirstat.AnalyzeDir(absDir, maxDepth, minPercent)
	if err != nil {
		fmt.Println("Error analyzing directory:", err)
		return
	}

	// Display the results
	display.PrintDiskItem(rootDiskItem, maxDepth, minPercent)
}
