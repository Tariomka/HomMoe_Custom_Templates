package tournament_variant

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
)

type RingClusterService struct {
	base.TopologyBase
}

func NewRingClusterService() *RingClusterService {
	return &RingClusterService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *RingClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans models.NeutralZonePlans,
	playerIndex int,
	playerLabel string) ([]entities.Zone, []entities.Connection) {
	ringLabels := this.createLabels(playerNeutralZonePlans, playerLabel)
	ringCount := len(ringLabels)
	if ringCount < 2 {
		return []entities.Zone{this.createSinglePlayerZone(configuration, playerLabel, playerIndex, tuning)}, nil
	}

	connNames := make([]string, ringCount)
	for index, label := range ringLabels {
		nextIndex := (index + 1) % ringCount
		connNames[index] = fmt.Sprintf("TRing-%s-%s", label, ringLabels[nextIndex])
	}

	zones := this.createZones(configuration, ringLabels, connNames, tuning, allNeutralZonePlans, playerIndex)
	connections := this.createConnections(ringLabels, connNames, tuning, allNeutralZonePlans, playerLabel)
	return zones, connections
}

func (this *RingClusterService) createLabels(playerNeutralZonePlans models.NeutralZonePlans, playerLabel string) []string {
	sortedNeutralZonePlans := models.NeutralZonePlans{}
	sortedNeutralZonePlans.AddPlans(playerNeutralZonePlans...)
	sortedNeutralZonePlans.SortByBalanceScoreAscending()

	zoneCount := len(sortedNeutralZonePlans)
	orderedNeutralZonePlans := make(models.NeutralZonePlans, zoneCount)
	lowIndex, highIndex := 0, zoneCount-1
	for index, zonePlan := range sortedNeutralZonePlans {
		if index%2 == 0 {
			orderedNeutralZonePlans[lowIndex] = zonePlan
			lowIndex++
		} else {
			orderedNeutralZonePlans[highIndex] = zonePlan
			highIndex--
		}
	}

	return append([]string{playerLabel},
		linq.FromSlice(orderedNeutralZonePlans).
			SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
			ToSlice()...)
}

func (this *RingClusterService) createZones(
	configuration config.GeneratorConfig,
	ringLabels, connectionNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	playerIndex int) []entities.Zone {
	count := len(ringLabels)

	var zones []entities.Zone
	for index, label := range ringLabels {
		prev := (index - 1 + count) % count
		seen := map[string]bool{}
		var myConns []string
		for _, name := range []string{connectionNames[prev], connectionNames[index]} {
			if !seen[name] {
				seen[name] = true
				myConns = append(myConns, name)
			}
		}
		if index == 0 {
			zones = append(zones, this.CreateSpawnZone(
				label, fmt.Sprintf("Player%d", playerIndex+1), myConns, configuration.ZoneConfiguration.PlayerZoneCastles,
				configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
				tuning.RemoteFootholdCount, configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones, this.CreateNeutralZone(
				linq.FromSlice(allNeutralZonePlans).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
				myConns, configuration.ZoneConfiguration.Advanced.NeutralZoneSize, tuning.RemoteFootholdCount,
				configuration.GenerateRoads, tuning, false))
		}
	}
	return zones
}

func (this *RingClusterService) createSinglePlayerZone(
	configuration config.GeneratorConfig,
	playerLabel string,
	playerIndex int,
	tuning models.GenerationTuning) entities.Zone {
	return this.CreateSpawnZone(
		playerLabel, fmt.Sprintf("Player%d", playerIndex+1), nil, configuration.ZoneConfiguration.PlayerZoneCastles,
		configuration.MatchPlayerCastleFactions, configuration.ZoneConfiguration.Advanced.PlayerZoneSize,
		tuning.RemoteFootholdCount, configuration.GenerateRoads, tuning)
}

func (this *RingClusterService) createConnections(
	ringLabels, connectionNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	playerLabel string) []entities.Connection {
	ringCount := len(ringLabels)

	var connections []entities.Connection
	for currentIndex := range ringCount {
		nextIndex := (currentIndex + 1) % ringCount
		labelFrom := ringLabels[currentIndex]
		labelTo := ringLabels[nextIndex]

		connectionBuilder := variant_content.NewConnectionBuilder().
			WithName(connectionNames[currentIndex]).
			WithConnectionTypeDirect().
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(labelFrom, labelTo, []string{playerLabel}, allNeutralZonePlans, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("tourney_ring_guard_%s_%s", labelFrom, labelTo))

		if currentIndex != 0 {
			zoneFrom := "Neutral-" + labelFrom
			connectionBuilder.WithFrom(zoneFrom).WithGuardZone(zoneFrom)
		} else {
			zoneFrom := "Spawn-" + labelFrom
			connectionBuilder.WithFrom(zoneFrom).WithGuardZone(zoneFrom)
		}

		if nextIndex != 0 {
			zoneTo := "Neutral-" + labelTo
			connectionBuilder.WithTo(zoneTo)
		} else {
			zoneTo := "Spawn-" + labelTo
			connectionBuilder.WithTo(zoneTo)
		}

		connections = append(connections, connectionBuilder.Build())
	}
	return connections
}
