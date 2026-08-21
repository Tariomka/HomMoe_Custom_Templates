package handlers

import (
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
)

type templateHandler struct {
	templateGenerator template_generator.ITemplateGenerator
	mapper            mappers.IGeneratorConfigMapper
	contentProvider   provider_interfaces.IMandatoryContentProvider
	connectionEditor  connection_editor.IConnectionEditorService
	zoneEditor        connection_editor.IZoneEditorService
	manualReapply     connection_editor.IManualReapplyService
	fileService       file_service.IFileService
	previewGenerator  preview_service.IPreviewGeneratorService
	stateHandler      handler_interfaces.IStateHandler
}

func NewTemplateHandler(
	templateGenerator template_generator.ITemplateGenerator,
	mapper mappers.IGeneratorConfigMapper,
	contentProvider provider_interfaces.IMandatoryContentProvider,
	connectionEditor connection_editor.IConnectionEditorService,
	zoneEditor connection_editor.IZoneEditorService,
	manualReapply connection_editor.IManualReapplyService,
	fileService file_service.IFileService,
	previewGenerator preview_service.IPreviewGeneratorService,
	stateHandler handler_interfaces.IStateHandler) handler_interfaces.ITemplateHandler {
	return &templateHandler{
		templateGenerator: templateGenerator,
		mapper:            mapper,
		contentProvider:   contentProvider,
		connectionEditor:  connectionEditor,
		zoneEditor:        zoneEditor,
		manualReapply:     manualReapply,
		fileService:       fileService,
		previewGenerator:  previewGenerator,
		stateHandler:      stateHandler,
	}
}

func (this *templateHandler) GenerateTemplate(stateDto editor_state_dto.EditorStateDto) (dtos.TemplateLoadDto, error) {
	validation := this.stateHandler.ValidateEditorState(stateDto, true)
	stateDto = validation.State

	configuration := this.mapper.FromEditorState(stateDto)
	if configuration.TemplateName == "" {
		return dtos.TemplateLoadDto{}, common_errors.ErrNoTemplateName
	}

	this.templateGenerator.SetConfiguration(configuration)
	template, generationWarnings := this.templateGenerator.Generate()
	if template == nil {
		return dtos.TemplateLoadDto{}, common_errors.ErrGeneratedTemplateInvalid
	}

	warnings := slices.Concat(validation.Warnings, generationWarnings)
	return dtos.TemplateLoadDto{Template: template, Warnings: warnings}, nil
}

func (this *templateHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
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
	if this.connectionEditor.ComputeHasErrors(newTemplate.Variants[0].Zones, newTemplate.Variants[0].Connections) {
		err = common_errors.ErrZonesMissing
	}

	return dtos.TemplateLoadDto{Template: &newTemplate}, err
}

func (this *templateHandler) ReapplyCastleSettings(
	request dtos.CastleSettingsReapplyRequestDto) []entities.Zone {
	configuration := this.mapper.FromEditorState(request.EditorState)
	this.manualReapply.ApplyCastleSettingChanges(request.Zones, request.Changes, configuration)
	return request.Zones
}

func (this *templateHandler) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	if templateDto.Template == nil {
		return "", common_errors.ErrNothingToSave
	}

	outputPath := strings.TrimSpace(templateDto.OutputPath)
	if outputPath == "" {
		return "", common_errors.ErrNoOutputPath
	}

	previewImage := this.previewGenerator.CreatePreviewImage(templateDto.Template, templateDto.Topology)
	return this.fileService.SaveTemplateWithPreview(outputPath, templateDto.Template, previewImage)
}
