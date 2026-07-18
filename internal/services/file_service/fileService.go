package file_service

import (
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

const (
	folderPermission        = 0o755
	fileReadWritePermission = 0o644

	pngExtension      = ".png"
	templateExtension = ".rmg.json"
	defaultName       = "Generated_Template"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

// LoadSettingsFile reads settings file from the given filepath and returns the parsed settings object.
func (this *FileService) LoadSettingsFile(filePath string) (*dtos.EditorStateDto, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	settingsFile := dtos.NewDefaultEditorStateDto()
	if err = json.Unmarshal(data, &settingsFile); err != nil {
		return nil, err
	}

	return &settingsFile, nil
}

// SaveSettings writes settings object to the given filepath as JSON.
func (this *FileService) SaveSettings(filePath string, editorState *dtos.EditorStateDto) error {
	data, err := json.MarshalIndent(editorState, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, data, fileReadWritePermission)
}

// SaveTemplate writes the given template as JSON into `dir/<template_name>.rmg.json`.
//
// The directory is created if missing and returns the final path on success.
func (this *FileService) SaveTemplate(directory string, template *entities.RmgTemplate) (string, error) {
	if err := os.MkdirAll(directory, folderPermission); err != nil {
		return "", err
	}

	safeName := helpers.SanitizeFilename(template.Name)
	if safeName == "" {
		safeName = defaultName
	}
	out := filepath.Join(directory, safeName+templateExtension)
	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return "", err
	}

	if err = os.WriteFile(out, data, fileReadWritePermission); err != nil {
		return "", err
	}

	return out, nil
}

// SavePreviewImage writes the given image as PNG into `dir/<template_name>.png`.
//
// The directory is created if missing and returns the final path on success.
func (this *FileService) SavePreviewImage(
	directory string,
	previewImage *image.RGBA,
	templateName string) (string, error) {
	if err := os.MkdirAll(directory, folderPermission); err != nil {
		return "", err
	}

	safeName := helpers.SanitizeFilename(templateName)
	if safeName == "" {
		safeName = defaultName
	}
	out := filepath.Join(directory, safeName+pngExtension)
	file, err := os.Create(out)
	if err != nil {
		return "", err
	}

	defer file.Close()
	if err = png.Encode(file, previewImage); err != nil {
		return "", err
	}

	return out, nil
}
