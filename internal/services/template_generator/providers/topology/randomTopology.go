package topology

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
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
	neutralZones neutralZone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string) entities.Variant {
	return this.createVariantFromLayout(
		configuration, playerLabels, neutralZones, tuning, holdCityNeutralLabel, this.createRandomLayout)
}

// createRandomLayout scatters all zones uniformly over the map and connects
// them through a Delaunay triangulation of the random positions.
func (this *RandomTopologyService) createRandomLayout(
	playerLabels []string,
	neutralZones neutralZone.Plans) ([]string, models.Positions, []models.ConnectionIndexes) {
	neutralLabels := make([]string, len(neutralZones))
	for i, nz := range neutralZones {
		neutralLabels[i] = nz.Label
	}
	allLabels := append(append([]string{}, playerLabels...), neutralLabels...)
	labelCount := len(allLabels)
	rand.Shuffle(labelCount, func(i, j int) { allLabels[i], allLabels[j] = allLabels[j], allLabels[i] })
	var positions models.Positions
	for range labelCount {
		positions.Add(data.NewVec2(rand.Float64()*0.9+0.05, rand.Float64()*0.9+0.05))
	}
	pairs := positions.CreateDelaunayTriangulation()
	return allLabels, positions, pairs
}

// createVariantFromLayout is the shared topology-variant pipeline: build the
// layout, derive connection names, create the zones with their generator
// positions, create the connections (plus portal and missing-connection
// fill-ins), and assemble the variant. Each topology service supplies only its
// own layout builder.
func (this *RandomTopologyService) createVariantFromLayout(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones neutralZone.Plans,
	tuning models.GenerationTuning,
	holdCityNeutralLabel string,
	buildLayout layoutFunc) entities.Variant {
	isIsolated := configuration.NoDirectPlayerConnections && len(playerLabels) > 1
	allLabels, positions, pairs := buildLayout(playerLabels, neutralZones)
	connectionNames := this.createConnectionNames(playerLabels, allLabels, pairs, isIsolated)

	zones := this.createZones(
		configuration, playerLabels, allLabels, tuning, neutralZones, holdCityNeutralLabel, connectionNames)
	for index := range zones {
		position := positions[index]
		zones[index].GeneratorPosition = &[2]float64{position.X, position.Y}
	}

	conns := this.createConnections(playerLabels, allLabels, tuning, isIsolated, neutralZones, connectionNames, pairs)
	if configuration.RandomPortals {
		conns = append(conns,
			this.CreateRandomPortalConnections(playerLabels, allLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if isIsolated {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLabels, zones, conns, tuning)...)
	}
	conns = append(conns,
		this.CreateMissingConnections(playerLabels, allLabels, positions, zones, conns, tuning, neutralZones)...)
	return this.CreateVariant(playerLabels, allLabels[0], len(allLabels), zones, conns)
}

func (this *RandomTopologyService) createConnectionNames(
	playerLabels, allLabels []string,
	triangulationPairs []models.ConnectionIndexes,
	isIsolated bool) map[int][]string {
	connectionNamesByZone := make(map[int][]string, len(allLabels))
	for _, pair := range triangulationPairs {
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

func (this *RandomTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, allLabels []string,
	tuning models.GenerationTuning,
	neutralZones neutralZone.Plans,
	holdCityNeutralLabel string,
	connectionNames map[int][]string) []entities.Zone {
	var zones []entities.Zone
	for index, label := range allLabels {
		myConns := connectionNames[index]
		if playerIndex := slices.Index(playerLabels, label); playerIndex >= 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), myConns,
					configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
					configuration.ZoneConfiguration.PlayerZoneSize, tuning.RemoteFootholdCount,
					configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(neutralZones).
						FirstOrDefault(func(x neutralZone.Plan) bool { return x.Label == label }),
					myConns, configuration.ZoneConfiguration.NeutralZoneSize, tuning.RemoteFootholdCount,
					configuration.GenerateRoads, tuning, label == holdCityNeutralLabel))
		}
	}
	return zones
}

func (this *RandomTopologyService) createConnections(
	playerLabels, allLabels []string,
	tuning models.GenerationTuning,
	isIsolated bool,
	neutralZones neutralZone.Plans,
	connectionNames map[int][]string,
	triangulationPairs []models.ConnectionIndexes) []entities.Connection {
	nameLookup := make(map[int]int, len(allLabels))

	var connections []entities.Connection
	for _, pair := range triangulationPairs {
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
