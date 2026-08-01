package guiHandler_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type handlerDependenciesStub struct {
	templateWorkflowCalled    bool
	statePersistenceCalled    bool
	templatePersistenceCalled bool
	previewCalled             bool
	contentRuleCalled         bool
	zoneEditorCalled          bool
}

func (this *handlerDependenciesStub) GenerateTemplate(
	dtos.EditorStateDto,
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
	stateDto dtos.EditorStateDto,
	_ bool,
) dtos.EditorStateValidationDto {
	return dtos.EditorStateValidationDto{State: stateDto}
}

func (this *handlerDependenciesStub) LoadState(
	_ string,
	_ bool,
) (*dtos.EditorStateDto, []string, error) {
	return nil, nil, nil
}

func (this *handlerDependenciesStub) SaveState(dtos.EditorStateSaveDto) (string, error) {
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
	models.ContentRuleRowSave,
) dtos.ContentRuleDescriptionDto {
	return dtos.ContentRuleDescriptionDto{}
}

func (this *handlerDependenciesStub) GetZoneEditorOptions(
	dtos.EditorStateDto,
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

func (this *handlerDependenciesStub) dependencies() handlers.GUIHandlerDependencies {
	return handlers.GUIHandlerDependencies{
		TemplateWorkflow:    this,
		StatePersistence:    this,
		TemplatePersistence: this,
		Preview:             this,
		ContentRule:         this,
		ZoneEditor:          this,
	}
}
