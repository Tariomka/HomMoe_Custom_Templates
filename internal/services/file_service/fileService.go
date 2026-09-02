package file_service

import (
	"image"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
)

// FileService decides where and under what name persisted files go; the
// repositories own the encoding, the extension and the atomic replacement.
type FileService struct {
	editorStateRepository repositories.IFileRepository[editor_state.EditorState]
	templateRepository    repositories.IFileRepository[template.RmgTemplate]
	previewRepository     repositories.IFileRepository[image.RGBA]
	editorStateMapper     mappers.IEditorStateMapper
}

func NewFileService(
	editorStateRepository repositories.IFileRepository[editor_state.EditorState],
	templateRepository repositories.IFileRepository[template.RmgTemplate],
	previewRepository repositories.IFileRepository[image.RGBA],
	editorStateMapper mappers.IEditorStateMapper) IFileService {
	return &FileService{
		editorStateRepository: editorStateRepository,
		templateRepository:    templateRepository,
		previewRepository:     previewRepository,
		editorStateMapper:     editorStateMapper,
	}
}

// LoadSettingsFile reads settings file from the given filepath and returns the parsed settings object.
func (this *FileService) LoadSettingsFile(filePath string) (*editor_state_model.EditorState, error) {
	entity := this.editorStateMapper.NewDefaultEntity()
	if err := this.editorStateRepository.Load(filePath, &entity); err != nil {
		return nil, err
	}

	return new(this.editorStateMapper.ToModel(entity)), nil
}

// SaveSettings writes the editor state next to filePath, named after the
// template, and returns the path actually written.
func (this *FileService) SaveSettings(
	filePath string,
	editorState *editor_state_model.EditorState) (string, error) {
	entity := this.editorStateMapper.ToEntity(*editorState)
	return this.editorStateRepository.Save(filepath.Dir(filePath), editorState.TemplateName, entity)
}

// SaveTemplateWithPreview writes the template and, when previewImage is not
// nil, its preview. A preview failure still returns the template path so the
// caller can report the partial success.
func (this *FileService) SaveTemplateWithPreview(
	directory string,
	template *template.RmgTemplate,
	previewImage *image.RGBA) (string, error) {
	templatePath, err := this.templateRepository.Save(directory, template.Name, *template)
	if err != nil {
		return "", err
	}

	if previewImage == nil {
		return templatePath, nil
	}

	if _, err = this.previewRepository.Save(directory, template.Name, *previewImage); err != nil {
		return templatePath, err
	}

	return templatePath, nil
}
