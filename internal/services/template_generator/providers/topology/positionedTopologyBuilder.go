package topology

import (
	"fmt"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type PositionedTopologyBuilder struct {
	base.TopologyBase
}

func NewPositionedTopologyBuilder(
	creationServices *zone_services.CreationServices,
) *PositionedTopologyBuilder {
	return &PositionedTopologyBuilder{
		TopologyBase: base.NewTopologyBaseWithCreationServices(creationServices),
	}
}

func (this *PositionedTopologyBuilder) BuildVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string,
	buildLayout PositionedTopologyLayoutBuilder,
	decorateZones PositionedTopologyZoneDecorator,
) entities.Variant {
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1
	allLabels, positions, pairs := buildLayout(playerLabels, neutralZones)
	connectionNames := this.createConnectionNames(playerLabels, allLabels, pairs, isIsolated)

	zones := this.createZones(
		configuration, playerLabels, allLabels, tuning, neutralZones, holdCityNeutralLabel, connectionNames)
	for index := range zones {
		position := positions[index]
		zones[index].GeneratorPosition = &[2]float64{position.X, position.Y}
	}
	if decorateZones != nil {
		decorateZones(zones, allLabels, playerLabels, neutralZones)
	}

	connections := this.createConnections(
		playerLabels, allLabels, tuning, isIsolated, neutralZones, connectionNames, pairs)
	if configuration.RandomPortals {
		connections = append(connections,
			this.CreateRandomPortalConnections(playerLabels, allLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isIsolated {
		connections = append(connections,
			this.CreateMissingPlayerConnections(playerLabels, zones, connections, tuning)...)
	}
	connections = append(connections,
		this.CreateMissingConnections(
			playerLabels, allLabels, positions, zones, connections, tuning, neutralZones)...)
	return this.CreateVariant(playerLabels, allLabels[0], len(allLabels), zones, connections)
}

func (this *PositionedTopologyBuilder) createConnectionNames(
	playerLabels, allLabels []string,
	pairs []models.ConnectionIndexes,
	isIsolated bool,
) map[int][]string {
	connectionNamesByZone := make(map[int][]string, len(allLabels))
	for _, pair := range pairs {
		indexA, indexB := pair.X, pair.Y
		labelFrom := allLabels[indexA]
		labelTo := allLabels[indexB]
		if isIsolated && slices.Contains(playerLabels, labelFrom) && slices.Contains(playerLabels, labelTo) {
			continue
		}

		connectionName := fmt.Sprintf("Rnd-%s-%s", labelFrom, labelTo)
		connectionNamesByZone[indexA] = append(connectionNamesByZone[indexA], connectionName)
		connectionNamesByZone[indexB] = append(connectionNamesByZone[indexB], connectionName)
	}

	return connectionNamesByZone
}

func (this *PositionedTopologyBuilder) createZones(
	configuration config.GeneratorConfig,
	playerLabels, allLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans,
	holdCityNeutralLabel string,
	connectionNames map[int][]string,
) []entities.Zone {
	var zones []entities.Zone
	for index, label := range allLabels {
		zoneConnectionNames := connectionNames[index]
		if playerIndex := slices.Index(playerLabels, label); playerIndex >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), zoneConnectionNames,
					configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
					configuration.ZoneConfiguration.PlayerZoneSize, tuning.RemoteFootholdCount,
					configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).
						FirstOrDefault(func(plan neutral_zone.Plan) bool { return plan.Label == label }),
					zoneConnectionNames, configuration.ZoneConfiguration.NeutralZoneSize, tuning.RemoteFootholdCount,
					configuration.GenerateRoads, tuning, label == holdCityNeutralLabel))
		}
	}
	return zones
}

func (this *PositionedTopologyBuilder) createConnections(
	playerLabels, allLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones neutral_zone.Plans,
	connectionNames map[int][]string,
	pairs []models.ConnectionIndexes,
) []entities.Connection {
	nameLookup := make(map[int]int, len(allLabels))

	var connections []entities.Connection
	for _, pair := range pairs {
		indexA, indexB := pair.X, pair.Y
		labelFrom := allLabels[indexA]
		labelTo := allLabels[indexB]
		if isIsolated && slices.Contains(playerLabels, labelFrom) && slices.Contains(playerLabels, labelTo) {
			continue
		}
		connectionName := connectionNames[indexA][nameLookup[indexA]]
		zoneNameFrom := this.ZoneLabelProvider.CreateZoneName(labelFrom, playerLabels)
		nameLookup[indexA]++
		nameLookup[indexB]++

		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(connectionName).
			WithFrom(zoneNameFrom).
			WithTo(this.ZoneLabelProvider.CreateZoneName(labelTo, playerLabels)).
			WithConnectionTypeDirect().
			WithGuardZone(zoneNameFrom).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(labelFrom, labelTo, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("rnd_guard_%s_%s", labelFrom, labelTo)).
			Build())
	}
	return connections
}
