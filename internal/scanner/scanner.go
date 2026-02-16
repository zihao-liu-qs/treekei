package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"unicode/utf8"

	"github.com/zihao-liu-qs/treekei/internal/tree"
)

type Scanner struct {
	ignoreNames map[string]bool
}

func NewScanner(showAll bool) *Scanner {
	ignoreNames := make(map[string]bool)

	if !showAll {
		ignoreNames[".git"] = true
		ignoreNames[".env"] = true
		ignoreNames[".DS_Store"] = true
	}

	return &Scanner{
		ignoreNames: ignoreNames,
	}
}

func (s *Scanner) ScanDirectory(path string) (*tree.Node, []error) {
	var errors []error

	info, err := os.Stat(path)
	if err != nil {
		return nil, []error{fmt.Errorf("cannot access %s: %v", path, err)}
	}

	node := &tree.Node{
		Name:  filepath.Base(path),
		Path:  path,
		IsDir: info.IsDir(),
	}

	if !info.IsDir() {
		lines, err := countLines(path)
		if err != nil {
			node.Lines = 0
		} else {
			node.Lines = lines
		}
		return node, errors
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return node, errors
	}

	totalLines := 0
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		childNode, childErrors := s.ScanDirectory(childPath)
		errors = append(errors, childErrors...)

		if s.shouldIgnore(childNode.Name) {
			continue
		}

		if childNode != nil {
			node.Children = append(node.Children, childNode)
			totalLines += childNode.Lines
		}
	}

	// sort.Slice(node.Children, func(i, j int) bool {
	// 	return node.Children[i].Name < node.Children[j].Name
	// })
	sort.Slice(node.Children, func(i, j int) bool {
		return node.Children[i].Lines > node.Children[j].Lines
	})

	node.Lines = totalLines
	return node, errors
}

func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && n == 0 {
		return 0, err
	}

	// Check if the content is valid UTF-8 (simple text detection)
	if !utf8.Valid(buf[:n]) {
		return 0, nil
	}

	// Reset file pointer to beginning
	file.Seek(0, 0)

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return lineCount, nil
}

func (s *Scanner) shouldIgnore(name string) bool {
	val, _ := s.ignoreNames[name]
	return val
}
