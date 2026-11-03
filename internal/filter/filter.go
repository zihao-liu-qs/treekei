package filter

import (
	"fmt"
	"regexp"
	"strings"
)

// Matcher holds compiled include/exclude regexes applied to paths relative to
// the scan root. A nil *Matcher keeps everything and prunes nothing.
type Matcher struct {
	include []*regexp.Regexp
	exclude []*regexp.Regexp
}

// Compile builds a Matcher from raw patterns. Patterns are matched as
// unanchored regular expressions (Go's RE2, so no backreferences or
// lookaround). When ignoreCase is set, every pattern is made case-insensitive.
func Compile(include, exclude []string, ignoreCase bool) (*Matcher, error) {
	inc, err := compileAll(include, ignoreCase)
	if err != nil {
		return nil, err
	}
	exc, err := compileAll(exclude, ignoreCase)
	if err != nil {
		return nil, err
	}
	return &Matcher{include: inc, exclude: exc}, nil
}

func compileAll(patterns []string, ignoreCase bool) ([]*regexp.Regexp, error) {
	var out []*regexp.Regexp
	for _, p := range patterns {
		pattern := p
		if ignoreCase {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("%q: %w", p, err)
		}
		out = append(out, re)
	}
	return out, nil
}

// IsEmpty reports whether no patterns are set.
func (m *Matcher) IsEmpty() bool {
	return len(m.include) == 0 && len(m.exclude) == 0
}

// MatchPath classifies a forward-slash normalized path relative to the scan
// root.
//
// For directories, include patterns are ignored; keep is false and prune is
// true when the directory matches an exclude pattern.
//
// For files, keep is true when the file matches an include pattern (or no
// include patterns are set) and does not match any exclude pattern.
func (m *Matcher) MatchPath(rel string, isDir bool) (keep bool, prune bool) {
	if m == nil {
		return true, false
	}

	rel = strings.TrimPrefix(rel, "./")
	if rel == "" || rel == "." {
		return true, false
	}

	if matchesAny(m.exclude, rel) {
		if isDir {
			return false, true
		}
		return false, false
	}

	if isDir {
		return true, false
	}

	if len(m.include) == 0 {
		return true, false
	}

	return matchesAny(m.include, rel), false
}

func matchesAny(res []*regexp.Regexp, s string) bool {
	for _, re := range res {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}
