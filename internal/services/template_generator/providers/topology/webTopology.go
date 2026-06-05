package topology

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
)

type SharedWebTopologyService struct {
	topologyBase
}

func NewSharedWebTopologyService() *SharedWebTopologyService {
	return &SharedWebTopologyService{
		topologyBase: newTopologyBase(),
	}
}

func (this *SharedWebTopologyService) CreateTopologyVariant(
	configuration config.GeneratorConfig,
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	tuning models.GenerationTuning,
	holdCityNeutralLetter string) template.Variant {
	neutralLabels := this.createLabels(playerLabels, neutralZones, configuration.Topology == config.TopologyBalanced)
	playerSpokes, neutralSpokes := this.createSpokes(playerLabels, neutralLabels)
	neutralConnNames := this.createRingConnectionNames(neutralLabels)

	zones := this.createZones(configuration, playerLabels, neutralLabels, tuning, neutralZones, holdCityNeutralLetter, playerSpokes, neutralSpokes, neutralConnNames)
	conns := this.createConnections(playerLabels, neutralLabels, tuning, neutralZones, playerSpokes, neutralConnNames)
	if configuration.RandomPortals {
		allLabels := append(append([]string{}, playerLabels...), neutralLabels...)
		conns = append(conns, this.CreateRandomPortalConnections(playerLabels, allLabels, tuning, configuration.MaxPortalConnections)...)
	}
	if configuration.NoDirectPlayerConnections && len(playerLabels) > 1 {
		conns = append(conns, this.CreateMissingPlayerConnections(playerLabels, zones, conns, tuning)...)
	}
	return this.CreateVariant(playerLabels, playerLabels[0], len(zones), zones, conns)
}

func (this *SharedWebTopologyService) createLabels(
	playerLabels []string,
	neutralZones models.NeutralZonePlans,
	isBalanced bool) []string {
	var neutrals []string
	if isBalanced {
		neutrals = this.zoneLabelProvider.CreateBalancedNeutralRingZoneLabels(neutralZones, len(playerLabels))
	} else {
		for _, zonePlan := range neutralZones {
			neutrals = append(neutrals, zonePlan.Label)
		}
	}
	return neutrals
}

func (this *SharedWebTopologyService) createSpokes(
	playerLabels, neutralLabels []string) (playerSpokes, neutralSpokes map[string][]string) {
	playerCount := len(playerLabels)
	neutralCount := len(neutralLabels)

	addSpoke := func(playerLabel, neutralLabel string) {
		connectionName := fmt.Sprintf("Web-%s-%s", playerLabel, neutralLabel)
		playerSpokes[playerLabel] = append(playerSpokes[playerLabel], connectionName)
		neutralSpokes[neutralLabel] = append(neutralSpokes[neutralLabel], connectionName)
	}

	for _, label := range playerLabels {
		playerSpokes[label] = nil
	}
	for _, label := range neutralLabels {
		neutralSpokes[label] = nil
	}

	for index, label := range playerLabels {
		n1 := (index * neutralCount / playerCount) % neutralCount
		n2 := ((index * neutralCount / playerCount) + 1) % neutralCount
		addSpoke(label, neutralLabels[n1])
		if n1 != n2 {
			addSpoke(label, neutralLabels[n2])
		}
	}

	return
}

func (this *SharedWebTopologyService) createRingConnectionNames(neutralLabels []string) []string {
	neutralCount := len(neutralLabels)
	neutralRingConnNames := make([]string, neutralCount)
	for i, label := range neutralLabels {
		neutralRingConnNames[i] = fmt.Sprintf("NRing-%s-%s", label, neutralLabels[(i+1)%neutralCount])
	}
	return neutralRingConnNames
}

func (this *SharedWebTopologyService) createZones(
	configuration config.GeneratorConfig,
	playerLabels, neutralLabels []string,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans,
	holdCityNeutralLabel string,
	playerSpokes, neutralSpokes map[string][]string,
	connectionNames []string) []template.Zone {
	neutralCount := len(neutralLabels)

	var zones []template.Zone
	for i, label := range neutralLabels {
		var neutralConnNames []string
		if neutralCount > 1 {
			neutralConnNames = append(neutralConnNames,
				connectionNames[(i-1+neutralCount)%neutralCount],
				connectionNames[i])
		}
		neutralConnNames = linq.FromSlice(append(neutralConnNames, neutralSpokes[label]...)).Distinct().ToSlice()
		zonePlan := linq.FromSlice(neutralZones).
			FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label })
		zone := this.CreateNeutralZone(
			zonePlan, neutralConnNames, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
			configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning,
			label == holdCityNeutralLabel)
		if zonePlan.CastleCount == 0 {
			zone.Roads = this.CreateConnectorZoneRoads(neutralConnNames, configuration.GenerateRoads)
		}
		zones = append(zones, zone)
	}

	for i, label := range playerLabels {
		zones = append(zones,
			this.CreateSpawnZone(
				label, fmt.Sprintf("Player%d", i+1), playerSpokes[label], configuration.ZoneConfiguration.PlayerZoneCastles,
				configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
				configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning))
	}

	return zones
}

func (this *SharedWebTopologyService) createConnections(
	playerLabels, neutralLabels []string,
	tuning models.GenerationTuning,
	neutralZones models.NeutralZonePlans,
	playerSpokes map[string][]string,
	connectionNames []string) []template.Connection {
	neutralCount := len(neutralLabels)

	var connections []template.Connection
	for _, label := range playerLabels {
		for _, connectionName := range playerSpokes[label] {
			nextLabel := strings.Split(connectionName, "-")[2]
			connections = append(connections, variant_content.NewConnectionBuilder().
				WithName(connectionName).
				WithFrom("Spawn-"+label).
				WithTo("Neutral-"+nextLabel).
				WithConnectionTypeDirect().
				WithGuardZone("Neutral-"+nextLabel).
				WithSimTurnSquad().
				WithGuardValue(this.GetBorderGuardValue(label, nextLabel, playerLabels, neutralZones, tuning)).
				WithGuardWeeklyIncrement(0.15).
				WithGuardMatchGroup(fmt.Sprintf("web_guard_%s_%s", label, nextLabel)).
				Build())
		}
	}

	if neutralCount < 2 {
		return connections
	}

	for i, label := range neutralLabels {
		next := (i + 1) % neutralCount
		nextLabel := neutralLabels[next]
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(connectionNames[i]).
			WithFrom("Neutral-"+label).
			WithTo("Neutral-"+nextLabel).
			WithConnectionTypeDirect().
			WithGuardZone("Neutral-"+label).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(label, nextLabel, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("nring_guard_%s_%s", label, nextLabel)).
			Build())
	}
	return connections
}
