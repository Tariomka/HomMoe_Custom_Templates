package guiHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type handlerDependenciesStub struct {
	templateWorkflowCalled    bool
	statePersistenceCalled    bool
	templatePersistenceCalled bool
	previewCalled             bool
	contentRuleCalled         bool
	zoneEditorCalled          bool
	bonusCalled               bool
	pickerCalled              bool
}

func (this *handlerDependenciesStub) GenerateTemplate(
	editor_state_dto.EditorStateDto,
) (dtos.TemplateLoadDto, error) {
	this.templateWorkflowCalled = true
	return dtos.TemplateLoadDto{}, nil
}

func (this *handlerDependenciesStub) UpdateTemplate(
	dtos.TemplateUpdateDto,
) (dtos.TemplateLoadDto, error) {
	return dtos.TemplateLoadDto{}, nil
}

func (this *handlerDependenciesStub) ReapplyCastleSettings(
	dtos.CastleSettingsReapplyRequestDto,
) []entities.Zone {
	return nil
}

func (this *handlerDependenciesStub) ValidateEditorState(
	stateDto editor_state_model.EditorState,
	_ bool,
) editor_state_dto.EditorStateValidationDto {
	return editor_state_dto.EditorStateValidationDto{State: stateDto}
}

func (this *handlerDependenciesStub) LoadState(
	_ string,
	_ bool,
) (*editor_state_dto.EditorStateDto, []string, error) {
	return nil, nil, nil
}

func (this *handlerDependenciesStub) SaveState(editor_state_dto.EditorStateSaveDto) (string, error) {
	this.statePersistenceCalled = true
	return "", nil
}

func (this *handlerDependenciesStub) SaveTemplate(dtos.TemplateSaveDto) (string, error) {
	this.templatePersistenceCalled = true
	return "", nil
}

func (this *handlerDependenciesStub) BuildPreviewLayout(
	dtos.PreviewLayoutRequestDto,
) (dtos.PreviewLayoutDto, error) {
	this.previewCalled = true
	return dtos.PreviewLayoutDto{}, nil
}

func (this *handlerDependenciesStub) GetContentRuleEditorOptions(
	models.SidMapping,
) dtos.ContentRuleEditorOptionsDto {
	this.contentRuleCalled = true
	return dtos.ContentRuleEditorOptionsDto{}
}

func (this *handlerDependenciesStub) DescribeContentRule(
	models.SidMapping,
	models.ContentRuleRow,
) dtos.ContentRuleDescriptionDto {
	return dtos.ContentRuleDescriptionDto{}
}

func (this *handlerDependenciesStub) GetZoneEditorOptions(
	editor_state_dto.EditorStateDto,
	int,
) dtos.ZoneEditorOptionsDto {
	return dtos.ZoneEditorOptionsDto{}
}

func (this *handlerDependenciesStub) CountZoneCastles(entities.Zone) int {
	this.zoneEditorCalled = true
	return 0
}

func (this *handlerDependenciesStub) GetZoneQuality(entities.Zone) neutral_zone.Quality {
	return neutral_zone.QualityUnknown
}

func (this *handlerDependenciesStub) GetZoneConnectionGuardQuality(
	string,
	string,
	[]entities.Zone,
	map[string]bool,
) neutral_zone.Quality {
	return neutral_zone.QualityUnknown
}

func (this *handlerDependenciesStub) ApplyZoneEditorQuality(
	request dtos.ZoneEditorQualityRequestDto,
) entities.Zone {
	return request.Zone
}

func (this *handlerDependenciesStub) DescribeZoneEditorGraph(
	[]entities.Zone,
	[]entities.Connection,
) dtos.ZoneEditorGraphDto {
	return dtos.ZoneEditorGraphDto{}
}

func (this *handlerDependenciesStub) ComputeHasErrors(
	[]entities.Zone,
	[]entities.Connection,
) bool {
	return false
}

func (this *handlerDependenciesStub) RebuildZoneConnectionRoads(
	[]entities.Zone,
	[]entities.Connection,
) {
}

func (this *handlerDependenciesStub) CreateZoneEditorConnection(
	dtos.ZoneEditorConnectionRequestDto,
) entities.Connection {
	return entities.Connection{}
}

func (this *handlerDependenciesStub) FindOpenZonePosition([][2]float64) [2]float64 {
	return [2]float64{}
}

func (this *handlerDependenciesStub) GetNextZoneLabel([]entities.Zone) string {
	return ""
}

func (this *handlerDependenciesStub) CreateZoneEditorNeutralZone(
	dtos.ZoneEditorNeutralZoneRequestDto,
) entities.Zone {
	return entities.Zone{}
}

