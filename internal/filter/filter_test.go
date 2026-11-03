package filter

import "testing"

func mustCompile(t *testing.T, include, exclude []string, ignoreCase bool) *Matcher {
	t.Helper()
	m, err := Compile(include, exclude, ignoreCase)
	if err != nil {
		t.Fatalf("Compile(%v, %v, %v) error: %v", include, exclude, ignoreCase, err)
	}
	return m
}

func TestMatchPath(t *testing.T) {
	tests := []struct {
		name    string
		include []string
		exclude []string
		ic      bool
		rel     string
		isDir   bool
		keep    bool
		prune   bool
	}{
		{name: "empty matcher keeps file", rel: "a/b.go", keep: true},
		{name: "empty matcher keeps dir", rel: "a", isDir: true, keep: true},
		{name: "include matches file", include: []string{`\.go$`}, rel: "a/b.go", keep: true},
		{name: "include skips file", include: []string{`\.go$`}, rel: "a/b.ts", keep: false},
		{name: "include ignored for dir", include: []string{`\.go$`}, rel: "a/b", isDir: true, keep: true},
		{name: "exclude file", exclude: []string{`_test\.go$`}, rel: "a/b_test.go", keep: false},
		{name: "exclude dir prunes", exclude: []string{"internal"}, rel: "internal", isDir: true, keep: false, prune: true},
		{name: "exclude wins over include", include: []string{`\.go$`}, exclude: []string{`_test`}, rel: "a/b_test.go", keep: false},
		{name: "ignore case on", include: []string{"readme"}, ic: true, rel: "README.md", keep: true},
		{name: "ignore case off", include: []string{"readme"}, ic: false, rel: "README.md", keep: false},
		{name: "nested anchored include", include: []string{`^internal/scanner/.*\.go$`}, rel: "internal/scanner/scanner.go", keep: true},
		{name: "nested anchored include skips other dir", include: []string{`^internal/scanner/.*\.go$`}, rel: "internal/tree/node.go", keep: false},
		{name: "root always kept", exclude: []string{"."}, rel: ".", isDir: true, keep: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mustCompile(t, tt.include, tt.exclude, tt.ic)
			keep, prune := m.MatchPath(tt.rel, tt.isDir)
			if keep != tt.keep || prune != tt.prune {
				t.Fatalf("MatchPath(%q, dir=%v) = keep %v, prune %v; want keep %v, prune %v",
					tt.rel, tt.isDir, keep, prune, tt.keep, tt.prune)
			}
		})
	}
}

func TestCompileInvalid(t *testing.T) {
	if _, err := Compile([]string{"["}, nil, false); err == nil {
		t.Fatal("expected error for invalid include regex")
	}
	if _, err := Compile(nil, []string{"("}, false); err == nil {
		t.Fatal("expected error for invalid exclude regex")
	}
}

func TestIsEmpty(t *testing.T) {
	if !mustCompile(t, nil, nil, false).IsEmpty() {
		t.Fatal("expected empty matcher")
	}
	if mustCompile(t, []string{"a"}, nil, false).IsEmpty() {
		t.Fatal("expected non-empty matcher")
	}
}

func TestNilMatcher(t *testing.T) {
	var m *Matcher
	if keep, prune := m.MatchPath("a/b.go", false); !keep || prune {
		t.Fatal("nil matcher should keep everything and prune nothing")
	}
}
