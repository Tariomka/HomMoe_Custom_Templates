package handlers

import (
	"log/slog"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
)

type GUIHandler struct {
	templateHandler    handler_interfaces.ITemplateHandler
	previewHandler     handler_interfaces.IPreviewHandler
	stateHandler       handler_interfaces.IStateHandler
	contentRuleHandler handler_interfaces.IContentRuleHandler
	zoneEditorHandler  handler_interfaces.IZoneEditorHandler
}

func NewDefaultGuiHandler() handler_interfaces.IGuiHandler {
	return NewGuiHandler(nil, nil, nil, nil, nil)
}

func NewGuiHandler(
	templateHandler handler_interfaces.ITemplateHandler,
	stateHandler handler_interfaces.IStateHandler,
	previewHandler handler_interfaces.IPreviewHandler,
	contentRuleHandler handler_interfaces.IContentRuleHandler,
	zoneEditorHandler handler_interfaces.IZoneEditorHandler) handler_interfaces.IGuiHandler {
	previewGenerator, err := preview_service.NewPreviewGenerator()
	if err != nil {
		slog.Error(
			"Preview Generator failed to initialize, preview images will not be generated",
			slog.String("error", err.Error()))
	}

	zoneClassifier := zone_services.NewZoneClassifier()
	castleFactory := zone_services.NewCastleFactory()
	roadFactory := zone_services.NewRoadFactory()
	zoneFactory := zone_services.NewZoneFactory(castleFactory, roadFactory)
	zoneEditor := connection_editor.NewZoneEditorService(castleFactory, roadFactory, zoneFactory)
	tuningFactory := generation_tuning.NewGenerationTuningFactory()
	fileService := file_service.NewFileService()
	mapper := mappers.NewConfigMapper()
	connectionEditor := connection_editor.NewConnectionEditorService(zoneClassifier)
	manualReapply := connection_editor.NewManualReapplyService(
		zoneEditor,
		zoneClassifier,
		tuningFactory)

	if previewHandler == nil {
		previewHandler = newPreviewHandler(preview_service.NewPreviewLayoutService())
	}
	if contentRuleHandler == nil {
		contentRuleHandler = newContentRuleHandler(content_rules.NewContentRuleService())
	}
	if stateHandler == nil {
		stateHandler = newStateHandler(fileService, validators.NewEditorStateValidator())
	}
	if zoneEditorHandler == nil {
		zoneEditorHandler = newZoneEditorHandler(
			mapper,
			zoneClassifier,
			connectionEditor,
			zoneEditor,
			tuningFactory)
	}
	if templateHandler == nil {
		templateHandler = newTemplateHandler(
			template_generator.NewTemplateGenerator(nil, castleFactory, roadFactory, zoneFactory),
			mapper,
			providers.NewMandatoryContentProvider(zoneClassifier, zoneEditor),
			connectionEditor,
			zoneEditor,
			manualReapply,
			fileService,
			previewGenerator,
			stateHandler)
	}

	return &GUIHandler{
		templateHandler:    templateHandler,
		stateHandler:       stateHandler,
		previewHandler:     previewHandler,
		contentRuleHandler: contentRuleHandler,
		zoneEditorHandler:  zoneEditorHandler,
	}
}

func (this *GUIHandler) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error) {
	return this.templateHandler.GenerateTemplate(stateDto)
}

func (this *GUIHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	return this.templateHandler.UpdateTemplate(templateDto)
}

func (this *GUIHandler) ReapplyCastleSettings(request dtos.CastleSettingsReapplyRequestDto) []entities.Zone {
	return this.templateHandler.ReapplyCastleSettings(request)
}

func (this *GUIHandler) GetZoneEditorOptions(state dtos.EditorStateDto, totalZoneCount int) dtos.ZoneEditorOptionsDto {
	return this.zoneEditorHandler.GetZoneEditorOptions(state, totalZoneCount)
}

func (this *GUIHandler) CountZoneCastles(zone entities.Zone) int {
	return this.zoneEditorHandler.CountZoneCastles(zone)
}

func (this *GUIHandler) GetZoneQuality(zone entities.Zone) neutral_zone.Quality {
	return this.zoneEditorHandler.GetZoneQuality(zone)
}

func (this *GUIHandler) GetZoneConnectionGuardQuality(
	from, to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool) neutral_zone.Quality {
	return this.zoneEditorHandler.GetZoneConnectionGuardQuality(from, to, zones, playerZoneNames)
}

func (this *GUIHandler) ApplyZoneEditorQuality(request dtos.ZoneEditorQualityRequestDto) entities.Zone {
	return this.zoneEditorHandler.ApplyZoneEditorQuality(request)
}

func (this *GUIHandler) DescribeZoneEditorGraph(
	zones []entities.Zone,
	connections []entities.Connection) dtos.ZoneEditorGraphDto {
	return this.zoneEditorHandler.DescribeZoneEditorGraph(zones, connections)
}

func (this *GUIHandler) CreateZoneEditorConnection(
	request dtos.ZoneEditorConnectionRequestDto) entities.Connection {
	return this.zoneEditorHandler.CreateZoneEditorConnection(request)
}

func (this *GUIHandler) FindOpenZonePosition(occupied [][2]float64) [2]float64 {
	return this.zoneEditorHandler.FindOpenZonePosition(occupied)
}

func (this *GUIHandler) GetNextZoneLabel(zones []entities.Zone) string {
	return this.zoneEditorHandler.GetNextZoneLabel(zones)
}

func (this *GUIHandler) CreateZoneEditorNeutralZone(request dtos.ZoneEditorNeutralZoneRequestDto) entities.Zone {
	return this.zoneEditorHandler.CreateZoneEditorNeutralZone(request)
}

func (this *GUIHandler) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return this.zoneEditorHandler.CanDeleteZone(zoneName, playerZoneNames)
}

func (this *GUIHandler) RemoveZoneEditorZone(request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto {
	return this.zoneEditorHandler.RemoveZoneEditorZone(request)
}

func (this *GUIHandler) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	return this.templateHandler.SaveTemplate(templateDto)
}

func (this *GUIHandler) BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error) {
	return this.previewHandler.BuildPreviewLayout(request)
}

func (this *GUIHandler) GetContentRuleEditorOptions(content models.SidMapping) dtos.ContentRuleEditorOptionsDto {
	return this.contentRuleHandler.GetContentRuleEditorOptions(content)
}

func (this *GUIHandler) DescribeContentRule(
	content models.SidMapping,
	savedRule models.ContentRuleRowSave) dtos.ContentRuleDescriptionDto {
	return this.contentRuleHandler.DescribeContentRule(content, savedRule)
}

func (this *GUIHandler) ValidateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool) dtos.EditorStateValidationDto {
	return this.stateHandler.ValidateEditorState(stateDto, fixIssues)
}

func (this *GUIHandler) LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error) {
	return this.stateHandler.LoadState(path, fixIssues)
}

func (this *GUIHandler) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	return this.stateHandler.SaveState(stateDto)
}
