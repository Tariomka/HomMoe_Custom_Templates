package tournament_variant

import (
	"fmt"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/builders/variant_content"
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
	playerLabel string) ([]template.Zone, []template.Connection) {
	hubName := "Hub-" + playerLabel
	spokeLabels := append([]string{playerLabel},
		linq.FromSlice(playerNeutralZonePlans).
			SelectString(func(x models.NeutralZonePlan) string { return x.Label }).
			ToSlice()...)
	spokeConnNames := make([]string, len(spokeLabels))
	for index, label := range spokeLabels {
		spokeConnNames[index] = fmt.Sprintf("THubSpoke-%s-%s", playerLabel, label)
	}

	zones := this.createZones(configuration, spokeLabels, spokeConnNames, tuning, allNeutralZonePlans, hubName, playerIndex)
	connections := this.createConnections(spokeLabels, spokeConnNames, tuning, allNeutralZonePlans, hubName, playerLabel)
	return zones, connections
}

func (this *HubClusterService) createZones(
	configuration config.GeneratorConfig,
	spokeLabels, spokeConnNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	hubName string,
	playerIndex int) []template.Zone {
	var zones []template.Zone
	hubZone := this.CreateHubZone(
		spokeConnNames, tuning, false, configuration.ZoneConfiguration.HubZoneSize,
		configuration.ZoneConfiguration.HubZoneCastles, configuration.GenerateRoads)
	hubZone.Name = hubName
	zones = append(zones, hubZone)

	for index, label := range spokeLabels {
		connectionNames := []string{spokeConnNames[index]}
		if index == 0 {
			zones = append(zones,
				this.CreateSpawnZone(
					label, fmt.Sprintf("Player%d", playerIndex+1), connectionNames,
					configuration.ZoneConfiguration.PlayerZoneCastles, configuration.MatchPlayerCastleFactions,
					configuration.ZoneConfiguration.Advanced.PlayerZoneSize, configuration.SpawnRemoteFootholds,
					configuration.GenerateRoads, tuning))
		} else {
			zones = append(zones,
				this.CreateNeutralZone(
					linq.FromSlice(allNeutralZonePlans).FirstOrDefault(func(x models.NeutralZonePlan) bool { return x.Label == label }),
					connectionNames, configuration.ZoneConfiguration.Advanced.NeutralZoneSize,
					configuration.SpawnRemoteFootholds, configuration.GenerateRoads, tuning, false))
		}
	}
	return zones
}

func (this *HubClusterService) createConnections(
	spokeLabels, spokeConnNames []string,
	tuning models.GenerationTuning,
	allNeutralZonePlans models.NeutralZonePlans,
	hubName, playerLabel string) []template.Connection {
	var connections []template.Connection
	for index, spokeLabel := range spokeLabels {
		spokeZone := "Spawn-" + spokeLabel
		if index != 0 {
			spokeZone = "Neutral-" + spokeLabel
		}
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(spokeConnNames[index]).
			WithFrom(hubName).
			WithTo(spokeZone).
			WithConnectionTypeDirect().
			WithGuardZone(hubName).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(playerLabel, spokeLabel, []string{playerLabel}, allNeutralZonePlans, tuning)).
			WithGuardWeeklyIncrement(0.15).
			WithGuardMatchGroup(fmt.Sprintf("tourney_hub_guard_%s_%s", playerLabel, spokeLabel)).
			Build())
	}

	// Proximity ring around spokes so the engine arranges them sensibly.
	for currentIndex, label := range spokeLabels {
		nextIndex := (currentIndex + 1) % len(spokeLabels)
		labelTo := spokeLabels[nextIndex]
		zoneFrom := "Spawn-" + label
		if currentIndex != 0 {
			zoneFrom = "Neutral-" + label
		}
		zoneTo := "Spawn-" + labelTo
		if nextIndex != 0 {
			zoneTo = "Neutral-" + labelTo
		}
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(fmt.Sprintf("TPseudo-%s-%s-%s", playerLabel, label, labelTo)).
			WithFrom(zoneFrom).
			WithTo(zoneTo).
			WithConnectionTypeProximity().
			Build())
	}
	return connections
}
