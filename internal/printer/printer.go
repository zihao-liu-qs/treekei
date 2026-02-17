package printer

import (
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/mattn/go-runewidth"

	"github.com/zihao-liu-qs/treekei/internal/tree"
)

func init() {
	runewidth.DefaultCondition.EastAsianWidth = false
}

type Printer struct {
	maxWidth int
	maxDepth int
	NoColor  bool
	DirOnly  bool
}

var (
	dirColor    = color.New(color.FgCyan)
	treeSymbol  = color.New(color.Faint)
	smallLines  = color.New(color.Faint)    // 0-100 lines
	normalLines = color.New(color.FgGreen)  // 101-300 lines
	largeLines  = color.New(color.FgYellow) // 301-600 lines
	hugeLines   = color.New(color.FgRed)    // 600+ lines
)

func NewPrinter(maxDepth int, noColor bool, dirOnly bool) *Printer {
	if noColor {
		color.NoColor = true
	}
	return &Printer{
		maxDepth: maxDepth,
		NoColor:  noColor,
		DirOnly:  dirOnly,
	}
}

func getLineColor(lines int) *color.Color {
	switch {
	case lines <= 100:
		return smallLines
	case lines <= 300:
		return normalLines
	case lines <= 600:
		return largeLines
	default:
		return hugeLines
	}
}

func (p *Printer) PrintTree(node *tree.Node, prefix string, isLast, isRoot bool, depth int) {
	if p.maxDepth > 0 && depth > p.maxDepth {
		return
	}

	// Skip files if dirOnly is enabled
	if p.DirOnly && !node.IsDir && !isRoot {
		return
	}

	if isRoot {
		p.calSetMaxWidth(node, 0, 0)
	}

	var line string
	var nameStr string

	if node.IsDir {
		nameStr = dirColor.Sprint(node.Name)
	} else {
		nameStr = node.Name
	}

	var lineColor *color.Color
	var linesStr string
	lineColor = getLineColor(node.Lines)
	if node.IsBinary {
		linesStr = lineColor.Sprintf("[bin]")
	} else {
		linesStr = lineColor.Sprintf("%d", node.Lines)
	}

	if isRoot {
		plainText := node.Name
		padding := p.maxWidth - runewidth.StringWidth(plainText)
		line = fmt.Sprintf("%s%*s%s\n", nameStr, padding, "", linesStr)
	} else {
		var symbol string
		if !isLast {
			symbol = treeSymbol.Sprint("├── ")
		} else {
			symbol = treeSymbol.Sprint("└── ")
		}

		coloredPrefix := treeSymbol.Sprint(prefix)

		plainText := prefix + "├── " + node.Name
		padding := p.maxWidth - runewidth.StringWidth(plainText)
		line = fmt.Sprintf("%s%s%s%*s%s\n", coloredPrefix, symbol, nameStr, padding, "", linesStr)
	}

	fmt.Print(line)

	// Filter children if dirOnly is enabled
	var children []*tree.Node
	if p.DirOnly {
		for _, child := range node.Children {
			if child.IsDir {
				children = append(children, child)
			}
		}
	} else {
		children = node.Children
	}

	var nextPrefix string
	var nextIsLast bool
	for id, childNode := range children {
		if id == len(children)-1 {
			nextIsLast = true
		}
		if isLast {
			nextPrefix = prefix + "    "
		} else {
			nextPrefix = prefix + "│   "
		}
		p.PrintTree(childNode, nextPrefix, nextIsLast, false, depth+1)
	}
}

func (p *Printer) calSetMaxWidth(node *tree.Node, prefixLen, depth int) {
	if p.maxDepth > 0 && depth > p.maxDepth {
		return
	}

	// Skip files if dirOnly is enabled
	if p.DirOnly && !node.IsDir {
		return
	}

	p.maxWidth = max(p.maxWidth, prefixLen+runewidth.StringWidth(node.Name+" "+strconv.Itoa(node.Lines))+4)

	for _, childNode := range node.Children {
		// Skip files in width calculation if dirOnly is enabled
		if p.DirOnly && !childNode.IsDir {
			continue
		}
		p.calSetMaxWidth(childNode, prefixLen+runewidth.StringWidth("├── "), depth+1)
	}
}
