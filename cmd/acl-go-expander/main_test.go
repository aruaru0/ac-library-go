package main

import (
	"strings"
	"testing"
)

func TestPseudoVersionHash(t *testing.T) {
	cases := []struct {
		version string
		want    string
	}{
		{"v1.2.3", ""},
		{"v0.0.0-20240101120000-abcdef123456", "abcdef123456"},
		{"v1.0.0-0.20240101120000-0123456789ab", "0123456789ab"},
		{"v0.0.0-20240101120000-abcdef1234567", ""},
		{"v0.0.0-20240101120000-ABCDEF123456", ""},
	}
	for _, c := range cases {
		if got := pseudoVersionHash(c.version); got != c.want {
			t.Errorf("pseudoVersionHash(%q) = %q, want %q", c.version, got, c.want)
		}
	}
}

func TestModuleRev(t *testing.T) {
	cases := []struct {
		name string
		mod  Module
		want string
	}{
		{"tag", Module{Version: "v1.2.3"}, "v1.2.3"},
		{"pseudo version", Module{Version: "v0.0.0-20240101120000-abcdef123456"}, "abcdef123456"},
		{"replace version", Module{Version: "", Replace: &Module{Version: "v2.0.0"}}, "v2.0.0"},
		{"no version falls back to main", Module{Version: "", Dir: "/nonexistent/ac-library-go"}, "main"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := moduleRev(c.mod); got != c.want {
				t.Errorf("moduleRev(%+v) = %q, want %q", c.mod, got, c.want)
			}
		})
	}
}

func TestURLModulePath(t *testing.T) {
	cases := []struct {
		name string
		mod  Module
		want string
	}{
		{"plain", Module{Path: "github.com/aruaru0/ac-library-go"}, "github.com/aruaru0/ac-library-go"},
		{
			"fork replace",
			Module{Path: "github.com/aruaru0/ac-library-go", Replace: &Module{Path: "github.com/fork/ac-library-go"}},
			"github.com/fork/ac-library-go",
		},
		{
			"local replace",
			Module{Path: "github.com/aruaru0/ac-library-go", Replace: &Module{Path: "../ac-library-go"}},
			"github.com/aruaru0/ac-library-go",
		},
		{
			"absolute local replace",
			Module{Path: "github.com/aruaru0/ac-library-go", Replace: &Module{Path: "/tmp/ac-library-go"}},
			"github.com/aruaru0/ac-library-go",
		},
		{
			"empty path falls back to library repo",
			Module{},
			"github.com/aruaru0/ac-library-go",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := urlModulePath(c.mod); got != c.want {
				t.Errorf("urlModulePath(%+v) = %q, want %q", c.mod, got, c.want)
			}
		})
	}
}

func TestSourceURL(t *testing.T) {
	pkg := &GoPackage{
		Dir: "/repo/dsu",
		Module: Module{
			Path:    "github.com/aruaru0/ac-library-go",
			Version: "v1.2.3",
			Dir:     "/repo",
		},
	}
	got := sourceURL(pkg, "/repo/dsu/dsu.go")
	want := "https://github.com/aruaru0/ac-library-go/blob/v1.2.3/dsu/dsu.go"
	if got != want {
		t.Errorf("sourceURL() = %q, want %q", got, want)
	}
}

func TestBuildBanner(t *testing.T) {
	banner := buildBanner("https://github.com/aruaru0/ac-library-go", "v1.2.3")
	if !strings.Contains(banner, "https://github.com/aruaru0/ac-library-go/blob/v1.2.3/cmd/acl-go-expander/main.go") {
		t.Errorf("banner does not contain tool permalink: %q", banner)
	}
	if !strings.HasPrefix(banner, "// ") {
		t.Errorf("banner should be a line comment block: %q", banner)
	}
}
