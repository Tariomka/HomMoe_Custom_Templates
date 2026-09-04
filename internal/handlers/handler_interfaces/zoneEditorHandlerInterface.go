package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IZoneEditorHandler interface {
	GetZoneEditorOptions(state editor_state_dto.EditorStateDto, totalZoneCount int) dtos.ZoneEditorOptionsDto
	CountZoneCastles(zone template_model.Zone) int
	GetZoneQuality(zone template_model.Zone) neutral_zone.Quality
	GetZoneConnectionGuardQuality(
		from, to string,
		zones []template_model.Zone,
		playerZoneNames map[string]bool) neutral_zone.Quality
	ApplyZoneEditorQuality(request dtos.ZoneEditorQualityRequestDto) template_model.Zone
	DescribeZoneEditorGraph(
		zones []template_model.Zone,
		connections []template_model.Connection) dtos.ZoneEditorGraphDto
	CreateZoneEditorConnection(request dtos.ZoneEditorConnectionRequestDto) template_model.Connection
	FindOpenZonePosition(occupied [][2]float64) [2]float64
	GetNextZoneLabel(zones []template_model.Zone) string
	CreateZoneEditorNeutralZone(request dtos.ZoneEditorNeutralZoneRequestDto) template_model.Zone
	CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool
	RemoveZoneEditorZone(request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto
	BuildZoneEditorGeometry(request dtos.ZoneEditorGeometryRequestDto) models.ZoneEditorGeometry
	HitTestZoneEditorNode(request dtos.ZoneEditorHitTestRequestDto) string
	HitTestZoneEditorEdge(position models.Position, edges []models.ZoneEditorEdge) int
	GetZoneEditorGridStep(zoneRadius float64) float64
	SnapZoneEditorPosition(request dtos.ZoneEditorSnapRequestDto) models.ZoneEditorSnapResult
}
