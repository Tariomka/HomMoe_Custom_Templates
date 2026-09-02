package handlers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type zoneEditorHandler struct {
	mapper           mappers.IGeneratorConfigMapper
	tierService      zone_interfaces.IZoneTierService
	connectionEditor connection_editor.IConnectionEditorService
	zoneEditor       connection_editor.IZoneEditorService
	geometry         connection_editor.IZoneEditorGeometryService
	tuningFactory    generation_tuning.IGenerationTuningFactory
}

func NewZoneEditorHandler(
	mapper mappers.IGeneratorConfigMapper,
	tierService zone_interfaces.IZoneTierService,
	connectionEditor connection_editor.IConnectionEditorService,
	zoneEditor connection_editor.IZoneEditorService,
	geometry connection_editor.IZoneEditorGeometryService,
	tuningFactory generation_tuning.IGenerationTuningFactory) handler_interfaces.IZoneEditorHandler {
	return &zoneEditorHandler{
		mapper:           mapper,
		tierService:      tierService,
		connectionEditor: connectionEditor,
		zoneEditor:       zoneEditor,
		geometry:         geometry,
		tuningFactory:    tuningFactory,
	}
}

func (this *zoneEditorHandler) GetZoneEditorOptions(
	state editor_state_dto.EditorStateDto,
	totalZoneCount int) dtos.ZoneEditorOptionsDto {
	configuration := this.mapper.FromEditorState(state.EditorState)
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
	return this.tierService.GetQuality(zone)
}

func (this *zoneEditorHandler) GetZoneConnectionGuardQuality(
	from, to string,
	zones []entities.Zone,
	playerZoneNames map[string]bool) neutral_zone.Quality {
	playerNames := make([]string, 0, len(playerZoneNames))
	for playerName := range playerZoneNames {
		playerNames = append(playerNames, playerName)
	}
	return this.tierService.GetConnectionGuardQuality(from, to, zones, playerNames)
}

func (this *zoneEditorHandler) ApplyZoneEditorQuality(request dtos.ZoneEditorQualityRequestDto) entities.Zone {
	this.zoneEditor.ApplyNeutralZoneQuality(
		&request.Zone,
		request.Quality,
		request.CastleCount,
		request.Tuning)
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

func (this *zoneEditorHandler) CreateZoneEditorConnection(
	request dtos.ZoneEditorConnectionRequestDto) entities.Connection {
	return this.connectionEditor.NewDefaultConnection(request.From, request.To, request.Zones, request.PlayerZoneNames)
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
		request.Tuning)
}

func (this *zoneEditorHandler) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return this.zoneEditor.CanDeleteZone(zoneName, playerZoneNames)
}

func (this *zoneEditorHandler) RemoveZoneEditorZone(
	request dtos.ZoneEditorRemoveRequestDto) dtos.ZoneEditorMutationDto {
	zones, connections := this.zoneEditor.RemoveZone(request.Zones, request.Connections, request.ZoneName)
	return dtos.ZoneEditorMutationDto{Zones: zones, Connections: connections}
}

func (this *zoneEditorHandler) BuildZoneEditorGeometry(
	request dtos.ZoneEditorGeometryRequestDto) models.ZoneEditorGeometry {
	return this.geometry.BuildGeometry(request.Zones, request.Connections, request.Topology, request.CanvasSide)
}

func (this *zoneEditorHandler) HitTestZoneEditorNode(request dtos.ZoneEditorHitTestRequestDto) string {
	return this.geometry.HitTestNode(request.Position, request.Positions, request.ZoneRadius)
}

func (this *zoneEditorHandler) HitTestZoneEditorEdge(position models.Position, edges []models.ZoneEditorEdge) int {
	return this.geometry.HitTestEdge(position, edges)
}

func (this *zoneEditorHandler) GetZoneEditorGridStep(zoneRadius float64) float64 {
	return this.geometry.GridStep(zoneRadius)
}

func (this *zoneEditorHandler) SnapZoneEditorPosition(
	request dtos.ZoneEditorSnapRequestDto) models.ZoneEditorSnapResult {
	return this.geometry.SnapPosition(request.Position, request.Positions, request.ZoneRadius, request.DraggedZone)
}
