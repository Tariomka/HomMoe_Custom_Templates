package base

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type TopologyBase struct {
	ZoneLabelProvider zones.IZoneLabelProvider
	roadFactory       *zones.RoadFactory
	zoneFactory       *zones.ZoneFactory
	connectionService *TopologyConnectionService
}

func NewTopologyBase(
	zoneFactory *zones.ZoneFactory,
	roadFactory *zones.RoadFactory,
	zoneLabelProvider zones.IZoneLabelProvider,
	connectionService *TopologyConnectionService) TopologyBase {
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

func (this *TopologyBase) CreateSpawnZone(
	label, playerName string,
	connectionNames []string,
	castleCount int,
	matchFactions bool,
	zoneSize float64,
	footholdCount int,
	generateRoads bool,
	tuning models.GenerationTuning) entities.Zone {
	return this.zoneFactory.CreateSpawnZone(
		label,
		playerName,
		connectionNames,
		castleCount,
		matchFactions,
		zoneSize,
		footholdCount,
		generateRoads,
		tuning,
	)
}

func (this *TopologyBase) CreateNeutralZone(
	plan neutral_zone.Plan,
	connectionNames []string,
	zoneSize float64,
	footholdCount int,
	generateRoads bool,
	tuning models.GenerationTuning,
	isHoldCity bool) entities.Zone {
	return this.zoneFactory.CreateNeutralZone(models.NeutralZoneCreation{
		Name:                 constants.NeutralZonePrefix + plan.Label,
		Quality:              plan.Quality,
		Size:                 zoneSize,
		ConnectionNames:      connectionNames,
		MandatoryContentName: "mandatory_content_neutral_" + plan.Label,
		CastleCount:          plan.CastleCount,
		HoldCity:             isHoldCity,
		OutpostCount:         tuning.AbandonedOutpostCount,
		FootholdCount:        footholdCount,
		GuardRandomization:   tuning.GuardRandomization,
		GenerateRoads:        generateRoads,
		Tuning:               tuning,
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
	return this.zoneFactory.CreateHubZone(models.HubZoneCreation{
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
	maxCount int) []entities.Connection {
	return this.connectionService.CreateRandomPortalConnections(
		playerLabels,
		orderedLabels,
		tuning,
		maxCount,
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
