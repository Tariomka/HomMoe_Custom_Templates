package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type GUIHandler struct {
	templateHandler    handler_interfaces.ITemplateHandler
	previewHandler     handler_interfaces.IPreviewHandler
	stateHandler       handler_interfaces.IStateHandler
	contentRuleHandler handler_interfaces.IZoneContentHandler
	zoneEditorHandler  handler_interfaces.IZoneEditorHandler
	bonusHandler       handler_interfaces.IBonusHandler
}

func NewGuiHandler(
	templateHandler handler_interfaces.ITemplateHandler,
	stateHandler handler_interfaces.IStateHandler,
	previewHandler handler_interfaces.IPreviewHandler,
	contentRuleHandler handler_interfaces.IZoneContentHandler,
	zoneEditorHandler handler_interfaces.IZoneEditorHandler,
	bonusHandler handler_interfaces.IBonusHandler) handler_interfaces.IGuiHandler {
	return &GUIHandler{
		templateHandler:    templateHandler,
		stateHandler:       stateHandler,
		previewHandler:     previewHandler,
		contentRuleHandler: contentRuleHandler,
		zoneEditorHandler:  zoneEditorHandler,
		bonusHandler:       bonusHandler,
	}
}

func (this *GUIHandler) GenerateTemplate(
	state editor_state_dto.EditorStateDto) (dtos.TemplateLoadDto, error) {
	return this.templateHandler.GenerateTemplate(state)
}

func (this *GUIHandler) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	return this.templateHandler.UpdateTemplate(templateDto)
}

func (this *GUIHandler) ReapplyCastleSettings(
	request dtos.CastleSettingsReapplyRequestDto) []template_model.Zone {
	return this.templateHandler.ReapplyCastleSettings(request)
}

func (this *GUIHandler) GetZoneEditorOptions(
	state editor_state_dto.EditorStateDto,
	totalZoneCount int) dtos.ZoneEditorOptionsDto {
	return this.zoneEditorHandler.GetZoneEditorOptions(state, totalZoneCount)
}

func (this *GUIHandler) CountZoneCastles(zone template_model.Zone) int {
	return this.zoneEditorHandler.CountZoneCastles(zone)
}

func (this *GUIHandler) GetZoneQuality(zone template_model.Zone) neutral_zone.Quality {
	return this.zoneEditorHandler.GetZoneQuality(zone)
}

func (this *GUIHandler) GetZoneConnectionGuardQuality(
	from, to string,
	zones []template_model.Zone,
	playerZoneNames map[string]bool) neutral_zone.Quality {
	return this.zoneEditorHandler.GetZoneConnectionGuardQuality(from, to, zones, playerZoneNames)
}

func (this *GUIHandler) ApplyZoneEditorQuality(request dtos.ZoneEditorQualityRequestDto) template_model.Zone {
	return this.zoneEditorHandler.ApplyZoneEditorQuality(request)
}

func (this *GUIHandler) DescribeZoneEditorGraph(
	zones []template_model.Zone,
	connections []entities.Connection) dtos.ZoneEditorGraphDto {
	return this.zoneEditorHandler.DescribeZoneEditorGraph(zones, connections)
}

func (this *GUIHandler) CreateZoneEditorConnection(request dtos.ZoneEditorConnectionRequestDto) entities.Connection {
	return this.zoneEditorHandler.CreateZoneEditorConnection(request)
}

func (this *GUIHandler) FindOpenZonePosition(occupied [][2]float64) [2]float64 {
	return this.zoneEditorHandler.FindOpenZonePosition(occupied)
}

func (this *GUIHandler) GetNextZoneLabel(zones []template_model.Zone) string {
	return this.zoneEditorHandler.GetNextZoneLabel(zones)
}

func (this *GUIHandler) CreateZoneEditorNeutralZone(request dtos.ZoneEditorNeutralZoneRequestDto) template_model.Zone {
	return this.zoneEditorHandler.CreateZoneEditorNeutralZone(request)
}

func (this *GUIHandler) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return this.zoneEditorHandler.CanDeleteZone(zoneName, playerZoneNames)
}

