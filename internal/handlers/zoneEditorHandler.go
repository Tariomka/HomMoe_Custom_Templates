package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type zoneEditorHandler struct {
	mapper           *mappers.GeneratorConfigMapper
	zoneClassifier   *zone_services.ZoneClassifier
	connectionEditor *connection_editor.ConnectionEditorService
	zoneEditor       *connection_editor.ZoneEditorService
	tuningFactory    *generation_tuning.GenerationTuningFactory
}

func newZoneEditorHandler(
	mapper *mappers.GeneratorConfigMapper,
	zoneClassifier *zone_services.ZoneClassifier,
	connectionEditor *connection_editor.ConnectionEditorService,
	zoneEditor *connection_editor.ZoneEditorService,
	tuningFactory *generation_tuning.GenerationTuningFactory) handler_interfaces.IZoneEditorHandler {
	return &zoneEditorHandler{
		mapper:           mapper,
		zoneClassifier:   zoneClassifier,
		connectionEditor: connectionEditor,
		zoneEditor:       zoneEditor,
		tuningFactory:    tuningFactory,
	}
}

func (this *zoneEditorHandler) GetZoneEditorOptions(
	state dtos.EditorStateDto,
	totalZoneCount int) dtos.ZoneEditorOptionsDto {
	configuration := this.mapper.FromEditorState(state)
	return dtos.ZoneEditorOptionsDto{
		Topology:      state.Topology,
		Tuning:        this.tuningFactory.Create(configuration, totalZoneCount),
		GenerateRoads: state.GenerateRoads,
	}
}

func (this *zoneEditorHandler) CountZoneCastles(zone entities.Zone) int {
	return this.zoneEditor.CountZoneCastles(zone)
}

func (this *zoneEditorHandler) GetZoneQuality(zone entities.Zone) neutral_zone.Quality {
	return this.zoneClassifier.GetQuality(zone)
}

func (this *zoneEditorHandler) GetZoneConnectionGuardQuality(
	from, to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool) neutral_zone.Quality {
	playerNames := make([]string, 0, len(playerZoneNames))
	for playerName := range playerZoneNames {
		playerNames = append(playerNames, playerName)
	}
	return this.zoneClassifier.GetConnectionGuardQuality(from, to, zones, playerNames)
}

func (this *zoneEditorHandler) ApplyZoneEditorQuality(
	request dtos.ZoneEditorQualityRequestDto) entities.Zone {
	this.zoneEditor.ApplyNeutralZoneQuality(
		&request.Zone,
		request.Quality,
		request.CastleCount,
		request.Tuning,
	)
	return request.Zone
}

func (this *zoneEditorHandler) DescribeZoneEditorGraph(
	zones []entities.Zone,
	connections []entities.Connection) dtos.ZoneEditorGraphDto {
	return dtos.ZoneEditorGraphDto{
		HasErrors:         this.connectionEditor.ComputeHasErrors(zones, connections),
		IsolatedZoneCount: len(this.connectionEditor.FindIsolatedZones(zones, connections)),
	}
}

func (this *zoneEditorHandler) ComputeHasErrors(
	zones []entities.Zone,
	connections []entities.Connection) bool {
	return this.connectionEditor.ComputeHasErrors(zones, connections)
}

func (this *zoneEditorHandler) CreateZoneEditorConnection(
	request dtos.ZoneEditorConnectionRequestDto) entities.Connection {
	return this.connectionEditor.NewDefaultConnection(
		request.From,
		request.To,
		request.Zones,
		request.PlayerZoneNames,
	)
}

func (this *zoneEditorHandler) FindOpenZonePosition(occupied [][2]float64) [2]float64 {
	return this.zoneEditor.FindOpenPosition(occupied)
}

func (this *zoneEditorHandler) GetNextZoneLabel(zones []entities.Zone) string {
	return this.zoneEditor.NextFreeZoneLabel(zones)
}

func (this *zoneEditorHandler) CreateZoneEditorNeutralZone(request dtos.ZoneEditorNeutralZoneRequestDto) entities.Zone {
	return this.zoneEditor.NewDefaultNeutralZone(
		request.Label,
		request.Quality,
		request.CastleCount,
		request.GenerateRoads,
		request.Tuning,
	)
}

func (this *zoneEditorHandler) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return this.zoneEditor.CanDeleteZone(zoneName, playerZoneNames)
}

func (this *zoneEditorHandler) RemoveZoneEditorZone(
	request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto {
	zones, connections := this.zoneEditor.RemoveZone(
		request.Zones,
		request.Connections,
		request.ZoneName,
	)
	return dtos.ZoneEditorMutationDto{Zones: zones, Connections: connections}
}

func (this *zoneEditorHandler) RebuildZoneConnectionRoads(
	zones []entities.Zone,
	connections []entities.Connection) {
	this.zoneEditor.RebuildZoneConnectionRoads(zones, connections)
}
