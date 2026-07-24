package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/mock"
)

// TemplateHandlerMock is a testify mock of interfaces.Backend, used
// to unit-test app/gui/drivers.State without the real generator stack.
type TemplateHandlerMock struct {
	mock.Mock

	ValidateEditorStateFunc         func(dtos.EditorStateDto, bool) dtos.EditorStateValidationDto
	BuildPreviewLayoutFunc          func(dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error)
	GetContentRuleEditorOptionsFunc func(models.SidMapping) dtos.ContentRuleEditorOptionsDto
	DescribeContentRuleFunc         func(models.SidMapping, models.ContentRuleRowSave) dtos.ContentRuleDescriptionDto
	ReapplyCastleSettingsFunc       func(dtos.CastleSettingsReapplyRequestDto) []entities.Zone
	GetZoneEditorOptionsFunc        func(dtos.EditorStateDto, int) dtos.ZoneEditorOptionsDto
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
}

func (this *TemplateHandlerMock) GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error) {
	arguments := this.Called(stateDto)
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
	state dtos.EditorStateDto,
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

func (this *TemplateHandlerMock) SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error) {
	arguments := this.Called(templateDto)
	return arguments.String(0), arguments.Error(1)
}

func (this *TemplateHandlerMock) LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error) {
	arguments := this.Called(path, fixIssues)
	state, _ := arguments.Get(0).(*dtos.EditorStateDto)
	warnings, _ := arguments.Get(1).([]string)
	return state, warnings, arguments.Error(2)
}

func (this *TemplateHandlerMock) SaveState(stateDto dtos.EditorStateSaveDto) (string, error) {
	arguments := this.Called(stateDto)
	return arguments.String(0), arguments.Error(1)
}

func (this *TemplateHandlerMock) ValidateEditorState(
	stateDto dtos.EditorStateDto,
	fixIssues bool,
) dtos.EditorStateValidationDto {
	if this.ValidateEditorStateFunc != nil {
		return this.ValidateEditorStateFunc(stateDto, fixIssues)
	}
	return dtos.EditorStateValidationDto{State: stateDto}
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
	savedRule models.ContentRuleRowSave,
) dtos.ContentRuleDescriptionDto {
	if this.DescribeContentRuleFunc != nil {
		return this.DescribeContentRuleFunc(content, savedRule)
	}
	return dtos.ContentRuleDescriptionDto{DisplayText: savedRule.Name, SavedRule: savedRule}
}
