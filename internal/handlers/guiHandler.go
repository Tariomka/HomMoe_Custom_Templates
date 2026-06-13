package handlers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
)

type GUIHandler struct {
	templateGenerator *template_generator.TemplateGenerator
}

func NewGuiHandler() *GUIHandler {
	return &GUIHandler{
		templateGenerator: template_generator.NewTemplateGenerator(nil),
	}
}

func (this *GUIHandler) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error) {
	configuration := services.SettingsToGenerator(&stateDto)
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

	out, err := services.WriteTemplate(outputPath, templateDto.Template)
	if err != nil {
		return "", err
	}

	_, err = services.WritePreviewPNG(outputPath, templateDto.Template, templateDto.Topology)
	if err != nil {
		return out, err
	}

	return out, nil
}

func (this *GUIHandler) LoadState(path string) (*dtos.EditorStateDto, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, common.ErrNoOutputPath
	}

	loaded, err := services.LoadSettingsFile(path)
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

	err := services.SaveSettingsFile(outputPath, stateDto.State)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}
