package base

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type TopologyBase struct {
	ZoneLabelProvider zone_interfaces.IZoneLabelProvider
	roadFactory       zone_interfaces.IRoadFactory
	zoneFactory       zone_interfaces.IZoneFactory
	connectionService ITopologyConnectionService
}

func NewTopologyBase(
	zoneFactory zone_interfaces.IZoneFactory,
	roadFactory zone_interfaces.IRoadFactory,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	connectionService ITopologyConnectionService) TopologyBase {
	return TopologyBase{
		ZoneLabelProvider: zoneLabelProvider,
		roadFactory:       roadFactory,
		zoneFactory:       zoneFactory,
		connectionService: connectionService,
	}
}

func (this *TopologyBase) CreateVariant(
	playerLabels []string,
	firstLabel string,
	zoneCount int,
	zones []entities.Zone,
	connections []entities.Connection) entities.Variant {
	orientationBuilder := variant_content.NewOrientationBuilder().
		WithBaseAngleMin(45).
		WithBaseAngleMax(45).
		WithRandomAngleAmplitude(360)

	if zoneCount > 0 {
		orientationBuilder.WithRandomAngleStep(360 / zoneCount)
	}

	if slices.Contains(playerLabels, firstLabel) {
		orientationBuilder.WithZeroAngleZone(constants.PlayerZonePrefix + firstLabel)
	} else {
		orientationBuilder.WithZeroAngleZone(constants.NeutralZonePrefix + firstLabel)
	}

	return variant_content.NewVariantBuilder().
		WithOrientation(orientationBuilder.Build()).
		WithBorder(variant_content.NewBorderBuilder().
			WithCornerRadius(0).
			WithObstaclesWidth(3).
			WithObstaclesNoise(1, 12).
			WithWaterWidth(0).
			WithWaterNoise(1, 12).
			WithWaterTypeWaterGrass().
			Build()).
		WithZones(zones...).
		WithConnections(connections...).
		Build()
}

// CreateClusterZone builds either the player's spawn zone or a neutral zone for
// the same label, which is the choice every topology service makes per zone.
func (this *TopologyBase) CreateClusterZone(
	configuration config.GeneratorConfig,
	label string,
	connectionNames []string,
	playerIndex int,
	isSpawn bool,
	isHoldCity bool,
	tuning models.GenerationTuning,
	allNeutralZonePlans neutral_zone.Plans) entities.Zone {
	if isSpawn {
		return this.CreateSpawnZone(models.SpawnZoneCreationRequest{
			Label:           label,
			PlayerName:      fmt.Sprintf("Player%d", playerIndex+1),
			ConnectionNames: connectionNames,
			CastleCount:     configuration.ZoneConfiguration.PlayerZoneCastles,
			MatchFactions:   configuration.MatchPlayerCastleFactions,
			Size:            configuration.ZoneConfiguration.PlayerZoneSize,
			FootholdCount:   tuning.RemoteFootholdCount,
			GenerateRoads:   configuration.GenerateRoads,
			Tuning:          tuning,
		})
	}

	return this.CreateNeutralZone(models.TopologyNeutralZoneCreationRequest{
		Plan: linq.FromSlice(allNeutralZonePlans).
			FirstOrDefault(func(plan neutral_zone.Plan) bool { return plan.Label == label }),
		ConnectionNames: connectionNames,
		Size:            configuration.ZoneConfiguration.NeutralZoneSize,
		FootholdCount:   tuning.RemoteFootholdCount,
		GenerateRoads:   configuration.GenerateRoads,
		HoldCity:        isHoldCity,
		Tuning:          tuning,
	})
}

func (this *TopologyBase) CreateSpawnZone(input models.SpawnZoneCreationRequest) entities.Zone {
	return this.zoneFactory.CreateSpawnZone(input)
}

func (this *TopologyBase) CreateNeutralZone(input models.TopologyNeutralZoneCreationRequest) entities.Zone {
	return this.zoneFactory.CreateNeutralZone(models.NeutralZoneCreationRequest{
		Name:                 constants.NeutralZonePrefix + input.Plan.Label,
		Quality:              input.Plan.Quality,
		Size:                 input.Size,
		ConnectionNames:      input.ConnectionNames,
		MandatoryContentName: "mandatory_content_neutral_" + input.Plan.Label,
		CastleCount:          input.Plan.CastleCount,
		HoldCity:             input.HoldCity,
		OutpostCount:         input.Tuning.AbandonedOutpostCount,
		FootholdCount:        input.FootholdCount,
		GuardRandomization:   input.Tuning.GuardRandomization,
		GenerateRoads:        input.GenerateRoads,
		Tuning:               input.Tuning,
	})
}

func (this *TopologyBase) CreateHubZone(
	name string,
	connectionNames []string,
	tuning models.GenerationTuning,
	isHoldCity bool,
	size float64,
	castleCount int,
	generateRoads bool,
	mandatoryContentName string) entities.Zone {
	return this.zoneFactory.CreateHubZone(models.HubZoneCreationRequest{
		Name:                 name,
		Size:                 size,
		ConnectionNames:      connectionNames,
		MandatoryContentName: mandatoryContentName,
		CastleCount:          castleCount,
		HoldCity:             isHoldCity,
		GuardRandomization:   0.05,
		GenerateRoads:        generateRoads,
		Tuning:               tuning,
	})
}

func (this *TopologyBase) CreateRandomPortalConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	maxCount int,
	neutralZones neutral_zone.Plans) []entities.Connection {
	return this.connectionService.CreateRandomPortalConnections(
		playerLabels,
		orderedLabels,
		tuning,
		maxCount,
		neutralZones,
	)
}

func (this *TopologyBase) CreateMissingPlayerConnections(
	playerLabels []string,
	zones []entities.Zone,
	connections []entities.Connection,
	tuning models.GenerationTuning) []entities.Connection {
	return this.connectionService.CreateMissingPlayerConnections(playerLabels, zones, connections, tuning)
}

func (this *TopologyBase) CreateMissingConnections(
	playerLabels, allLabels []string,
	positions models.Positions,
	zones []entities.Zone,
	connections []entities.Connection,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans) []entities.Connection {
	return this.connectionService.CreateMissingConnections(
		playerLabels,
		allLabels,
		positions,
		zones,
		connections,
		tuning,
		neutralZones,
	)
}

func (this *TopologyBase) CreateConnectorZoneRoads(connectionNames []string, generateRoads bool) []entities.Road {
	return this.roadFactory.CreateConnectorZoneRoads(connectionNames, generateRoads)
}

func (this *TopologyBase) GetBorderGuardValue(
	labelA, labelB string,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning) int {
	return this.connectionService.GetBorderGuardValue(labelA, labelB, playerLabels, neutralZones, tuning)
}
