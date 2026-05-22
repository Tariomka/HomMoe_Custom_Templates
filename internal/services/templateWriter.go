package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

// SanitizeFilename replaces filesystem-unsafe runes in name with underscores
// and trims surrounding whitespace.
func SanitizeFilename(name string) string {
	bad := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	out := []rune(strings.TrimSpace(name))
	for index, runeValue := range out {
		for _, badRune := range bad {
			if runeValue == badRune {
				out[index] = '_'
			}
		}
	}
	return string(out)
}

// WriteTemplate writes the given template as indented JSON into
// dir/<safeName>.rmg.json. The directory is created if missing. Returns the
// final path on success.
func WriteTemplate(dir string, template *template.RmgTemplateModel) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	safeName := SanitizeFilename(template.Name)
	if safeName == "" {
		safeName = "Generated_Template"
	}
	out := filepath.Join(dir, safeName+".rmg.json")
	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return "", err
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		return "", err
	}
	return out, nil
}
