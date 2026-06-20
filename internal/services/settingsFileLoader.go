package services

import (
	"encoding/json"
	"os"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
)

// LoadSettingsFile reads a .gen.json file and returns the parsed SettingsFile.
func LoadSettingsFile(path string) (*dtos.EditorStateDto, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	settingsFile := dtos.NewDefaultEditorStateDto()
	if err := json.Unmarshal(data, &settingsFile); err != nil {
		return nil, err
	}

	return &settingsFile, nil
}

// SaveSettingsFile writes a SettingsFile to disk as indented JSON.
func SaveSettingsFile(path string, editorState *dtos.EditorStateDto) error {
	data, err := json.MarshalIndent(editorState, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
