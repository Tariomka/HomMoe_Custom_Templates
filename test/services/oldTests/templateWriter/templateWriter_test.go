package templateWriter_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// ── SanitizeFilename ─────────────────────────────────────────────────

func TestSanitizeFilename_PlainNameUnchanged(t *testing.T) {
	if got := helpers.SanitizeFilename("clean_name-1"); got != "clean_name-1" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_TrimsWhitespace(t *testing.T) {
	if got := helpers.SanitizeFilename("   spaced   "); got != "spaced" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesForwardSlash(t *testing.T) {
	if got := helpers.SanitizeFilename("a/b"); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesBackslash(t *testing.T) {
	if got := helpers.SanitizeFilename(`a\b`); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesColon(t *testing.T) {
	if got := helpers.SanitizeFilename("a:b"); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesAsterisk(t *testing.T) {
	if got := helpers.SanitizeFilename("a*b"); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesQuestionMark(t *testing.T) {
	if got := helpers.SanitizeFilename("a?b"); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesDoubleQuote(t *testing.T) {
	if got := helpers.SanitizeFilename(`a"b`); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesLessThan(t *testing.T) {
	if got := helpers.SanitizeFilename("a<b"); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesGreaterThan(t *testing.T) {
	if got := helpers.SanitizeFilename("a>b"); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesPipe(t *testing.T) {
	if got := helpers.SanitizeFilename("a|b"); got != "a_b" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_ReplacesMultipleBadRunesInOneCall(t *testing.T) {
	if got := helpers.SanitizeFilename(`bad/\*?":<>|name`); got != "bad_________name" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_EmptyString(t *testing.T) {
	if got := helpers.SanitizeFilename(""); got != "" {
		t.Errorf("got %q", got)
	}
}

func TestSanitizeFilename_OnlyWhitespace(t *testing.T) {
	if got := helpers.SanitizeFilename("   "); got != "" {
		t.Errorf("got %q", got)
	}
}

// ── WriteTemplate ────────────────────────────────────────────────────

func TestWriteTemplate_CreatesFileWithSanitisedName(t *testing.T) {
	dir := t.TempDir()
	tmpl := &template.RmgTemplate{Name: "My/Template"}
	path, err := services.WriteTemplate(dir, tmpl)
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "My_Template.rmg.json")
	if path != want {
		t.Errorf("path = %q, want %q", path, want)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("file not written: %v", err)
	}
}

func TestWriteTemplate_FallsBackToGeneratedTemplateOnEmptyName(t *testing.T) {
	dir := t.TempDir()
	tmpl := &template.RmgTemplate{Name: ""}
	path, err := services.WriteTemplate(dir, tmpl)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasSuffix(path, "Generated_Template.rmg.json") {
		t.Errorf("path = %q, want suffix Generated_Template.rmg.json", path)
	}
}

func TestWriteTemplate_FallsBackWhenNameIsAllInvalidThenEmpty(t *testing.T) {
	dir := t.TempDir()
	tmpl := &template.RmgTemplate{Name: "   "}
	path, err := services.WriteTemplate(dir, tmpl)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.Base(path) != "Generated_Template.rmg.json" {
		t.Errorf("path = %q", path)
	}
}

func TestWriteTemplate_CreatesNestedMissingDirectory(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "a", "b", "c")
	tmpl := &template.RmgTemplate{Name: "T"}
	if _, err := services.WriteTemplate(dir, tmpl); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("dir not created: %v", err)
	}
}

func TestWriteTemplate_ProducesIndentedJSON(t *testing.T) {
	dir := t.TempDir()
	tmpl := &template.RmgTemplate{Name: "T", SizeX: 10}
	path, err := services.WriteTemplate(dir, tmpl)
	if err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "\n  ") {
		t.Error("expected indented JSON")
	}
	var round template.RmgTemplate
	if err := json.Unmarshal(data, &round); err != nil {
		t.Errorf("round-trip parse failed: %v", err)
	}
}

func TestWriteTemplate_MkdirError(t *testing.T) {
	// Create a file then try to use it as a directory parent.
	parent := filepath.Join(t.TempDir(), "blocker")
	if err := os.WriteFile(parent, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(parent, "child")
	tmpl := &template.RmgTemplate{Name: "T"}
	if _, err := services.WriteTemplate(dir, tmpl); err == nil {
		t.Error("expected mkdir error")
	}
}
