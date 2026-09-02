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
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
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

	newTemplate := this.templateMapper.ToEntity(*templateDto.Template)
	newTemplate.Variants[0].Zones = templateDto.Zones
	newTemplate.Variants[0].Connections = templateDto.Connections

	this.zoneEditor.RebuildZoneConnectionRoads(
		newTemplate.Variants[0].Zones,
		newTemplate.Variants[0].Connections)

	// Rebuild mandatory content from the final zones so a zone re-tiered in the
	// manual editor gets the content of its new quality instead of the original tier.
	if templateDto.EditorState != nil {
		configuration := this.mapper.FromEditorState(templateDto.EditorState.EditorState)
		newTemplate.MandatoryContent = this.contentProvider.CreateContentsForZones(
			*configuration, newTemplate.Variants[0].Zones)
	}

	var err error
	if this.connectionEditor.ComputeHasErrors(newTemplate.Variants[0].Zones, newTemplate.Variants[0].Connections) {
		err = common_errors.ErrZonesMissing
	}

	updated := this.templateMapper.ToModel(newTemplate)
	carryZoneTiersByName(templateDto.Template.Variants[0].Zones, updated.Variants[0].Zones)

	return dtos.TemplateLoadDto{Template: &updated}, err
}

// carryZoneTiersByName re-attaches the tier a zone already carried after the
// edited zones round-tripped through the entity, which has nowhere to keep it.
//
// This is phase-4 scaffolding: it exists only because TemplateUpdateDto.Zones
// is still []entities.Zone, and it is the last name lookup in the batch. When
// the zone editor moves onto template_model, delete it - do not spread it.
func carryZoneTiersByName(previous []template_model.Zone, updated []template_model.Zone) {
	qualities := make(map[string]*neutral_zone.Quality, len(previous))
	for _, zone := range previous {
		if zone.Quality != nil {
			qualities[zone.Name] = zone.Quality
		}
	}

	for index := range updated {
		if quality, ok := qualities[updated[index].Name]; ok {
			updated[index].Quality = quality
		}
	}
}

func (this *templateHandler) ReapplyCastleSettings(
	request dtos.CastleSettingsReapplyRequestDto) []entities.Zone {
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
	// required, so this is where the model goes back to being an entity.
	template := this.templateMapper.ToEntity(*templateDto.Template)
	previewImage := this.previewGenerator.CreatePreviewImage(&template, templateDto.Topology)
	return this.fileService.SaveTemplateWithPreview(outputPath, &template, previewImage)
}
