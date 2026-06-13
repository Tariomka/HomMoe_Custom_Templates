package services

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

// WriteTemplate writes the given template as indented JSON into
// dir/<safeName>.rmg.json. The directory is created if missing. Returns the
// final path on success.
func WriteTemplate(dir string, template *template.RmgTemplate) (string, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	safeName := helpers.SanitizeFilename(template.Name)
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
