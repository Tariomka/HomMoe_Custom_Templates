package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/stretchr/testify/mock"
)

// TemplateHandlerMock is a testify mock of interfaces.IBackend, used
// to unit-test app/gui/drivers.State without the real generator stack.
type TemplateHandlerMock struct {
	mock.Mock

	ValidateEditorStateFunc         func(editor_state_model.EditorState, bool) editor_state_dto.EditorStateValidationDto
	BuildPreviewLayoutFunc          func(dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error)
	GetContentRuleEditorOptionsFunc func(models.SidMapping) dtos.ContentRuleEditorOptionsDto
	DescribeContentRuleFunc         func(models.SidMapping, models.ContentRuleRow) dtos.ContentRuleDescriptionDto
	ReapplyCastleSettingsFunc       func(dtos.CastleSettingsReapplyRequestDto) []entities.Zone
	GetZoneEditorOptionsFunc        func(editor_state_model.EditorState, int) dtos.ZoneEditorOptionsDto
	CountZoneCastlesFunc            func(entities.Zone) int
	GetZoneQualityFunc              func(entities.Zone) neutral_zone.Quality
	GetZoneConnectionQualityFunc    func(string, string, []entities.Zone, map[string]bool) neutral_zone.Quality
	ApplyZoneEditorQualityFunc      func(dtos.ZoneEditorQualityRequestDto) entities.Zone
	DescribeZoneEditorGraphFunc     func([]entities.Zone, []entities.Connection) dtos.ZoneEditorGraphDto
	CreateZoneEditorConnectionFunc  func(dtos.ZoneEditorConnectionRequestDto) entities.Connection
	FindOpenZonePositionFunc        func([][2]float64) [2]float64
	GetNextZoneLabelFunc            func([]entities.Zone) string
	CreateZoneEditorNeutralZoneFunc func(dtos.ZoneEditorNeutralZoneRequestDto) entities.Zone
	CanDeleteZoneFunc               func(string, map[string]bool) bool
	RemoveZoneEditorZoneFunc        func(dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto
	BuildZoneEditorGeometryFunc     func(dtos.ZoneEditorGeometryRequestDto) models.ZoneEditorGeometry
	HitTestZoneEditorNodeFunc       func(dtos.ZoneEditorHitTestRequestDto) string
	HitTestZoneEditorEdgeFunc       func(models.Position, []models.ZoneEditorEdge) int
	GetZoneEditorGridStepFunc       func(float64) float64
	SnapZoneEditorPositionFunc      func(dtos.ZoneEditorSnapRequestDto) models.ZoneEditorSnapResult
	DescribeExistingBonusesFunc     func([]config.BonusEntry) dtos.ExistingBonusesDto
	BuildBonusEntriesFunc           func(dtos.BonusCompositionRequestDto) dtos.BonusCompositionResultDto
	FilterNewBonusEntriesFunc       func([]config.BonusEntry, map[string]bool) []config.BonusEntry
	GetSpellCountLabelFunc          func(int) string
}

func (this *TemplateHandlerMock) GenerateTemplate(
	state editor_state_model.EditorState,
) (dtos.TemplateLoadDto, error) {
	arguments := this.Called(state)
	template, _ := arguments.Get(0).(dtos.TemplateLoadDto)
	return template, arguments.Error(1)
}

func (this *TemplateHandlerMock) UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error) {
	arguments := this.Called(templateDto)
	template, _ := arguments.Get(0).(dtos.TemplateLoadDto)
	return template, arguments.Error(1)
}

func (this *TemplateHandlerMock) ReapplyCastleSettings(
	request dtos.CastleSettingsReapplyRequestDto,
) []entities.Zone {
	if this.ReapplyCastleSettingsFunc != nil {
		return this.ReapplyCastleSettingsFunc(request)
	}
	return request.Zones
}

func (this *TemplateHandlerMock) GetZoneEditorOptions(
	state editor_state_model.EditorState,
	totalZoneCount int,
) dtos.ZoneEditorOptionsDto {
	if this.GetZoneEditorOptionsFunc != nil {
		return this.GetZoneEditorOptionsFunc(state, totalZoneCount)
	}
	return dtos.ZoneEditorOptionsDto{}
}

func (this *TemplateHandlerMock) CountZoneCastles(zone entities.Zone) int {
	if this.CountZoneCastlesFunc != nil {
		return this.CountZoneCastlesFunc(zone)
	}
	return 0
}

func (this *TemplateHandlerMock) GetZoneQuality(zone entities.Zone) neutral_zone.Quality {
	if this.GetZoneQualityFunc != nil {
		return this.GetZoneQualityFunc(zone)
	}
	return neutral_zone.QualityUnknown
}

func (this *TemplateHandlerMock) GetZoneConnectionGuardQuality(
	from, to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool,
) neutral_zone.Quality {
	if this.GetZoneConnectionQualityFunc != nil {
		return this.GetZoneConnectionQualityFunc(from, to, zones, playerZoneNames)
	}
	return neutral_zone.QualityUnknown
}

