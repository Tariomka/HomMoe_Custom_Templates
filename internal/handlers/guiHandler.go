package handlers

import (
	"errors"
	"log/slog"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
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
	templateWorkflow    ITemplateWorkflow
	templatePersistence ITemplatePersistence
	preview             IPreview
	statePersistence    IStatePersistence
	contentRule         IContentRule
	zoneEditor          IZoneEditor
}

func NewGuiHandler() *GUIHandler {
	previewGenerator, err := preview_service.NewPreviewGenerator()
	if err != nil {
		slog.Error(
			"Preview Generator failed to initialize, preview images will not be generated",
			slog.String("error", err.Error()))
	}
	zoneClassifier := zone_services.NewZoneClassifier()
	creationServices := zone_services.NewCreationServices(nil, nil)
	zoneEditor := connection_editor.NewZoneEditorServiceWithCreationServices(creationServices)
	tuningFactory := generation_tuning.NewGenerationTuningFactory()
	fileService := file_service.NewFileService()
	stateValidation := newStateValidationHandler(validators.NewEditorStateValidator())
	mapper := mappers.NewConfigMapper()
	connectionEditor := connection_editor.NewConnectionEditorService(zoneClassifier)
	contentProvider := providers.NewMandatoryContentProviderWithDependencies(zoneClassifier, zoneEditor)
	manualReapply := connection_editor.NewManualReapplyServiceWithDependencies(
		zoneEditor,
		zoneClassifier,
		tuningFactory,
	)

	dependencies := GUIHandlerDependencies{
		TemplateWorkflow: newTemplateWorkflowHandler(
			template_generator.NewTemplateGeneratorWithCreationServices(nil, creationServices),
			mapper,
			contentProvider,
			connectionEditor,
			zoneEditor,
			manualReapply,
			stateValidation,
		),
		TemplatePersistence: newTemplatePersistenceHandler(fileService, previewGenerator),
		Preview:             newPreviewHandler(preview_service.NewPreviewLayoutService()),
		StatePersistence:    newStatePersistenceHandler(fileService, stateValidation),
		ContentRule:         newContentRuleHandler(content_rules.NewContentRuleService()),
		ZoneEditor: newZoneEditorHandler(
			mapper,
			zoneClassifier,
			connectionEditor,
			zoneEditor,
			tuningFactory,
		),
	}
	handler, err := NewGuiHandlerWithDependencies(dependencies)
	if err != nil {
		panic(err)
	}
	return handler
}

func NewGuiHandlerWithDependencies(
	dependencies GUIHandlerDependencies,
) (*GUIHandler, error) {
	if dependencies.TemplateWorkflow == nil {
		return nil, errors.New("template workflow handler is required")
	}
	if dependencies.StatePersistence == nil {
		return nil, errors.New("state persistence handler is required")
	}
	if dependencies.TemplatePersistence == nil {
		return nil, errors.New("template persistence handler is required")
	}
	if dependencies.Preview == nil {
		return nil, errors.New("preview handler is required")
	}
	if dependencies.ContentRule == nil {
		return nil, errors.New("content rule handler is required")
	}
	if dependencies.ZoneEditor == nil {
		return nil, errors.New("zone editor handler is required")
	}

	return &GUIHandler{
		templateWorkflow:    dependencies.TemplateWorkflow,
		statePersistence:    dependencies.StatePersistence,
		templatePersistence: dependencies.TemplatePersistence,
		preview:             dependencies.Preview,
		contentRule:         dependencies.ContentRule,
		zoneEditor:          dependencies.ZoneEditor,
	}, nil
}

func (this *GUIHandler) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error) {
	return this.templateWorkflow.GenerateTemplate(stateDto)
}

func (this *GUIHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	return this.templateWorkflow.UpdateTemplate(templateDto)
}

func (this *GUIHandler) ReapplyCastleSettings(request dtos.CastleSettingsReapplyRequestDto) []entities.Zone {
	return this.templateWorkflow.ReapplyCastleSettings(request)
}

func (this *GUIHandler) GetZoneEditorOptions(state dtos.EditorStateDto, totalZoneCount int) dtos.ZoneEditorOptionsDto {
	return this.zoneEditor.GetZoneEditorOptions(state, totalZoneCount)
}

func (this *GUIHandler) CountZoneCastles(zone entities.Zone) int {
	return this.zoneEditor.CountZoneCastles(zone)
}

func (this *GUIHandler) GetZoneQuality(zone entities.Zone) neutral_zone.Quality {
	return this.zoneEditor.GetZoneQuality(zone)
}

func (this *GUIHandler) GetZoneConnectionGuardQuality(
	from, to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool) neutral_zone.Quality {
	return this.zoneEditor.GetZoneConnectionGuardQuality(from, to, zones, playerZoneNames)
}

func (this *GUIHandler) ApplyZoneEditorQuality(request dtos.ZoneEditorQualityRequestDto) entities.Zone {
	return this.zoneEditor.ApplyZoneEditorQuality(request)
}

func (this *GUIHandler) DescribeZoneEditorGraph(
	zones []entities.Zone,
	connections []entities.Connection) dtos.ZoneEditorGraphDto {
	return this.zoneEditor.DescribeZoneEditorGraph(zones, connections)
}

func (this *GUIHandler) CreateZoneEditorConnection(
	request dtos.ZoneEditorConnectionRequestDto) entities.Connection {
	return this.zoneEditor.CreateZoneEditorConnection(request)
}

func (this *GUIHandler) FindOpenZonePosition(occupied [][2]float64) [2]float64 {
	return this.zoneEditor.FindOpenZonePosition(occupied)
}

func (this *GUIHandler) GetNextZoneLabel(zones []entities.Zone) string {
	return this.zoneEditor.GetNextZoneLabel(zones)
}

func (this *GUIHandler) CreateZoneEditorNeutralZone(request dtos.ZoneEditorNeutralZoneRequestDto) entities.Zone {
	return this.zoneEditor.CreateZoneEditorNeutralZone(request)
}

func (this *GUIHandler) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return this.zoneEditor.CanDeleteZone(zoneName, playerZoneNames)
}

func (this *GUIHandler) RemoveZoneEditorZone(request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto {
	return this.zoneEditor.RemoveZoneEditorZone(request)
}

func (this *GUIHandler) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	return this.templatePersistence.SaveTemplate(templateDto)
}

func (this *GUIHandler) BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error) {
	return this.preview.BuildPreviewLayout(request)
}

func (this *GUIHandler) GetContentRuleEditorOptions(content models.SidMapping) dtos.ContentRuleEditorOptionsDto {
	return this.contentRule.GetContentRuleEditorOptions(content)
}

func (this *GUIHandler) DescribeContentRule(
	content models.SidMapping,
	savedRule models.ContentRuleRowSave) dtos.ContentRuleDescriptionDto {
	return this.contentRule.DescribeContentRule(content, savedRule)
}

func (this *GUIHandler) ValidateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool) dtos.EditorStateValidationDto {
	return this.templateWorkflow.ValidateEditorState(stateDto, fixIssues)
}

func (this *GUIHandler) LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error) {
	return this.statePersistence.LoadState(path, fixIssues)
}

func (this *GUIHandler) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	return this.statePersistence.SaveState(stateDto)
}
