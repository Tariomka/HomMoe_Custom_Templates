package topology

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

type RandomTopologyService struct {
	base.TopologyBase
}

func NewRandomTopologyService() *RandomTopologyService {
	return &RandomTopologyService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *RandomTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) template.Variant {
	neutralLabels := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLabels[i] = nz.Label
	}
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1
	allLabels := append(append([]string{}, playerLabels...), neutralLabels...)
	labelCount := len(allLabels)
	rand.Shuffle(labelCount, func(i, j int) { allLabels[i], allLabels[j] = allLabels[j], allLabels[i] })
	var positions models.Positions
	for i := 0; i < labelCount; i++ {
		positions.Add(models.NewPosition(rand.Float64()*0.9+0.05, rand.Float64()*0.9+0.05))
	}
	pairs := positions.CreateDelaunayTriangulation()
	connectionNames := this.createConnectionNames(playerLabels, allLabels, pairs, isIsolated)

	zones := this.createZones(configuration, playerLabels, allLabels, tuning, neutralZones, holdCityNeutralLabel, connectionNames)
	for index := range zones {
		position := positions[index]
		zones[index].GeneratorPosition = &[2]float64{position.X, position.Y}
	}

	conns := this.createConnections(playerLabels, allLabels, tuning, isIsolated, neutralZones, connectionNames, pairs)
	if configuration.RandomPortals {
		conns = append(conns, this.CreateRandomPortalConnections(playerLabels, allLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isIsolated {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLabels, zones, conns, tuning)...)
	}
	conns = this.CreateMissingConnections(playerLabels, allLabels, positions, zones, conns, tuning, neutralZones)
	return this.CreateVariant(playerLabels, allLabels[0], labelCount, zones, conns)
}

func (this *RandomTopologyService) createConnectionNames(
	playerLabels, allLabels []string,
	triangulationPairs [][2]int,
	isIsolated bool) map[int][]string {
	connectionNamesByZone := make(map[int][]string, len(allLabels))
	for _, pair := range triangulationPairs {
		indexA, indexB := pair[0], pair[1]
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

func (this *RandomTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, allLabels []string,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans,
	holdCityNeutralLabel string,
	connectionNames map[int][]string) []template.Zone {
	var zones []template.Zone
	for index, label := range allLabels {
		myConns := connectionNames[index]
		if playerIndex := slices.Index(playerLabels, label); playerIndex >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles,
					configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, label == holdCityNeutralLabel))
		}
	}
	return zones
}

func (this *RandomTopologyService) createConnections(
	playerLabels, allLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones models.NeutralZonePlans,
	connectionNames map[int][]string,
	triangulationPairs [][2]int) []template.Connection {
	nameLookup := make(map[int]int, len(allLabels))

	var connections []template.Connection
	for _, pair := range triangulationPairs {
		indexA, indexB := pair[0], pair[1]
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