func (this *handlerDependenciesStub) CanDeleteZone(string, map[string]bool) bool {
	return false
}

func (this *handlerDependenciesStub) RemoveZoneEditorZone(
	dtos.ZoneEditorRemoveRequestDto,
) dtos.ZoneEditorMutationDto {
	return dtos.ZoneEditorMutationDto{}
}

func (this *handlerDependenciesStub) BuildZoneEditorGeometry(
	dtos.ZoneEditorGeometryRequestDto,
) models.ZoneEditorGeometry {
	return models.ZoneEditorGeometry{}
}

func (this *handlerDependenciesStub) HitTestZoneEditorNode(dtos.ZoneEditorHitTestRequestDto) string {
	return ""
}

func (this *handlerDependenciesStub) HitTestZoneEditorEdge(models.Position, []models.ZoneEditorEdge) int {
	return -1
}

func (this *handlerDependenciesStub) GetZoneEditorGridStep(float64) float64 {
	return 0
}

func (this *handlerDependenciesStub) SnapZoneEditorPosition(
	request dtos.ZoneEditorSnapRequestDto,
) models.ZoneEditorSnapResult {
	return models.ZoneEditorSnapResult{Position: request.Position}
}

func (this *handlerDependenciesStub) DescribeExistingBonuses([]config.BonusEntry) dtos.ExistingBonusesDto {
	return dtos.ExistingBonusesDto{}
}

func (this *handlerDependenciesStub) BuildBonusEntries(
	dtos.BonusCompositionRequestDto,
) dtos.BonusCompositionResultDto {
	return dtos.BonusCompositionResultDto{}
}

func (this *handlerDependenciesStub) FilterNewBonusEntries(
	entries []config.BonusEntry,
	_ map[string]bool,
) []config.BonusEntry {
	return entries
}

func (this *handlerDependenciesStub) GetSpellCountLabel(int) string {
	this.bonusCalled = true
	return ""
}

func (this *handlerDependenciesStub) ComposeContentRule(
	dtos.ContentRuleCompositionRequestDto,
) dtos.ContentRuleCompositionResultDto {
	this.contentRuleCalled = true
	return dtos.ContentRuleCompositionResultDto{}
}

func (this *handlerDependenciesStub) UpsertContentRule(
	rules []models.ContentRuleRow,
	rule models.ContentRuleRow,
) []models.ContentRuleRow {
	this.contentRuleCalled = true
	return append(rules, rule)
}

func (this *handlerDependenciesStub) GetDefaultContentRules(models.SidMapping) []models.ContentRuleRow {
	this.contentRuleCalled = true
	return nil
}

func (this *handlerDependenciesStub) GetContentRuleMarkers(
	models.SidMapping,
	[]models.ContentRuleRow,
) string {
	this.contentRuleCalled = true
	return ""
}

func (this *handlerDependenciesStub) GetContentRowDisplayName(
	content models.SidMapping,
	_ []models.ContentRuleRow,
) string {
	this.contentRuleCalled = true
	return content.Name
}

func (this *handlerDependenciesStub) SortContentItemsByName(items []models.SidMapping) []models.SidMapping {
	this.contentRuleCalled = true
	return items
}

func (this *handlerDependenciesStub) ClampContentCount(count int, _ int) int {
	this.contentRuleCalled = true
	return count
}

func (this *handlerDependenciesStub) BuildItemPickerEntries(
	[]dtos.PickerItemDto,
) []dtos.PickerEntryDto {
	this.pickerCalled = true
	return nil
}

func (this *handlerDependenciesStub) BuildSpellPickerEntries(
	[]dtos.PickerSpellDto,
) []dtos.PickerEntryDto {
	this.pickerCalled = true
	return nil
}

func (this *handlerDependenciesStub) BuildValueOverridePickerEntries([]string) []dtos.PickerEntryDto {
	this.pickerCalled = true
	return nil
}

func (this *handlerDependenciesStub) NormalizePickerFilter(text string) string {
	this.pickerCalled = true
	return text
}

func (this *handlerDependenciesStub) GetVisiblePickerRows(
	[]dtos.PickerEntryDto,
	string,
	bool,
) []dtos.PickerRowDto {
	this.pickerCalled = true
	return nil
}

func (this *handlerDependenciesStub) GetSelectedPickerIDs(
	[]dtos.PickerEntryDto,
	map[string]bool,
) []string {
	this.pickerCalled = true
	return nil
}

func (this *handlerDependenciesStub) newHandler() handler_interfaces.IGuiHandler {
	return handlers.NewGuiHandler(this, this, this, this, this, this, this)
}
