package handlers

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
)

type templateWorkflowHandler struct {
	templateGenerator *template_generator.TemplateGenerator
	mapper            *mappers.GeneratorConfigMapper
	contentProvider   *providers.MandatoryContentProvider
	connectionEditor  *connection_editor.ConnectionEditorService
	zoneEditor        *connection_editor.ZoneEditorService
	manualReapply     *connection_editor.ManualReapplyService
	stateValidation   *stateValidationHandler
}

func newTemplateWorkflowHandler(
	templateGenerator *template_generator.TemplateGenerator,
	mapper *mappers.GeneratorConfigMapper,
	contentProvider *providers.MandatoryContentProvider,
	connectionEditor *connection_editor.ConnectionEditorService,
	zoneEditor *connection_editor.ZoneEditorService,
	manualReapply *connection_editor.ManualReapplyService,
	stateValidation *stateValidationHandler,
) *templateWorkflowHandler {
	return &templateWorkflowHandler{
		templateGenerator: templateGenerator,
		mapper:            mapper,
		contentProvider:   contentProvider,
		connectionEditor:  connectionEditor,
		zoneEditor:        zoneEditor,
		manualReapply:     manualReapply,
		stateValidation:   stateValidation,
	}
}

func (this *templateWorkflowHandler) GenerateTemplate(
	stateDto dtos.EditorStateDto,
) (dtos.TemplateLoadDto, error) {
	validation := this.stateValidation.validateEditorState(stateDto, true)
	stateDto = validation.State

	configuration := this.mapper.FromEditorState(stateDto)
	if configuration.TemplateName == "" {
		return dtos.TemplateLoadDto{}, common_errors.ErrNoTemplateName
	}

	this.templateGenerator.SetConfiguration(configuration)
	template := this.templateGenerator.Generate()
	if template == nil {
		return dtos.TemplateLoadDto{}, common_errors.ErrGeneratedTemplateInvalid
	}

	return dtos.TemplateLoadDto{Template: template, Warnings: validation.Warnings}, nil
}

func (this *templateWorkflowHandler) UpdateTemplate(
	templateDto dtos.TemplateUpdateDto,
) (dtos.TemplateLoadDto, error) {
	if templateDto.Template == nil || len(templateDto.Template.Variants) == 0 {
		return dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid
	}

	newTemplate := *templateDto.Template
	newTemplate.Variants = slices.Clone(templateDto.Template.Variants)
	newTemplate.Variants[0].Zones = templateDto.Zones
	newTemplate.Variants[0].Connections = templateDto.Connections

	this.zoneEditor.RebuildZoneConnectionRoads(
		newTemplate.Variants[0].Zones,
		newTemplate.Variants[0].Connections)

	// Rebuild mandatory content from the final zones so a zone re-tiered in the
	// manual editor gets the content of its new quality instead of the original tier.
	if templateDto.EditorState != nil {
		configuration := this.mapper.FromEditorState(*templateDto.EditorState)
		newTemplate.MandatoryContent = this.contentProvider.CreateContentsForZones(
			*configuration, newTemplate.Variants[0].Zones)
	}

	var err error
	if this.connectionEditor.ComputeHasErrors(
		newTemplate.Variants[0].Zones,
		newTemplate.Variants[0].Connections,
	) {
		err = common_errors.ErrZonesMissing
	}

	return dtos.TemplateLoadDto{Template: &newTemplate}, err
}

func (this *templateWorkflowHandler) ReapplyCastleSettings(
	request dtos.CastleSettingsReapplyRequestDto,
) []entities.Zone {
	configuration := this.mapper.FromEditorState(request.EditorState)
	this.manualReapply.ApplyCastleSettingChanges(request.Zones, request.Changes, configuration)
	return request.Zones
}

func (this *templateWorkflowHandler) ValidateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool,
) dtos.EditorStateValidationDto {
	return this.stateValidation.validateEditorState(stateDto, fixIssues)
}
