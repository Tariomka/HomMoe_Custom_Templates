package handlers

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
)

type GUIHandler struct {
	templateGenerator *template_generator.TemplateGenerator
	mapper            *mappers.GeneratorConfigMapper
	contentProvider   *providers.MandatoryContentProvider
	fileService       *file_service.FileService
	previewGenerator  *preview_service.PreviewGeneratorService
}

func NewGuiHandler() *GUIHandler {
	previewGenerator, err := preview_service.NewPreviewGenerator()
	if err != nil {
		fmt.Printf("Preview Generator failed to initialize, preview images will not be generated. Error: %v\n", err)
	}

	return &GUIHandler{
		templateGenerator: template_generator.NewTemplateGenerator(nil),
		mapper:            mappers.NewConfigMapper(),
		contentProvider:   providers.NewMandatoryContentProvider(),
		fileService:       file_service.NewFileService(),
		previewGenerator:  previewGenerator,
	}
}

func (this *GUIHandler) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error) {
	configuration := this.mapper.FromEditorState(stateDto)
	if configuration.TemplateName == "" {
		return dtos.TemplateLoadDto{}, common.ErrNoTemplateName
	}

	this.templateGenerator.SetConfiguration(configuration)
	template := this.templateGenerator.Generate()
	if template == nil {
		return dtos.TemplateLoadDto{}, common.ErrGeneratedTemplateInvalid
	}

	return dtos.TemplateLoadDto{Template: template}, nil
}

func (this *GUIHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	if templateDto.Template == nil || len(templateDto.Template.Variants) == 0 {
		return dtos.TemplateLoadDto{}, common.ErrProvidedTemplateInvalid
	}

	templateDto.Template.Variants[0].Zones = templateDto.Zones
	templateDto.Template.Variants[0].Connections = templateDto.Connections

	connection_editor.RebuildZoneConnectionRoads(
		templateDto.Template.Variants[0].Zones,
		templateDto.Template.Variants[0].Connections)

	// Rebuild mandatory content from the final zones so a zone re-tiered in the
	// manual editor (e.g. Medium -> High) gets the content of its new quality
	// instead of keeping the content keyed to its original generation tier.
	if templateDto.Config != nil {
		templateDto.Template.MandatoryContent = this.contentProvider.CreateContentsForZones(
			*templateDto.Config, templateDto.Template.Variants[0].Zones)
	}

	var err error
	if connection_editor.ComputeHasErrors(templateDto.Zones, templateDto.Connections) {
		err = common.ErrZonesMissing
	}

	return dtos.TemplateLoadDto{Template: templateDto.Template}, err
}

func (this *GUIHandler) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	if templateDto.Template == nil {
		return "", common.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(templateDto.OutputPath)
	if outputPath == "" {
		return "", common.ErrNoOutputPath
	}

	out, err := this.fileService.SaveTemplate(outputPath, templateDto.Template)
	if err != nil {
		return "", err
	}

	if this.previewGenerator != nil {
		// previewImage := previewGenerator.CreatePreviewImage(templateDto.Template, templateDto.Topology)
		// _, err = this.fileService.SavePreviewImage(outputPath, previewImage, templateDto.Template.Name)
		_, err = services.WritePreviewPNG(outputPath, templateDto.Template, templateDto.Topology)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func (this *GUIHandler) LoadState(path string) (*dtos.EditorStateDto, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, common.ErrNoOutputPath
	}

	loaded, err := this.fileService.LoadSettingsFile(path)
	if err != nil {
		return nil, err
	}

	return loaded, nil
}

func (this *GUIHandler) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	if stateDto.State == nil {
		return "", common.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(stateDto.OutputPath)
	if outputPath == "" {
		return "", common.ErrNoOutputPath
	}

	err := this.fileService.SaveSettings(outputPath, stateDto.State)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