func (this *GUIHandler) RemoveZoneEditorZone(request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto {
	return this.zoneEditorHandler.RemoveZoneEditorZone(request)
}

func (this *GUIHandler) BuildZoneEditorGeometry(request dtos.ZoneEditorGeometryRequestDto) models.ZoneEditorGeometry {
	return this.zoneEditorHandler.BuildZoneEditorGeometry(request)
}

func (this *GUIHandler) HitTestZoneEditorNode(request dtos.ZoneEditorHitTestRequestDto) string {
	return this.zoneEditorHandler.HitTestZoneEditorNode(request)
}

func (this *GUIHandler) HitTestZoneEditorEdge(position models.Position, edges []models.ZoneEditorEdge) int {
	return this.zoneEditorHandler.HitTestZoneEditorEdge(position, edges)
}

func (this *GUIHandler) GetZoneEditorGridStep(zoneRadius float64) float64 {
	return this.zoneEditorHandler.GetZoneEditorGridStep(zoneRadius)
}

func (this *GUIHandler) SnapZoneEditorPosition(request dtos.ZoneEditorSnapRequestDto) models.ZoneEditorSnapResult {
	return this.zoneEditorHandler.SnapZoneEditorPosition(request)
}

func (this *GUIHandler) DescribeExistingBonuses(existing []config.BonusEntry) dtos.ExistingBonusesDto {
	return this.bonusHandler.DescribeExistingBonuses(existing)
}

func (this *GUIHandler) BuildBonusEntries(request dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto {
	return this.bonusHandler.BuildBonusEntries(request)
}

func (this *GUIHandler) FilterNewBonusEntries(
	entries []config.BonusEntry,
	existingKeys map[string]bool) []config.BonusEntry {
	return this.bonusHandler.FilterNewBonusEntries(entries, existingKeys)
}

func (this *GUIHandler) GetSpellCountLabel(count int) string {
	return this.bonusHandler.GetSpellCountLabel(count)
}

func (this *GUIHandler) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	return this.templateHandler.SaveTemplate(templateDto)
}

func (this *GUIHandler) BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) dtos.PreviewLayoutDto {
	return this.previewHandler.BuildPreviewLayout(request)
}

func (this *GUIHandler) GetContentRuleEditorOptions(content models.SidMapping) dtos.ContentRuleEditorOptionsDto {
	return this.contentRuleHandler.GetContentRuleEditorOptions(content)
}

func (this *GUIHandler) DescribeContentRule(
	content models.SidMapping,
	savedRule editor_state_model.ContentRuleRow) dtos.ContentRuleDescriptionDto {
	return this.contentRuleHandler.DescribeContentRule(content, savedRule)
}

func (this *GUIHandler) ComposeContentRule(
	request dtos.ContentRuleCompositionRequestDto) dtos.ContentRuleCompositionResultDto {
	return this.contentRuleHandler.ComposeContentRule(request)
}

func (this *GUIHandler) UpsertContentRule(
	rules []editor_state_model.ContentRuleRow,
	rule editor_state_model.ContentRuleRow) []editor_state_model.ContentRuleRow {
	return this.contentRuleHandler.UpsertContentRule(rules, rule)
}

func (this *GUIHandler) GetDefaultContentRules(content models.SidMapping) []editor_state_model.ContentRuleRow {
	return this.contentRuleHandler.GetDefaultContentRules(content)
}

func (this *GUIHandler) GetContentRuleMarkers(
	content models.SidMapping,
	rules []editor_state_model.ContentRuleRow) string {
	return this.contentRuleHandler.GetContentRuleMarkers(content, rules)
}

func (this *GUIHandler) GetContentRowDisplayName(
	content models.SidMapping,
	rules []editor_state_model.ContentRuleRow) string {
	return this.contentRuleHandler.GetContentRowDisplayName(content, rules)
}

func (this *GUIHandler) SortContentItemsByName(items []models.SidMapping) []models.SidMapping {
	return this.contentRuleHandler.SortContentItemsByName(items)
}

func (this *GUIHandler) ClampContentCount(count int, maxCount int) int {
	return this.contentRuleHandler.ClampContentCount(count, maxCount)
}

func (this *GUIHandler) ValidateEditorState(
	state editor_state_model.EditorState,
	fixIssues bool) editor_state_dto.EditorStateValidationDto {
	return this.stateHandler.ValidateEditorState(state, fixIssues)
}

func (this *GUIHandler) LoadState(path string, fixIssues bool) (*editor_state_dto.EditorStateValidationDto, error) {
	return this.stateHandler.LoadState(path, fixIssues)
}

func (this *GUIHandler) SaveState(stateDto editor_state_dto.EditorStateSaveDto) (string, error) {
	return this.stateHandler.SaveState(stateDto)
}
