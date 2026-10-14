package handlers

import (
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor/connection_editor_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type templateHandler struct {
	templateGenerator template_generator.ITemplateGenerator
	mapper            mappers.IGeneratorConfigMapper
	contentProvider   provider_interfaces.IMandatoryContentProvider
	connectionEditor  connection_editor_interfaces.IConnectionEditorService
	zoneEditor        connection_editor_interfaces.IZoneEditorService
	roadPolicy        zone_interfaces.IRoadPolicyService
	manualEditDelta   connection_editor_interfaces.IManualEditDeltaService
	fileService       file_service.IFileService
	previewGenerator  preview_service.IPreviewGeneratorService
	stateHandler      handler_interfaces.IStateHandler
}

func NewTemplateHandler(
	templateGenerator template_generator.ITemplateGenerator,
	mapper mappers.IGeneratorConfigMapper,
	contentProvider provider_interfaces.IMandatoryContentProvider,
	connectionEditor connection_editor_interfaces.IConnectionEditorService,
	zoneEditor connection_editor_interfaces.IZoneEditorService,
	roadPolicy zone_interfaces.IRoadPolicyService,
	manualEditDelta connection_editor_interfaces.IManualEditDeltaService,
	fileService file_service.IFileService,
	previewGenerator preview_service.IPreviewGeneratorService,
	stateHandler handler_interfaces.IStateHandler) handler_interfaces.ITemplateHandler {
	return &templateHandler{
		templateGenerator: templateGenerator,
		mapper:            mapper,
		contentProvider:   contentProvider,
		connectionEditor:  connectionEditor,
		zoneEditor:        zoneEditor,
		roadPolicy:        roadPolicy,
		manualEditDelta:   manualEditDelta,
		fileService:       fileService,
		previewGenerator:  previewGenerator,
		stateHandler:      stateHandler,
	}
}

func (this *templateHandler) GenerateTemplate(
	state editor_state_dto.EditorStateDto) (dtos.TemplateLoadDto, error) {
	generated, _, warnings, err := this.generate(state)
	if err != nil {
		return dtos.TemplateLoadDto{}, err
	}

	return dtos.TemplateLoadDto{Template: generated, Warnings: warnings}, nil
}

func (this *templateHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	if templateDto.Template == nil || len(templateDto.Template.Variants) == 0 {
		return dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid
	}

	zones := templateDto.Zones
	connections := templateDto.Connections
	this.zoneEditor.EnsureConnectionNames(connections)

	updated := templateDto.Template.Clone()
	updated.Variants[0].Zones = zones
	updated.Variants[0].Connections = connections

	if templateDto.EditorState != nil {
		configuration := this.mapper.FromEditorState(templateDto.EditorState.EditorState)
		updated.MandatoryContent = this.contentProvider.CreateContentsForZones(*configuration, zones)
		this.roadPolicy.Reconcile(models.RoadReconciliationRequest{
			Zones:                zones,
			Connections:          connections,
			MandatoryContent:     updated.MandatoryContent,
			GenerateRoads:        configuration.GenerateRoads,
			SpawnRemoteFootholds: configuration.SpawnRemoteFootholds,
			RemoteFootholdCount:  configuration.RemoteFootholdCount,
		})
	}

	var err error
	if this.connectionEditor.ComputeHasErrors(zones, connections) {
		err = common_errors.ErrZonesMissing
	}

	return dtos.TemplateLoadDto{Template: &updated}, err
}

func (this *templateHandler) CaptureManualEdits(
	request dtos.ManualEditCaptureRequestDto) (dtos.ManualEditCaptureDto, error) {
	return this.captureManualEdits(request.State, request.Zones, request.Connections)
}

func (this *templateHandler) ExpandManualEdits(request dtos.ManualEditExpandRequestDto) dtos.ManualEditExpansionDto {
	validation := this.stateHandler.ValidateEditorState(request.State.EditorState, true)
	configuration := this.mapper.FromEditorState(validation.State)
	freshZones, freshConnections := variantLayout(request.Template)
	return dtos.ManualEditExpansionDto(this.manualEditDelta.Expand(models.ManualEditExpandRequest{
		FreshZones:       freshZones,
		FreshConnections: freshConnections,
		ManualEdits:      request.State.ManualEditSettings,
		CastleChanges:    request.CastleChanges,
		Configuration:    configuration,
		NonDeterministic: configuration.HasNonDeterministicLayout(),
	}))
}

func (this *templateHandler) ConvertManualEdits(
	request dtos.ManualEditConvertRequestDto) (dtos.ManualEditCaptureDto, error) {
	zones, connections := request.ManualEdits.ToSnapshot()
	return this.captureManualEdits(request.Baseline, zones, connections)
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

// generate validates the state and generates a template from it; manual
// edits are not generator input.
func (this *templateHandler) generate(
	state editor_state_dto.EditorStateDto) (*template_model.Template, *config.GeneratorConfig, []string, error) {
	validation := this.stateHandler.ValidateEditorState(state.EditorState, true)
	if validation.Rejection != nil {
		return nil, nil, nil, validation.Rejection
	}

	configuration := this.mapper.FromEditorState(validation.State)
	if configuration.TemplateName == "" {
		return nil, nil, nil, common_errors.ErrNoTemplateName
	}

	this.templateGenerator.SetConfiguration(configuration)
	generated, generationWarnings, err := this.templateGenerator.Generate()
	if err != nil {
		return nil, nil, nil, err
	}
	if generated == nil {
		return nil, nil, nil, common_errors.ErrGeneratedTemplateInvalid
	}

	return generated, configuration, slices.Concat(validation.Warnings, generationWarnings), nil
}

func (this *templateHandler) captureManualEdits(
	settings editor_state_dto.EditorStateDto,
	zones []template_model.Zone,
	connections []template_model.Connection) (dtos.ManualEditCaptureDto, error) {
	fresh, configuration, _, err := this.generate(settings)
	if err != nil {
		return dtos.ManualEditCaptureDto{}, err
	}

	freshZones, freshConnections := variantLayout(fresh)
	return dtos.ManualEditCaptureDto{ManualEdits: this.manualEditDelta.Capture(models.ManualEditCaptureRequest{
		FreshZones:       freshZones,
		FreshConnections: freshConnections,
		Zones:            zones,
		Connections:      connections,
		Configuration:    configuration,
		NonDeterministic: configuration.HasNonDeterministicLayout(),
	})}, nil
}

// variantLayout returns the zones and connections of the template's first
// variant, or none when it has no variant.
func variantLayout(template *template_model.Template) ([]template_model.Zone, []template_model.Connection) {
	if template == nil || len(template.Variants) == 0 {
		return nil, nil
	}

	return template.Variants[0].Zones, template.Variants[0].Connections
}
