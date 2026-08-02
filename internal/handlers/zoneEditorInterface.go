package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type IZoneEditor interface {
	GetZoneEditorOptions(state dtos.EditorStateDto, totalZoneCount int) dtos.ZoneEditorOptionsDto
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
}