func (this *TemplateHandlerMock) ApplyZoneEditorQuality(
	request dtos.ZoneEditorQualityRequestDto,
) entities.Zone {
	if this.ApplyZoneEditorQualityFunc != nil {
		return this.ApplyZoneEditorQualityFunc(request)
	}
	return request.Zone
}

func (this *TemplateHandlerMock) DescribeZoneEditorGraph(
	zones []entities.Zone,
	connections []entities.Connection,
) dtos.ZoneEditorGraphDto {
	if this.DescribeZoneEditorGraphFunc != nil {
		return this.DescribeZoneEditorGraphFunc(zones, connections)
	}
	return dtos.ZoneEditorGraphDto{}
}

func (this *TemplateHandlerMock) CreateZoneEditorConnection(
	request dtos.ZoneEditorConnectionRequestDto,
) entities.Connection {
	if this.CreateZoneEditorConnectionFunc != nil {
		return this.CreateZoneEditorConnectionFunc(request)
	}
	return entities.Connection{}
}

func (this *TemplateHandlerMock) FindOpenZonePosition(occupied [][2]float64) [2]float64 {
	if this.FindOpenZonePositionFunc != nil {
		return this.FindOpenZonePositionFunc(occupied)
	}
	return [2]float64{}
}

func (this *TemplateHandlerMock) GetNextZoneLabel(zones []entities.Zone) string {
	if this.GetNextZoneLabelFunc != nil {
		return this.GetNextZoneLabelFunc(zones)
	}
	return ""
}

func (this *TemplateHandlerMock) CreateZoneEditorNeutralZone(
	request dtos.ZoneEditorNeutralZoneRequestDto,
) entities.Zone {
	if this.CreateZoneEditorNeutralZoneFunc != nil {
		return this.CreateZoneEditorNeutralZoneFunc(request)
	}
	return entities.Zone{}
}

func (this *TemplateHandlerMock) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	if this.CanDeleteZoneFunc != nil {
		return this.CanDeleteZoneFunc(zoneName, playerZoneNames)
	}
	return false
}

func (this *TemplateHandlerMock) RemoveZoneEditorZone(
	request dtos.ZoneEditorRemoveRequestDto,
) dtos.ZoneEditorMutationDto {
	if this.RemoveZoneEditorZoneFunc != nil {
		return this.RemoveZoneEditorZoneFunc(request)
	}
	return dtos.ZoneEditorMutationDto{Zones: request.Zones, Connections: request.Connections}
}

func (this *TemplateHandlerMock) BuildZoneEditorGeometry(
	request dtos.ZoneEditorGeometryRequestDto,
) models.ZoneEditorGeometry {
	if this.BuildZoneEditorGeometryFunc != nil {
		return this.BuildZoneEditorGeometryFunc(request)
	}
	return models.ZoneEditorGeometry{Positions: map[string]models.Position{}}
}

func (this *TemplateHandlerMock) HitTestZoneEditorNode(request dtos.ZoneEditorHitTestRequestDto) string {
	if this.HitTestZoneEditorNodeFunc != nil {
		return this.HitTestZoneEditorNodeFunc(request)
	}
	return ""
}

func (this *TemplateHandlerMock) HitTestZoneEditorEdge(
	position models.Position,
	edges []models.ZoneEditorEdge,
) int {
	if this.HitTestZoneEditorEdgeFunc != nil {
		return this.HitTestZoneEditorEdgeFunc(position, edges)
	}
	return -1
}

func (this *TemplateHandlerMock) GetZoneEditorGridStep(zoneRadius float64) float64 {
	if this.GetZoneEditorGridStepFunc != nil {
		return this.GetZoneEditorGridStepFunc(zoneRadius)
	}
	return 0
}

func (this *TemplateHandlerMock) SnapZoneEditorPosition(
	request dtos.ZoneEditorSnapRequestDto,
) models.ZoneEditorSnapResult {
	if this.SnapZoneEditorPositionFunc != nil {
		return this.SnapZoneEditorPositionFunc(request)
	}
	return models.ZoneEditorSnapResult{Position: request.Position}
}

func (this *TemplateHandlerMock) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	arguments := this.Called(templateDto)
	return arguments.String(0), arguments.Error(1)
}

func (this *TemplateHandlerMock) LoadState(
	path string,
	fixIssues bool,
) (*editor_state_model.EditorState, []string, error) {
	arguments := this.Called(path, fixIssues)
	state, _ := arguments.Get(0).(*editor_state_model.EditorState)
	warnings, _ := arguments.Get(1).([]string)
	return state, warnings, arguments.Error(2)
}

func (this *TemplateHandlerMock) SaveState(stateDto editor_state_dto.EditorStateSaveDto) (string, error) {
	arguments := this.Called(stateDto)
	return arguments.String(0), arguments.Error(1)
}

func (this *TemplateHandlerMock) ValidateEditorState(
	state editor_state_model.EditorState,
	fixIssues bool,
) editor_state_dto.EditorStateValidationDto {
	if this.ValidateEditorStateFunc != nil {
		return this.ValidateEditorStateFunc(state, fixIssues)
	}
	return editor_state_dto.EditorStateValidationDto{State: state}
}

