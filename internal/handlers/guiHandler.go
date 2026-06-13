package handlers

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
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

func (this *GUIHandler) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateDto, error) {
	configuration := services.SettingsToGenerator(&stateDto)
	if configuration.TemplateName == "" {
		return dtos.TemplateDto{}, common.ErrNoTemplateName
	}

	this.templateGenerator.SetConfiguration(configuration)
	template := this.templateGenerator.Generate()
	if template == nil {
		return dtos.TemplateDto{}, common.ErrGeneratedTemplateInvalid
	}

	return dtos.TemplateDto{Template: template}, nil
}

func (this *GUIHandler) SaveTemplate(templateDto dtos.TemplateDto) (string, error) {
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

	return out, nil
}

func (this *GUIHandler) SaveState(templateDto dtos.TemplateDto) (string, error) {
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

	return out, nil
}
