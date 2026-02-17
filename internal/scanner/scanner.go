package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/zihao-liu-qs/treekei/internal/tree"
)

type Scanner struct {
	ignoreNames map[string]bool
	sortBy      SortMode
	allowedExts map[string]bool
}

type SortMode string

const (
	SortByName  SortMode = "name"
	SortByLines SortMode = "lines"
)

// langExtensions maps language names to their file extensions
var langExtensions = map[string][]string{
	"go":         {".go"},
	"python":     {".py"},
	"py":         {".py"},
	"javascript": {".js", ".mjs", ".cjs"},
	"js":         {".js", ".mjs", ".cjs"},
	"typescript": {".ts", ".tsx"},
	"ts":         {".ts", ".tsx"},
	"rust":       {".rs"},
	"rs":         {".rs"},
	"java":       {".java"},
	"c":          {".c", ".h"},
	"cpp":        {".cpp", ".cc", ".cxx", ".hpp"},
	"ruby":       {".rb"},
	"rb":         {".rb"},
	"swift":      {".swift"},
	"kotlin":     {".kt", ".kts"},
	"kt":         {".kt", ".kts"},
	"css":        {".css", ".scss", ".sass", ".less"},
	"html":       {".html", ".htm"},
	"shell":      {".sh", ".bash", ".zsh"},
	"sh":         {".sh", ".bash", ".zsh"},
}

// ParseLangs converts a comma-separated language string into a set of allowed extensions.
// Returns nil if langs is empty (no filter).
func ParseLangs(langs string) map[string]bool {
	if langs == "" {
		return nil
	}
	allowed := make(map[string]bool)
	for _, lang := range strings.Split(langs, ",") {
		lang = strings.TrimSpace(strings.ToLower(lang))
		if exts, ok := langExtensions[lang]; ok {
			for _, ext := range exts {
				allowed[ext] = true
			}
		}
	}
	return allowed
}

func NewScanner(showAll bool, sortBy SortMode, langs string) *Scanner {
	ignoreNames := make(map[string]bool)

	if !showAll {
		ignoreNames[".git"] = true
		ignoreNames[".env"] = true
		ignoreNames[".DS_Store"] = true
	}

	return &Scanner{
		ignoreNames: ignoreNames,
		sortBy:      sortBy,
		allowedExts: ParseLangs(langs),
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
		lines, isBinary, err := countLines(path)
		if err != nil {
			node.Lines = 0
		} else {
			node.Lines = lines
		}
		if isBinary {
			node.IsBinary = true
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

		// Skip files that don't match the language filter
		if !childNode.IsDir && s.shouldSkipByLang(childNode.Name) {
			continue
		}

		if childNode != nil {
			node.Children = append(node.Children, childNode)
			totalLines += childNode.Lines
		}
	}

	switch s.sortBy {
	case SortByName:
		sort.Slice(node.Children, func(i, j int) bool {
			return node.Children[i].Name < node.Children[j].Name
		})
	default:
		sort.Slice(node.Children, func(i, j int) bool {
			return node.Children[i].Lines > node.Children[j].Lines
		})
	}

	node.Lines = totalLines
	return node, errors
}

func countLines(path string) (lines int, isBinary bool, err error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, false, err
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && n == 0 {
		return 0, false, err
	}

	// Check if the content is valid UTF-8 (simple text detection)
	if !utf8.Valid(buf[:n]) {
		return 0, true, nil
	}

	// Reset file pointer to beginning
	file.Seek(0, 0)

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, false, err
	}

	return lineCount, false, nil
}

func (s *Scanner) shouldIgnore(name string) bool {
	val, _ := s.ignoreNames[name]
	return val
}

// shouldSkipByLang returns true if the file should be excluded based on language filter.
func (s *Scanner) shouldSkipByLang(name string) bool {
	if s.allowedExts == nil {
		return false
	}
	ext := strings.ToLower(filepath.Ext(name))
	return !s.allowedExts[ext]
}
