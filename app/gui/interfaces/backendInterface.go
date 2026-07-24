package interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

type IBackend interface {
	ITemplateWorkflowHandler
	IStatePersistenceHandler
	IStateValidationHandler
	IPreviewHandler
	IContentRuleHandler
	IZoneEditorHandler
}

type ITemplateWorkflowHandler interface {
	GenerateTemplate(stateDto dtos.EditorStateDto) (dtos.TemplateLoadDto, error)
	UpdateTemplate(templateDto dtos.TemplateUpdateDto) (dtos.TemplateLoadDto, error)
	ReapplyCastleSettings(request dtos.CastleSettingsReapplyRequestDto) []entities.Zone
	SaveTemplate(templateDto dtos.TemplateSaveDto) (string, error)
}

type IStatePersistenceHandler interface {
	LoadState(path string, fixIssues bool) (*dtos.EditorStateDto, []string, error)
	SaveState(stateDto dtos.EditorStateSaveDto) (string, error)
}

type IStateValidationHandler interface {
	ValidateEditorState(state dtos.EditorStateDto, fixIssues bool) dtos.EditorStateValidationDto
}

type IPreviewHandler interface {
	BuildPreviewLayout(request dtos.PreviewLayoutRequestDto) (dtos.PreviewLayoutDto, error)
}

type IContentRuleHandler interface {
	GetContentRuleEditorOptions(content models.SidMapping) dtos.ContentRuleEditorOptionsDto
	DescribeContentRule(
		content models.SidMapping,
		savedRule models.ContentRuleRowSave,
	) dtos.ContentRuleDescriptionDto
}

type IZoneEditorHandler interface {
	GetZoneEditorOptions(state dtos.EditorStateDto, totalZoneCount int) dtos.ZoneEditorOptionsDto
	CountZoneCastles(zone entities.Zone) int
	ApplyZoneEditorQuality(request dtos.ZoneEditorQualityRequestDto) entities.Zone
	DescribeZoneEditorGraph(zones []entities.Zone, connections []entities.Connection) dtos.ZoneEditorGraphDto
	CreateZoneEditorConnection(request dtos.ZoneEditorConnectionRequestDto) entities.Connection
	FindOpenZonePosition(occupied [][2]float64) [2]float64
	GetNextZoneLabel(zones []entities.Zone) string
	CreateZoneEditorNeutralZone(request dtos.ZoneEditorNeutralZoneRequestDto) entities.Zone
	CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool
	RemoveZoneEditorZone(request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto
}
