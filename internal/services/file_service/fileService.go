package file_service

import (
	"image"
	"path/filepath"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
)

type FileService struct {
	editorStateRepository repositories.IFileRepository[editor_state.EditorState]
	templateRepository    repositories.IFileRepository[template.RmgTemplate]
	previewRepository     repositories.IFileRepository[image.RGBA]
	editorStateMapper     mappers.IEditorStateMapper
	templateMapper        mappers.ITemplateMapper
}

func NewFileService(
	editorStateRepository repositories.IFileRepository[editor_state.EditorState],
	templateRepository repositories.IFileRepository[template.RmgTemplate],
	previewRepository repositories.IFileRepository[image.RGBA],
	editorStateMapper mappers.IEditorStateMapper,
	templateMapper mappers.ITemplateMapper) IFileService {
	return &FileService{
		editorStateRepository: editorStateRepository,
		templateRepository:    templateRepository,
		previewRepository:     previewRepository,
		editorStateMapper:     editorStateMapper,
		templateMapper:        templateMapper,
	}
}

func (this *FileService) LoadSettingsFile(filePath string) (*editor_state_model.EditorState, error) {
	entity := this.editorStateMapper.NewDefaultEntity()
	if err := this.editorStateRepository.Load(filePath, &entity); err != nil {
		return nil, err
	}

	return new(this.editorStateMapper.ToModel(entity)), nil
}

func (this *FileService) SaveSettings(
	filePath string,
	editorState *editor_state_model.EditorState) (string, error) {
	entity := this.editorStateMapper.ToEntity(*editorState)
	return this.editorStateRepository.Save(filepath.Dir(filePath), editorState.TemplateName, entity)
}

func (this *FileService) SaveTemplateWithPreview(
	directory string,
	template *template_model.Template,
	previewImage *image.RGBA) (string, error) {
	entity := this.templateMapper.ToEntity(*template)
	templatePath, err := this.templateRepository.Save(directory, entity.Name, entity)
	if err != nil {
		return "", err
	}

	if previewImage == nil {
		return templatePath, nil
	}

	if _, err = this.previewRepository.Save(directory, entity.Name, *previewImage); err != nil {
		return templatePath, err
	}

	return templatePath, nil
}
