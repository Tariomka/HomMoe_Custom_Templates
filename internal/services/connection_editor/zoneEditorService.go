package connection_editor

import (
	"fmt"
	"math"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type ZoneEditorService struct {
	castleFactory zone_interfaces.ICastleFactory
	zoneFactory   zone_interfaces.IZoneFactory
	roadPolicy    zone_interfaces.IRoadPolicyService
}

func NewZoneEditorService(
	castleFactory zone_interfaces.ICastleFactory,
	roadPolicy zone_interfaces.IRoadPolicyService,
	zoneFactory zone_interfaces.IZoneFactory) IZoneEditorService {
	return &ZoneEditorService{
		zoneFactory:   zoneFactory,
		castleFactory: castleFactory,
		roadPolicy:    roadPolicy,
	}
}

func (this *ZoneEditorService) EnsureConnectionNames(connections []template_model.Connection) {
	used := make(map[string]bool, len(connections))
	for _, connection := range connections {
		if connection.Name != "" {
			used[connection.Name] = true
		}
	}
	for i := range connections {
		if connections[i].Name != "" {
			continue
		}

		prefix := constants.GetManualConnectionNameFor(
			helpers.GetZoneLabel(connections[i].From),
			helpers.GetZoneLabel(connections[i].To))
		name := prefix
		for suffix := 2; used[name]; suffix++ {
			name = fmt.Sprintf("%s-%d", prefix, suffix)
		}
		connections[i].Name = name
		used[name] = true
	}
}

func (this *ZoneEditorService) RebuildZoneConnectionRoads(
	zones []template_model.Zone,
	connections []template_model.Connection) {
	this.EnsureConnectionNames(connections)
	this.roadPolicy.RebuildZoneConnectionRoads(zones, connections)
}

func (this *ZoneEditorService) ApplyConnectionRoadPolicy(
	connection *template_model.Connection,
	generateRoads bool) {
	this.roadPolicy.StampConnectionRoad(connection, generateRoads)
}

func (this *ZoneEditorService) ChangeConnectionType(
	connection *template_model.Connection,
	connectionType string,
	generateRoads bool) {
	connection.ConnectionType = connectionType
	this.roadPolicy.StampConnectionRoad(connection, generateRoads)
}

func (this *ZoneEditorService) NextFreeZoneLabel(zones []template_model.Zone) string {
	used := make(map[string]bool, len(zones))
	for _, zone := range zones {
		used[helpers.GetZoneLabel(zone.Name)] = true
	}
	for _, label := range constants.GetZoneLabels() {
		if !used[label] {
			return label
		}
	}

	return ""
}

func (this *ZoneEditorService) NewDefaultNeutralZone(
	label string,
	quality neutral_zone.Quality,
	castleCount int,
	generateRoads bool,
	tuning models.GenerationTuning) template_model.Zone {
	return this.zoneFactory.CreateNeutralZone(models.NeutralZoneCreationRequest{
		Name:               constants.GetNeutralZoneNameFor(label),
		Quality:            quality,
		Size:               1.0,
		CastleCount:        castleCount,
		OutpostCount:       tuning.AbandonedOutpostCount,
		FootholdCount:      tuning.RemoteFootholdCount,
		GuardRandomization: tuning.GuardRandomization,
		GenerateRoads:      generateRoads,
		Tuning:             tuning,
	})
}

func (this *ZoneEditorService) CountZoneCastles(zone template_model.Zone) int {
	count := 0
	for _, mainObject := range zone.MainObjects {
		if mainObject.Type == registry.GetMainObjectTypeValues().City {
			count++
		}
	}
	return count
}

func (this *ZoneEditorService) ApplyNeutralZoneQuality(
	zone *template_model.Zone,
	quality neutral_zone.Quality,
	castleCount int,
	tuning models.GenerationTuning) {
	profile := common_zones.GetNeutralZoneProfile(quality)
	zone.Quality = &quality
	zone.Layout = profile.Layout
	zone.GuardMultiplier = tuning.ScaleByNeutralGuardStrengthPrecise(profile.GuardMultiplier)
	zone.GuardReactionDistribution = profile.GuardReactionDistribution
	zone.GuardedContentPool = profile.GuardedContentPool
	zone.UnguardedContentPool = profile.UnguardedContentPool
	zone.ResourcesContentPool = profile.ResourcesContentPool
	zone.GuardedContentValue = tuning.ScaleByStructureDensity(
		float64(profile.GuardedContentValue) * tuning.ContentScale)
	zone.GuardedContentValuePerArea = tuning.ScaleByStructureDensity(
		float64(profile.GuardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))
	zone.UnguardedContentValue = tuning.ScaleByStructureDensity(
		float64(profile.UnguardedContentValue) * tuning.ContentScale)
	zone.UnguardedContentValuePerArea = tuning.ScaleByStructureDensity(
		float64(profile.UnguardedContentValuePerArea) * math.Sqrt(tuning.ContentScale))
	zone.ResourcesValue = tuning.ScaleByResourceDensity(float64(profile.ResourcesValue) * tuning.ContentScale)
	zone.ResourcesValuePerArea = tuning.ScaleByResourceDensity(
		float64(profile.ResourcesValuePerArea) * math.Sqrt(tuning.ContentScale))
	zone.MainObjects = this.castleFactory.CreateNeutralZoneCastles(profile, tuning, castleCount, false)

	// Regenerate the castle<->castle roads so the rebuilt castles are
	// road-connected. Other roads (connection and foothold roads) are left for
	// RebuildZoneConnectionRoads to finalize once the edit is applied.
	this.RebuildCastleRoads(zone)
}

// CanDeleteZone reports whether the zone may be removed in the editor. Spawn
// zones are owned by the General tab's player count and cannot be deleted.
func (this *ZoneEditorService) CanDeleteZone(zoneName string, playerZoneNames map[string]bool) bool {
	return !playerZoneNames[zoneName]
}

// RemoveZone returns the zone and connection lists without the named zone and
// without any connection referencing it.
func (this *ZoneEditorService) RemoveZone(
	zones []template_model.Zone,
	connections []template_model.Connection,
	zoneName string) ([]template_model.Zone, []template_model.Connection) {
	keptZones := make([]template_model.Zone, 0, len(zones))
	for _, zone := range zones {
		if zone.Name != zoneName {
			keptZones = append(keptZones, zone)
		}
	}
	keptConnections := make([]template_model.Connection, 0, len(connections))
	for _, connection := range connections {
		if connection.From != zoneName && connection.To != zoneName {
			keptConnections = append(keptConnections, connection)
		}
	}
	return keptZones, keptConnections
}

// FindOpenPosition returns a normalized position on a coarse interior grid
// that maximizes the distance to the occupied positions.
func (this *ZoneEditorService) FindOpenPosition(occupied []data.Vec2[float64]) data.Vec2[float64] {
	const gridSteps = 7
	best := data.NewVec2(0.5, 0.5)
	bestScore := -1.0
	for row := range gridSteps {
		for col := range gridSteps {
			candidate := data.NewVec2(
				0.1+0.8*float64(col)/float64(gridSteps-1),
				0.1+0.8*float64(row)/float64(gridSteps-1))
			score := math.MaxFloat64
			for _, position := range occupied {
				score = math.Min(score, candidate.Subtract(position).Distance())
			}
			if len(occupied) == 0 {
				score = candidate.Subtract(data.NewVec2(0.5, 0.5)).Distance()
			}
			if score > bestScore {
				bestScore = score
				best = candidate
			}
		}
	}
	return best
}

// RebuildCastleRoads reconciles the zone's castle<->castle roads with its
// current main objects, preserving every other road.
func (this *ZoneEditorService) RebuildCastleRoads(zone *template_model.Zone) {
	this.roadPolicy.RebuildCastleRoads(zone)
}
