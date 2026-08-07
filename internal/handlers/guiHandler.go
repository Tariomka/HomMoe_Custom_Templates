package handlers

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type GUIHandler struct {
	templateHandler    handler_interfaces.ITemplateHandler
	previewHandler     handler_interfaces.IPreviewHandler
	stateHandler       handler_interfaces.IStateHandler
	contentRuleHandler handler_interfaces.IContentRuleHandler
	zoneEditorHandler  handler_interfaces.IZoneEditorHandler
}

func NewGuiHandler(
	templateHandler handler_interfaces.ITemplateHandler,
	stateHandler handler_interfaces.IStateHandler,
	previewHandler handler_interfaces.IPreviewHandler,
	contentRuleHandler handler_interfaces.IContentRuleHandler,
	zoneEditorHandler handler_interfaces.IZoneEditorHandler) handler_interfaces.IGuiHandler {
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

func (this *GUIHandler) BuildZoneEditorGeometry(
	request dtos.ZoneEditorGeometryRequestDto) models.ZoneEditorGeometry {
	return this.zoneEditorHandler.BuildZoneEditorGeometry(request)
}

func (this *GUIHandler) HitTestZoneEditorNode(request dtos.ZoneEditorHitTestRequestDto) string {
	return this.zoneEditorHandler.HitTestZoneEditorNode(request)
}

func (this *GUIHandler) HitTestZoneEditorEdge(position image.Point, edges []models.ZoneEditorEdge) int {
	return this.zoneEditorHandler.HitTestZoneEditorEdge(position, edges)
}

func (this *GUIHandler) GetZoneEditorGridStep(zoneRadius int) float64 {
	return this.zoneEditorHandler.GetZoneEditorGridStep(zoneRadius)
}

func (this *GUIHandler) SnapZoneEditorPosition(
	request dtos.ZoneEditorSnapRequestDto) models.ZoneEditorSnapResult {
	return this.zoneEditorHandler.SnapZoneEditorPosition(request)
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
