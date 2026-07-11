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

type HubClusterService struct {
	base.TopologyBase
}

func NewHubClusterService() *HubClusterService {
	return &HubClusterService{
		TopologyBase: base.NewTopologyBase(),
	}
}

func (this *HubClusterService) CreateClusterVariant(
	configuration config.GeneratorConfig,
	tuning models.GenerationTuning,
	allNeutralZonePlans, playerNeutralZonePlans models.NeutralZonePlans,
	playerIndex int,
	playerLabel string) ([]entities.Zone, []entities.Connection) {
	hubName := "Hub-" + playerLabel
	spokeLabels := append([]string{playerLabel},
		linq.FromSlice(playerNeutralZonePlans).
			SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
			ToSlice()...)
	spokeConnNames := make([]string, len(spokeLabels))
	for index, label := range spokeLabels {
		spokeConnNames[index] = fmt.Sprintf("THubSpoke-%s-%s", playerLabel, label)
	}

	zones := this.createZones(
		configuration, spokeLabels, spokeConnNames, tuning, allNeutralZonePlans, hubName, playerIndex)
	connections := this.createConnections(
		spokeLabels, spokeConnNames, tuning, allNeutralZonePlans, hubName, playerLabel)
	return zones, connections
}

func (this *HubClusterService) createZones(
	configuration config.GeneratorConfig,
	spokeLabels, spokeConnNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	hubName string,
	playerIndex int) []entities.Zone {
	var zones []entities.Zone
	hubContentName := ""
	if len(configuration.HubZoneMandatoryContent) > 0 {
		hubContentName = "mandatory_content_hub"
	}
	hubZone := this.CreateHubZone(
		spokeConnNames, tuning, false, configuration.ZoneConfiguration.HubZoneSize,
		configuration.ZoneConfiguration.HubZoneCastles, configuration.GenerateRoads, hubContentName)
	hubZone.Name = hubName
	zones = append(zones, hubZone)

	for index, label := range spokeLabels {
		connectionNames := []string{spokeConnNames[index]}
		if index == 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), connectionNames,
					configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
					configuration.ZoneConfiguration.Advanced.PlayerZoneSize, tuning.RemoteFootholdCount,
					configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(allNeutralZonePlans).
						FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					connectionNames, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					tuning.RemoteFootholdCount, configuration.GenerateRoads, tuning, false))
		}
	}
	return zones
}

func (this *HubClusterService) createConnections(
	spokeLabels, spokeConnNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	hubName, playerLabel string) []entities.Connection {
	var connections []entities.Connection
	for index, spokeLabel := range spokeLabels {
		connectionBuilder := variant_content.NewConnectionBuilder().
			WithName(spokeConnNames[index]).
			WithFrom(hubName).
			WithConnectionTypeDirect().
			WithGuardZone(hubName).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(playerLabel, spokeLabel, []string{playerLabel}, allNeutralZonePlans, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("tourney_hub_guard_%s_%s", playerLabel, spokeLabel))

		if index != 0 {
			spokeZone := "Neutral-" + spokeLabel
			connectionBuilder = connectionBuilder.WithTo(spokeZone).WithGuardZone(spokeZone)
		} else {
			spokeZone := "Spawn-" + spokeLabel
			connectionBuilder = connectionBuilder.WithTo(spokeZone).WithGuardZone(hubName)
		}

		connections = append(connections, connectionBuilder.Build())
	}

	// Proximity ring around spokes so the engine arranges them sensibly.
	for currentIndex, label := range spokeLabels {
		nextIndex := (currentIndex + 1) % len(spokeLabels)
		labelTo := spokeLabels[nextIndex]
		connectionBuilder := variant_content.NewConnectionBuilder().
			WithName(fmt.Sprintf("THubRing-%s-%s-%s", playerLabel, label, labelTo)).
			WithConnectionTypeProximity()

		if currentIndex != 0 {
			connectionBuilder = connectionBuilder.WithFrom("Neutral-" + label)
		} else {
			connectionBuilder = connectionBuilder.WithFrom("Spawn-" + label)
		}

		if nextIndex != 0 {
			connectionBuilder = connectionBuilder.WithTo("Neutral-" + labelTo)
		} else {
			connectionBuilder = connectionBuilder.WithTo("Spawn-" + labelTo)
		}

		connections = append(connections, connectionBuilder.Build())
	}
	return connections
}
