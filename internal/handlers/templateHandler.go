package handlers

import (
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
)

type templateHandler struct {
	templateGenerator template_generator.ITemplateGenerator
	mapper            mappers.IGeneratorConfigMapper
	templateMapper    mappers.ITemplateMapper
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
	templateMapper mappers.ITemplateMapper,
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
		templateMapper:    templateMapper,
		contentProvider:   contentProvider,
		connectionEditor:  connectionEditor,
		zoneEditor:        zoneEditor,
		manualReapply:     manualReapply,
		fileService:       fileService,
		previewGenerator:  previewGenerator,
		stateHandler:      stateHandler,
	}
}

func (this *templateHandler) GenerateTemplate(
	state editor_state_dto.EditorStateDto) (dtos.TemplateLoadDto, error) {
	validation := this.stateHandler.ValidateEditorState(state.EditorState, true)

	configuration := this.mapper.FromEditorState(validation.State)
	if configuration.TemplateName == "" {
		return dtos.TemplateLoadDto{}, common_errors.ErrNoTemplateName
	}

	this.templateGenerator.SetConfiguration(configuration)
	generated, generationWarnings := this.templateGenerator.Generate()
	if generated == nil {
		return dtos.TemplateLoadDto{}, common_errors.ErrGeneratedTemplateInvalid
	}

	warnings := slices.Concat(validation.Warnings, generationWarnings)
	return dtos.TemplateLoadDto{Template: generated, Warnings: warnings}, nil
}

func (this *templateHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	if templateDto.Template == nil || len(templateDto.Template.Variants) == 0 {
		return dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid
	}

	zones := templateDto.Zones
	connections := templateDto.Connections
	this.zoneEditor.RebuildZoneConnectionRoads(zones, connections)

	newTemplate := this.templateMapper.ToEntity(*templateDto.Template)
	newTemplate.Variants[0].Zones = template_model.ToZoneEntities(zones)
	newTemplate.Variants[0].Connections = connections

	// Rebuild mandatory content from the final zones so a zone re-tiered in the
	// manual editor gets the content of its new quality instead of the original tier.
	if templateDto.EditorState != nil {
		configuration := this.mapper.FromEditorState(templateDto.EditorState.EditorState)
		newTemplate.MandatoryContent = this.contentProvider.CreateContentsForZones(*configuration, zones)
	}

	var err error
	if this.connectionEditor.ComputeHasErrors(zones, connections) {
		err = common_errors.ErrZonesMissing
	}

	// Re-attaching the applied zones rather than the ones the round trip produced
	// is what keeps each zone's recorded tier: the .rmg.json entity has nowhere
	// to put it, so a zone that went through ToEntity comes back with none.
	updated := this.templateMapper.ToModel(newTemplate)
	updated.Variants[0].Zones = zones

	return dtos.TemplateLoadDto{Template: &updated}, err
}

func (this *templateHandler) ReapplyCastleSettings(
	request dtos.CastleSettingsReapplyRequestDto) []template_model.Zone {
	configuration := this.mapper.FromEditorState(request.EditorState.EditorState)
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

	// Writing the .rmg.json is one of the two places the wire format is genuinely
	// required, so this is where the model goes back to being an entity. The
	// preview renders from the model, which is what lets it colour a zone by the
	// tier the user picked rather than by the one its content pools imply.
	previewImage := this.previewGenerator.CreatePreviewImage(templateDto.Template, templateDto.Topology)
	template := this.templateMapper.ToEntity(*templateDto.Template)
	return this.fileService.SaveTemplateWithPreview(outputPath, &template, previewImage)
}