func (this *TemplateHandlerMock) BuildPreviewLayout(
	request dtos.PreviewLayoutRequestDto,
) (dtos.PreviewLayoutDto, error) {
	if this.BuildPreviewLayoutFunc != nil {
		return this.BuildPreviewLayoutFunc(request)
	}
	return dtos.PreviewLayoutDto{}, nil
}

func (this *TemplateHandlerMock) GetContentRuleEditorOptions(
	content models.SidMapping,
) dtos.ContentRuleEditorOptionsDto {
	if this.GetContentRuleEditorOptionsFunc != nil {
		return this.GetContentRuleEditorOptionsFunc(content)
	}
	return dtos.ContentRuleEditorOptionsDto{}
}

func (this *TemplateHandlerMock) DescribeContentRule(
	content models.SidMapping,
	savedRule models.ContentRuleRow,
) dtos.ContentRuleDescriptionDto {
	if this.DescribeContentRuleFunc != nil {
		return this.DescribeContentRuleFunc(content, savedRule)
	}
	return dtos.ContentRuleDescriptionDto{DisplayText: savedRule.Name, SavedRule: savedRule}
}

func (this *TemplateHandlerMock) DescribeExistingBonuses(
	existing []config.BonusEntry,
) dtos.ExistingBonusesDto {
	if this.DescribeExistingBonusesFunc != nil {
		return this.DescribeExistingBonusesFunc(existing)
	}
	return dtos.ExistingBonusesDto{Keys: map[string]bool{}}
}

func (this *TemplateHandlerMock) BuildBonusEntries(
	request dtos.BonusCompositionRequestDto,
) dtos.BonusCompositionResultDto {
	if this.BuildBonusEntriesFunc != nil {
		return this.BuildBonusEntriesFunc(request)
	}
	return dtos.BonusCompositionResultDto{}
}

func (this *TemplateHandlerMock) FilterNewBonusEntries(
	entries []config.BonusEntry,
	existingKeys map[string]bool,
) []config.BonusEntry {
	if this.FilterNewBonusEntriesFunc != nil {
		return this.FilterNewBonusEntriesFunc(entries, existingKeys)
	}
	return entries
}

func (this *TemplateHandlerMock) GetSpellCountLabel(count int) string {
	if this.GetSpellCountLabelFunc != nil {
		return this.GetSpellCountLabelFunc(count)
	}
	return ""
}

func (this *TemplateHandlerMock) ComposeContentRule(
	request dtos.ContentRuleCompositionRequestDto,
) dtos.ContentRuleCompositionResultDto {
	return zone_content.NewZoneContentEditorService().ComposeContentRule(request)
}

func (this *TemplateHandlerMock) UpsertContentRule(
	rules []models.ContentRuleRow,
	rule models.ContentRuleRow,
) []models.ContentRuleRow {
	return zone_content.NewZoneContentEditorService().UpsertContentRule(rules, rule)
}

func (this *TemplateHandlerMock) GetDefaultContentRules(models.SidMapping) []models.ContentRuleRow {
	return nil
}

func (this *TemplateHandlerMock) GetContentRuleMarkers(models.SidMapping, []models.ContentRuleRow) string {
	return ""
}

func (this *TemplateHandlerMock) GetContentRowDisplayName(
	content models.SidMapping,
	_ []models.ContentRuleRow,
) string {
	return content.Name
}

func (this *TemplateHandlerMock) SortContentItemsByName(items []models.SidMapping) []models.SidMapping {
	return zone_content.NewZoneContentEditorService().SortContentItemsByName(items)
}

func (this *TemplateHandlerMock) ClampContentCount(count int, maxCount int) int {
	return zone_content.NewZoneContentEditorService().ClampContentCount(count, maxCount)
}

func (this *TemplateHandlerMock) BuildItemPickerEntries(items []dtos.PickerItemDto) []dtos.PickerEntryDto {
	return pickers.NewPickerEntryService().BuildItemPickerEntries(items)
}

func (this *TemplateHandlerMock) BuildSpellPickerEntries(spells []dtos.PickerSpellDto) []dtos.PickerEntryDto {
	return pickers.NewPickerEntryService().BuildSpellPickerEntries(spells)
}

func (this *TemplateHandlerMock) BuildValueOverridePickerEntries(sids []string) []dtos.PickerEntryDto {
	return pickers.NewPickerEntryService().BuildValueOverridePickerEntries(sids)
}

func (this *TemplateHandlerMock) NormalizePickerFilter(text string) string {
	return pickers.NewPickerEntryService().NormalizePickerFilter(text)
}

func (this *TemplateHandlerMock) GetVisiblePickerRows(
	entries []dtos.PickerEntryDto,
	filter string,
	grouped bool) []dtos.PickerRowDto {
	return pickers.NewPickerEntryService().GetVisiblePickerRows(entries, filter, grouped)
}

func (this *TemplateHandlerMock) GetSelectedPickerIDs(
	entries []dtos.PickerEntryDto,
	selected map[string]bool) []string {
	return pickers.NewPickerEntryService().GetSelectedPickerIDs(entries, selected)
}
