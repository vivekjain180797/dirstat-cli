package display

import (
	"fmt"
	"project/dirstat"
)

// DisplayInfo holds display-related info like indentation level.
type DisplayInfo struct {
	Level    int
	Fraction float64
	Last     bool
	Indents  string
}

// PrintDiskItem prints a directory tree.
func PrintDiskItem(item *dirstat.DiskItem, maxDepth int, minPercent float64) {
	info := DisplayInfo{
		Fraction: 100.0,
		Level:    0,
		Last:     true,
		Indents:  "",
	}

	printIndentedItem(item, &info)

	if len(item.Children) > 0 {
		for _, child := range item.Children {
			newInfo := info.addChild(100.0 * float64(child.Size) / float64(item.Size))
			printIndentedItem(child, &newInfo)
		}
	}
}

func printIndentedItem(item *dirstat.DiskItem, info *DisplayInfo) {
	prefix := "├──"
	if info.Last {
		prefix = "└──"
	}

	fmt.Printf("%s%s [%.2f%%] %s\n", info.Indents, prefix, info.Fraction, item.Name)
}

// Adds a child with proper indentation.
func (info *DisplayInfo) addChild(fraction float64) DisplayInfo {
	indents := info.Indents
	if info.Last {
		indents += "    "
	} else {
		indents += "│   "
	}

	return DisplayInfo{
		Level:    info.Level + 1,
		Fraction: fraction,
		Last:     false,
		Indents:  indents,
	}
}
