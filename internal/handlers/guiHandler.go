package handlers

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
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
		slog.Error(
			"Preview Generator failed to initialize, preview images will not be generated",
			slog.String("error", err.Error()))
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
	warnings := []string{}
	for _, issue := range validators.ValidateEditorState(&stateDto) {
		issue.Fix(&stateDto)
		warnings = append(warnings, issue.Message)
	}

	configuration := this.mapper.FromEditorState(stateDto)
	if configuration.TemplateName == "" {
		return dtos.TemplateLoadDto{}, common_errors.ErrNoTemplateName
	}

	this.templateGenerator.SetConfiguration(configuration)
	template := this.templateGenerator.Generate()
	if template == nil {
		return dtos.TemplateLoadDto{}, common_errors.ErrGeneratedTemplateInvalid
	}

	return dtos.TemplateLoadDto{Template: template, Warnings: warnings}, nil
}

func (this *GUIHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	if templateDto.Template == nil || len(templateDto.Template.Variants) == 0 {
		return dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid
	}

	if templateDto.Config.PlayerCount < 2 || templateDto.Config.PlayerCount > 8 {
		return dtos.TemplateLoadDto{}, common_errors.ErrInvalidPlayerCount
	}

	newTemplate := *templateDto.Template
	newTemplate.Variants = slices.Clone(templateDto.Template.Variants)
	newTemplate.Variants[0].Zones = templateDto.Zones
	newTemplate.Variants[0].Connections = templateDto.Connections

	connection_editor.RebuildZoneConnectionRoads(
		newTemplate.Variants[0].Zones,
		newTemplate.Variants[0].Connections)

	// Rebuild mandatory content from the final zones so a zone re-tiered in the
	// manual editor (e.g. Medium -> High) gets the content of its new quality
	// instead of keeping the content keyed to its original generation tier.
	if templateDto.Config != nil {
		newTemplate.MandatoryContent = this.contentProvider.CreateContentsForZones(
			*templateDto.Config, newTemplate.Variants[0].Zones)
	}

	var err error
	if connection_editor.ComputeHasErrors(newTemplate.Variants[0].Zones, newTemplate.Variants[0].Connections) {
		err = common_errors.ErrZonesMissing
	}

	return dtos.TemplateLoadDto{Template: &newTemplate}, err
}

func (this *GUIHandler) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	if templateDto.Template == nil {
		return "", common_errors.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(templateDto.OutputPath)
	if outputPath == "" {
		return "", common_errors.ErrNoOutputPath
	}

	out, err := this.fileService.SaveTemplate(outputPath, templateDto.Template)
	if err != nil {
		return "", err
	}

	if this.previewGenerator != nil {
		previewImage := this.previewGenerator.CreatePreviewImage(templateDto.Template, templateDto.Topology)
		_, err = this.fileService.SavePreviewImage(outputPath, previewImage, templateDto.Template.Name)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

// LoadState reads an editor state from the given .gen.json path and
// validates it against the editor's allowed values. When fixIssues is true,
// every detected issue is corrected in the returned state; the returned
// warnings describe the issues found either way.
func (this *GUIHandler) LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, nil, common_errors.ErrNoOutputPath
	}

	loaded, err := this.fileService.LoadSettingsFile(path)
	if err != nil {
		return nil, nil, err
	}

	issues := validators.ValidateEditorState(loaded)
	warnings := make([]string, 0, len(issues))
	for _, issue := range issues {
		if fixIssues {
			issue.Fix(loaded)
		}
		warnings = append(warnings, issue.Message)
	}

	return loaded, warnings, nil
}

func (this *GUIHandler) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	if stateDto.State == nil {
		return "", common_errors.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(stateDto.OutputPath)
	if outputPath == "" {
		return "", common_errors.ErrNoOutputPath
	}

	err := this.fileService.SaveSettings(outputPath, stateDto.State)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
