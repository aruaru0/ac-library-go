package main

import (
	"os"
	"path/filepath"
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
	dir := t.TempDir()
	licenseText := "MIT License\n\nCopyright (c) 2026 aruaru0\n\n" +
		"Permission is hereby granted, free of charge, to any person obtaining a copy\n" +
		"of this software and associated documentation files (the \"Software\"), to deal\n" +
		"in the Software without restriction, including without limitation the rights\n" +
		"to use, copy, modify, merge, publish, distribute, sublicense, and/or sell\n" +
		"copies of the Software.\n"
	if err := os.WriteFile(filepath.Join(dir, "LICENSE"), []byte(licenseText), 0o644); err != nil {
		t.Fatal(err)
	}

	mod := Module{
		Path:    "github.com/aruaru0/ac-library-go",
		Version: "v1.2.3",
		Dir:     dir,
	}
	banner := buildBanner(mod)

	if !strings.Contains(banner, "https://github.com/aruaru0/ac-library-go/blob/v1.2.3/cmd/acl-go-expander/main.go") {
		t.Errorf("banner does not contain tool permalink: %q", banner)
	}
	if !strings.Contains(banner, "https://github.com/aruaru0/ac-library-go/blob/v1.2.3/LICENSE") {
		t.Errorf("banner does not contain license permalink: %q", banner)
	}
	if !strings.Contains(banner, "// SPDX-License-Identifier: MIT") {
		t.Errorf("banner does not contain SPDX identifier: %q", banner)
	}
	if !strings.Contains(banner, "// Copyright (c) 2026 aruaru0") {
		t.Errorf("banner does not contain the copyright notice: %q", banner)
	}
	if !strings.Contains(banner, "Embedded code: ac-library-go (MIT License)") {
		t.Errorf("banner does not identify the license: %q", banner)
	}
	if strings.Contains(banner, "Permission is hereby granted") {
		t.Errorf("banner should not embed the full license text: %q", banner)
	}
	if !strings.HasPrefix(banner, "// ") {
		t.Errorf("banner should be a line comment block: %q", banner)
	}
}

func TestBuildBannerWithoutLicense(t *testing.T) {
	mod := Module{
		Path:    "github.com/aruaru0/ac-library-go",
		Version: "v1.2.3",
		Dir:     t.TempDir(),
	}
	banner := buildBanner(mod)
	if strings.Contains(banner, "Embedded code license") || strings.Contains(banner, "SPDX-License-Identifier") {
		t.Errorf("banner should not contain a license block when no license file exists: %q", banner)
	}
}

func TestBuildBannerUnknownLicenseFallsBackToFullText(t *testing.T) {
	dir := t.TempDir()
	licenseText := "The Foo License\n\nCopyright (c) 2026 Someone\n\nAll rights reserved.\n"
	if err := os.WriteFile(filepath.Join(dir, "LICENSE"), []byte(licenseText), 0o644); err != nil {
		t.Fatal(err)
	}

	mod := Module{Path: "github.com/aruaru0/ac-library-go", Version: "v1.2.3", Dir: dir}
	banner := buildBanner(mod)

	if !strings.Contains(banner, "Embedded code license: https://github.com/aruaru0/ac-library-go/blob/v1.2.3/LICENSE") {
		t.Errorf("banner should link the license for an unknown license: %q", banner)
	}
	if !strings.Contains(banner, "// The Foo License") {
		t.Errorf("banner should embed the full text for an unknown license: %q", banner)
	}
}

func TestSPDXIdentifier(t *testing.T) {
	cases := []struct {
		name string
		text string
		want string
		ok   bool
	}{
		{"mit title", "MIT License\n\nCopyright (c) 2026 x\n", "MIT", true},
		{"mit permission clause", "Permission is hereby granted, free of charge, to any person\n", "MIT", true},
		{"unknown", "The Foo License\nAll rights reserved.\n", "", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, ok := spdxIdentifier([]byte(c.text))
			if got != c.want || ok != c.ok {
				t.Errorf("spdxIdentifier(%q) = (%q, %v), want (%q, %v)", c.text, got, ok, c.want, c.ok)
			}
		})
	}
}

func TestCopyrightNotice(t *testing.T) {
	data := []byte("MIT License\n\n   Copyright (c) 2026 aruaru0\n\nPermission is hereby granted...\n")
	if got, want := copyrightNotice(data), "Copyright (c) 2026 aruaru0"; got != want {
		t.Errorf("copyrightNotice() = %q, want %q", got, want)
	}
}
