package handlers

import (
	"log/slog"
	"slices"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
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
	templateGenerator *template_generator.TemplateGenerator
	mapper            *mappers.GeneratorConfigMapper
	contentProvider   *providers.MandatoryContentProvider
	fileService       *file_service.FileService
	previewGenerator  *preview_service.PreviewGeneratorService
	previewLayout     *preview_service.PreviewLayoutService
	zoneClassifier    *zone_services.ZoneClassifier
	tuningFactory     *generation_tuning.GenerationTuningFactory
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
		previewLayout:     preview_service.NewPreviewLayoutService(),
		zoneClassifier:    zone_services.NewZoneClassifier(),
		tuningFactory:     generation_tuning.NewGenerationTuningFactory(),
	}
}

func (this *GUIHandler) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error) {
	validation := this.ValidateEditorState(stateDto, true)
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

func (this *GUIHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	if templateDto.Template == nil || len(templateDto.Template.Variants) == 0 {
		return dtos.TemplateLoadDto{}, common_errors.ErrProvidedTemplateInvalid
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
	if templateDto.EditorState != nil {
		configuration := this.mapper.FromEditorState(*templateDto.EditorState)
		newTemplate.MandatoryContent = this.contentProvider.CreateContentsForZones(
			*configuration, newTemplate.Variants[0].Zones)
	}

	var err error
	if connection_editor.ComputeHasErrors(newTemplate.Variants[0].Zones, newTemplate.Variants[0].Connections) {
		err = common_errors.ErrZonesMissing
	}

	return dtos.TemplateLoadDto{Template: &newTemplate}, err
}

func (this *GUIHandler) ReapplyCastleSettings(
	request dtos.CastleSettingsReapplyRequestDto,
) []entities.Zone {
	configuration := this.mapper.FromEditorState(request.EditorState)
	connection_editor.ApplyCastleSettingChanges(request.Zones, request.Changes, configuration)
	return request.Zones
}

func (this *GUIHandler) GetZoneEditorOptions(
	state dtos.EditorStateDto,
	totalZoneCount int,
) dtos.ZoneEditorOptionsDto {
	configuration := this.mapper.FromEditorState(state)
	return dtos.ZoneEditorOptionsDto{
		Topology:      state.Topology,
		Tuning:        this.tuningFactory.Create(configuration, totalZoneCount),
		GenerateRoads: state.GenerateRoads,
	}
}

func (this *GUIHandler) CountZoneCastles(zone entities.Zone) int {
	return connection_editor.CountZoneCastles(zone)
}

func (this *GUIHandler) GetZoneQuality(zone entities.Zone) neutral_zone.Quality {
	return this.zoneClassifier.GetQuality(zone)
}

func (this *GUIHandler) GetZoneConnectionGuardQuality(
	from, to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool,
) neutral_zone.Quality {
	playerNames := make([]string, 0, len(playerZoneNames))
	for playerName := range playerZoneNames {
		playerNames = append(playerNames, playerName)
	}
	return this.zoneClassifier.GetConnectionGuardQuality(from, to, zones, playerNames)
}

func (this *GUIHandler) ApplyZoneEditorQuality(
	request dtos.ZoneEditorQualityRequestDto,
) entities.Zone {
	connection_editor.ApplyNeutralZoneQuality(
		&request.Zone,
		request.Quality,
		request.CastleCount,
		request.Tuning,
	)
	return request.Zone
}

func (this *GUIHandler) DescribeZoneEditorGraph(
	zones []entities.Zone,
	connections []entities.Connection,
) dtos.ZoneEditorGraphDto {
	return dtos.ZoneEditorGraphDto{
		HasErrors:         connection_editor.ComputeHasErrors(zones, connections),
		IsolatedZoneCount: len(connection_editor.FindIsolatedZones(zones, connections)),
	}
}

func (this *GUIHandler) CreateZoneEditorConnection(
	request dtos.ZoneEditorConnectionRequestDto,
) entities.Connection {
	return connection_editor.NewDefaultConnection(
		request.From,
		request.To,
		request.Zones,
		request.PlayerZoneNames,
	)
}

func (this *GUIHandler) FindOpenZonePosition(occupied [][2]float64) [2]float64 {
	return connection_editor.FindOpenPosition(occupied)
}

func (this *GUIHandler) GetNextZoneLabel(zones []entities.Zone) string {
	return connection_editor.NextFreeZoneLabel(zones)
}

func (this *GUIHandler) CreateZoneEditorNeutralZone(
	request dtos.ZoneEditorNeutralZoneRequestDto,
) entities.Zone {
	return connection_editor.NewDefaultNeutralZone(
		request.Label,
		request.Quality,
		request.CastleCount,
		request.GenerateRoads,
		request.Tuning,
	)
}

func (this *GUIHandler) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return connection_editor.CanDeleteZone(zoneName, playerZoneNames)
}

func (this *GUIHandler) RemoveZoneEditorZone(
	request dtos.ZoneEditorRemoveRequestDto,
) dtos.ZoneEditorMutationDto {
	zones, connections := connection_editor.RemoveZone(
		request.Zones,
		request.Connections,
		request.ZoneName,
	)
	return dtos.ZoneEditorMutationDto{Zones: zones, Connections: connections}
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

func (this *GUIHandler) BuildPreviewLayout(
	request dtos.PreviewLayoutRequestDto,
) (dtos.PreviewLayoutDto, error) {
	layout := this.previewLayout.BuildPreviewLayout(request.Template, request.Topology, request.CanvasSide)
	return dtos.PreviewLayoutDto{Layout: layout}, nil
}

func (this *GUIHandler) GetContentRuleEditorOptions(
	content models.SidMapping,
) dtos.ContentRuleEditorOptionsDto {
	rules := []dtos.ContentRuleOptionDto{
		{
			Key:         dtos.ContentRuleKeyDistanceToRoad,
			Name:        content_rules.RuleDistanceToRoadName,
			Description: content_rules.RuleDistanceToRoadDescription,
			Marker:      content_rules.RuleDistanceToRoadMarker,
			EditorKind:  dtos.ContentRuleEditorKindDistance,
			EditorLabel: "Distance",
		},
		{
			Key:         dtos.ContentRuleKeyDistanceToTown,
			Name:        content_rules.RuleDistanceToTownName,
			Description: content_rules.RuleDistanceToTownDescription,
			Marker:      content_rules.RuleDistanceToTownMarker,
			EditorKind:  dtos.ContentRuleEditorKindDistance,
			EditorLabel: "Distance",
		},
		{
			Key:         dtos.ContentRuleKeyGuarded,
			Name:        content_rules.RuleGuardedName,
			Description: content_rules.RuleGuardedDescription,
			Marker:      content_rules.RuleGuardedMarker,
			EditorKind:  dtos.ContentRuleEditorKindBoolean,
			EditorLabel: "Guarded",
		},
		{
			Key:         dtos.ContentRuleKeySoloEncounter,
			Name:        content_rules.RuleSoloEncounterName,
			Description: content_rules.RuleSoloEncounterDescription,
			Marker:      content_rules.RuleSoloEncounterMarker,
			EditorKind:  dtos.ContentRuleEditorKindBoolean,
			EditorLabel: "Solo encounter",
		},
	}

	variants := contentRuleVariantOptions(content)
	if len(variants) > 0 {
		rules = append(rules, dtos.ContentRuleOptionDto{
			Key:         dtos.ContentRuleKeyVariant,
			Name:        content_rules.RuleVariantName,
			Description: content_rules.RuleVariantDescription,
			Marker:      content_rules.RuleVariantMarker,
			EditorKind:  dtos.ContentRuleEditorKindVariant,
			EditorLabel: "Variant",
		})
	}

	return dtos.ContentRuleEditorOptionsDto{
		Rules:     rules,
		Distances: content_rules.GetDistanceDisplayNames(),
		Variants:  variants,
	}
}

func (this *GUIHandler) DescribeContentRule(
	content models.SidMapping,
	savedRule models.ContentRuleRowSave,
) dtos.ContentRuleDescriptionDto {
	description := dtos.ContentRuleDescriptionDto{
		Key:         contentRuleKeyFromName(savedRule.Name),
		DisplayText: savedRule.Name,
		SavedRule:   savedRule,
	}
	rule := content_rules.CreateRuleFromSavedRule(savedRule, content)
	if rule == nil {
		return description
	}

	description.DisplayText = rule.DisplayText()
	description.Marker = rule.Marker()
	description.Valid = true
	if savedRule.VariantID != nil {
		variant, ok := content_rules.GetVariantForContentByID(content, *savedRule.VariantID)
		if ok {
			description.VariantLabel, _ = variant.GetVariantByID(*savedRule.VariantID)
		}
	}
	return description
}

func contentRuleVariantOptions(content models.SidMapping) []dtos.ContentRuleVariantOptionDto {
	variants := content_rules.GetVariantsForContent(content)
	options := make([]dtos.ContentRuleVariantOptionDto, 0, len(variants))
	for _, variant := range variants {
		for _, tuple := range variant.Variants {
			options = append(options, dtos.ContentRuleVariantOptionDto{ID: tuple.Key, Label: tuple.Value})
		}
	}
	return options
}

func contentRuleKeyFromName(name string) dtos.ContentRuleKey {
	switch {
	case strings.EqualFold(name, content_rules.RuleDistanceToRoadName):
		return dtos.ContentRuleKeyDistanceToRoad
	case strings.EqualFold(name, content_rules.RuleDistanceToTownName):
		return dtos.ContentRuleKeyDistanceToTown
	case strings.EqualFold(name, content_rules.RuleGuardedName):
		return dtos.ContentRuleKeyGuarded
	case strings.EqualFold(name, content_rules.RuleSoloEncounterName):
		return dtos.ContentRuleKeySoloEncounter
	case strings.EqualFold(name, content_rules.RuleVariantName):
		return dtos.ContentRuleKeyVariant
	default:
		return ""
	}
}

func (this *GUIHandler) ValidateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool,
) dtos.EditorStateValidationDto {
	issues := validators.ValidateEditorState(&stateDto)
	warnings := make([]string, 0, len(issues))
	for _, issue := range issues {
		if fixIssues {
			issue.Fix(&stateDto)
		}
		warnings = append(warnings, issue.Message)
	}

	return dtos.EditorStateValidationDto{State: stateDto, Warnings: warnings}
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

	validation := this.ValidateEditorState(*loaded, fixIssues)
	return &validation.State, validation.Warnings, nil
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
