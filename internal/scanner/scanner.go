package scanner

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/zihao-liu-qs/treekei/internal/filter"
	"github.com/zihao-liu-qs/treekei/internal/tree"
)

type Scanner struct {
	ignoreNames map[string]bool
	sortBy      SortMode
	allowedExts map[string]bool
	matcher     *filter.Matcher
	root        string
}

type SortMode string

const (
	SortByName  SortMode = "name"
	SortByLines SortMode = "lines"
)

// Options configures a Scanner.
type Options struct {
	ShowAll bool
	SortBy  SortMode
	Langs   string
	Matcher *filter.Matcher
}

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
	"csharp":     {".cs"},
	"cs":         {".cs"},
	"ruby":       {".rb"},
	"rb":         {".rb"},
	"swift":      {".swift"},
	"kotlin":     {".kt", ".kts"},
	"kt":         {".kt", ".kts"},
	"php":        {".php"},
	"scala":      {".scala"},
	"dart":       {".dart"},
	"lua":        {".lua"},
	"perl":       {".pl", ".pm"},
	"r":          {".r"},
	"css":        {".css", ".scss", ".sass", ".less"},
	"html":       {".html", ".htm"},
	"vue":        {".vue"},
	"svelte":     {".svelte"},
	"shell":      {".sh", ".bash", ".zsh"},
	"sh":         {".sh", ".bash", ".zsh"},
	"yaml":       {".yaml", ".yml"},
	"yml":        {".yaml", ".yml"},
	"json":       {".json"},
	"toml":       {".toml"},
	"xml":        {".xml"},
	"sql":        {".sql"},
	"proto":      {".proto"},
	"markdown":   {".md", ".markdown"},
	"md":         {".md", ".markdown"},
}

// SupportedLangs returns the sorted list of language names accepted by the
// --lang filter.
func SupportedLangs() []string {
	names := make([]string, 0, len(langExtensions))
	for name := range langExtensions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
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

func NewScanner(opts Options) *Scanner {
	ignoreNames := make(map[string]bool)

	if !opts.ShowAll {
		ignoreNames[".git"] = true
		ignoreNames[".env"] = true
		ignoreNames[".DS_Store"] = true
	}

	return &Scanner{
		ignoreNames: ignoreNames,
		sortBy:      opts.SortBy,
		allowedExts: ParseLangs(opts.Langs),
		matcher:     opts.Matcher,
	}
}

func (s *Scanner) ScanDirectory(path string) (*tree.Node, []error) {
	s.root = path
	return s.scan(path)
}

func (s *Scanner) scan(path string) (*tree.Node, []error) {
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
		rel := relPath(s.root, childPath)

		// A directory matched by an exclude pattern is pruned: no traversal,
		// no line counting.
		if entry.IsDir() {
			if _, prune := s.matcher.MatchPath(rel, true); prune {
				continue
			}
		}

		childNode, childErrors := s.scan(childPath)
		errors = append(errors, childErrors...)

		if childNode == nil || s.shouldIgnore(childNode.Name) {
			continue
		}

		if !childNode.IsDir {
			// Skip files that don't match the language filter or the path filter
			if s.shouldSkipByLang(childNode.Name) {
				continue
			}
			if keep, _ := s.matcher.MatchPath(rel, false); !keep {
				continue
			}
		}

		node.Children = append(node.Children, childNode)
		totalLines += childNode.Lines
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

// relPath returns path relative to root using forward slashes, matching the
// form that filter patterns are written against.
func relPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
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

	// Sample the first bytes to detect binary, but strip a trailing
	// multi-byte rune that may be cut off by the read boundary, so a
	// valid UTF-8 character is not mistaken for binary data.
	sample := buf[:n]
	end := len(sample)
	for end > 0 && !utf8.RuneStart(sample[end-1]) { // skip trailing continuation bytes
		end--
	}
	if end > 0 {
		lead := sample[end-1]
		cont := len(sample) - end // continuation bytes after the lead
		need := 0                 // continuation bytes this lead expects
		switch {
		case lead < utf8.RuneSelf: // ASCII, 1 byte
			need = 0
		case lead < 0xE0: // 2-byte
			need = 1
		case lead < 0xF0: // 3-byte
			need = 2
		default: // 4-byte
			need = 3
		}
		if cont < need { // rune was truncated; drop the lead byte too
			sample = sample[:end-1]
		}
	}

	if !utf8.Valid(sample) {
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
