package printer

import (
	"fmt"
	"strconv"

	"github.com/mattn/go-runewidth"

	"github.com/zihao-liu-qs/treekei/internal/tree"
)

func init() {
	runewidth.DefaultCondition.EastAsianWidth = false
}

type Printer struct {
	maxWidth int
}

func (p *Printer) PrintTree(node *tree.Node, prefix string, isLast, isRoot bool) {
	if isRoot {
		p.calSetMaxWidth(node, 0)
	}

	if isRoot {
		fmt.Printf("%-s%*d\n", node.Name, p.maxWidth-runewidth.StringWidth(node.Name), node.Lines)
	} else if !isLast {
		fmt.Printf("%-s%*d\n", prefix+"├── "+node.Name, p.maxWidth-runewidth.StringWidth(prefix+"├── "+node.Name), node.Lines)
	} else {
		fmt.Printf("%-s%*d\n", prefix+"└── "+node.Name, p.maxWidth-runewidth.StringWidth(prefix+"└── "+node.Name), node.Lines)
	}

	for id, childNode := range node.Children {
		if id == len(node.Children)-1 {
			if isLast {
				p.PrintTree(childNode, prefix+"    ", true, false)
			} else {
				p.PrintTree(childNode, prefix+"│   ", true, false)
			}
		} else {
			if isLast {
				p.PrintTree(childNode, prefix+"    ", false, false)
			} else {
				p.PrintTree(childNode, prefix+"│   ", false, false)
			}
		}
	}
}

func (p *Printer) calSetMaxWidth(node *tree.Node, prefixLen int) {
	p.maxWidth = max(p.maxWidth, prefixLen+runewidth.StringWidth(node.Name+" "+strconv.Itoa(node.Lines))+4)
	for _, childNode := range node.Children {
		p.calSetMaxWidth(childNode, prefixLen+runewidth.StringWidth("├── "))
	}
}
