package handler_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type IZoneEditorHandler interface {
	GetZoneEditorOptions(state editor_state_model.EditorState, totalZoneCount int) dtos.ZoneEditorOptionsDto
	CountZoneCastles(zone entities.Zone) int
	GetZoneQuality(zone entities.Zone) neutral_zone.Quality
	GetZoneConnectionGuardQuality(
		from, to string,
		zones []entities.Zone,
		playerZoneNames map[string]bool) neutral_zone.Quality
	ApplyZoneEditorQuality(request dtos.ZoneEditorQualityRequestDto) entities.Zone
	DescribeZoneEditorGraph(zones []entities.Zone, connections []entities.Connection) dtos.ZoneEditorGraphDto
	CreateZoneEditorConnection(request dtos.ZoneEditorConnectionRequestDto) entities.Connection
	FindOpenZonePosition(occupied [][2]float64) [2]float64
	GetNextZoneLabel(zones []entities.Zone) string
	CreateZoneEditorNeutralZone(request dtos.ZoneEditorNeutralZoneRequestDto) entities.Zone
	CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool
	RemoveZoneEditorZone(request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto
	BuildZoneEditorGeometry(request dtos.ZoneEditorGeometryRequestDto) models.ZoneEditorGeometry
	HitTestZoneEditorNode(request dtos.ZoneEditorHitTestRequestDto) string
	HitTestZoneEditorEdge(position models.Position, edges []models.ZoneEditorEdge) int
	GetZoneEditorGridStep(zoneRadius float64) float64
	SnapZoneEditorPosition(request dtos.ZoneEditorSnapRequestDto) models.ZoneEditorSnapResult
}
