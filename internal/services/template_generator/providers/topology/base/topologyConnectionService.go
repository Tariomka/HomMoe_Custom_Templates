package base

import (
	"fmt"
	"math/rand/v2"
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_connections"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/geometry_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/placement_rule"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type TopologyConnectionService struct {
	zoneLabelProvider zone_interfaces.IZoneLabelProvider
}

func NewTopologyConnectionService(zoneLabelProvider zone_interfaces.IZoneLabelProvider) ITopologyConnectionService {
	return &TopologyConnectionService{zoneLabelProvider: zoneLabelProvider}
}

func (this *TopologyConnectionService) CreateRandomPortalConnections(
	playerLabels, orderedLabels []string,
	tuning models.GenerationTuning,
	maxCount int,
	neutralZones neutral_zone.Plans) []entities.Connection {
	count := len(orderedLabels)
	if count < 2 {
		return nil
	}
	destinations := buildNonAdjacentDerangement(count)
	indices := make([]int, count)
	for index := range count {
		indices[index] = index
	}
	rand.Shuffle(len(indices), func(firstIndex, secondIndex int) {
		indices[firstIndex], indices[secondIndex] = indices[secondIndex], indices[firstIndex]
	})

	rule := placement_rule.NewPlacementRuleBuilder().BuildNearCrossroadsRule(2)
	var connections []entities.Connection
	for index := range min(count, maxCount) {
		labelIndex := indices[index]
		fromLabel := orderedLabels[labelIndex]
		toLabel := orderedLabels[destinations[labelIndex]]
		fromName := this.zoneLabelProvider.CreateZoneName(fromLabel, playerLabels)
		toName := this.zoneLabelProvider.CreateZoneName(toLabel, playerLabels)
		connections = append(connections, variant_content.NewConnectionBuilder().
			WithName(constants.GetPortalConnectionNameFor(fromLabel, toLabel)).
			WithFrom(fromName).
			WithTo(toName).
			WithConnectionTypePortal().
			WithPortalPlacementRulesFrom(rule).
			WithPortalPlacementRulesTo(rule).
			WithRoad(true).
			WithGuardValue(this.GetBorderGuardValue(fromLabel, toLabel, playerLabels, neutralZones, tuning)).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			Build())
	}
	return connections
}

func (this *TopologyConnectionService) CreateMissingPlayerConnections(
	playerLabels []string,
	zones []entities.Zone,
	connections []entities.Connection,
	tuning models.GenerationTuning) []entities.Connection {
	if len(playerLabels) < 2 {
		return nil
	}

	connectionNames := map[string]bool{}
	for _, connection := range connections {
		if connection.Name != "" {
			connectionNames[connection.Name] = true
		}
	}
	var additionalConnections []entities.Connection
	for _, label := range playerLabels {
		zoneName := constants.GetPlayerZoneNameFor(label)
		zone, ok := linq.FromSlice(zones).First(func(candidate entities.Zone) bool {
			return candidate.Name == zoneName
		})
		if !ok || spawnZoneHasConnection(zone, connectionNames) {
			continue
		}

		partner := linq.FromSlice(playerLabels).FirstOrDefault(func(candidate string) bool {
			return candidate != label
		})
		if partner == "" {
			continue
		}

		labelFrom, labelTo := label, partner
		if label > partner {
			labelFrom, labelTo = partner, label
		}
		fallbackName := constants.GetFallbackConnectionNameFor(labelFrom, labelTo)
		if connectionNames[fallbackName] {
			continue
		}

		additionalConnections = append(additionalConnections, variant_content.NewConnectionBuilder().
			WithName(fallbackName).
			WithFrom(zoneName).
			WithTo(constants.GetPlayerZoneNameFor(partner)).
			WithConnectionTypeDirect().
			WithGuardZone(zoneName).
			WithSimTurnSquad().
			WithGuardValue(this.GetBorderGuardValue(label, partner, playerLabels, nil, tuning)).
			WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
			WithGuardMatchGroup("fallback_guard_"+fallbackName).
			Build())
		connectionNames[fallbackName] = true
		appendSpawnFallbackRoads(zones, label, partner, fallbackName)
	}
	return additionalConnections
}

func (this *TopologyConnectionService) CreateMissingConnections(
	playerLabels, allLabels []string,
	positions models.Positions,
	zones []entities.Zone,
	connections []entities.Connection,
	tuning models.GenerationTuning,
	neutralZones neutral_zone.Plans,
) []entities.Connection {
	if len(allLabels) < 2 {
		return nil
	}

	adjacency := this.buildZoneAdjacency(playerLabels, allLabels, connections)
	connectionNames := map[string]bool{}
	for connection := range linq.FromSlice(connections).
		Where(func(candidate entities.Connection) bool { return candidate.Name != "" }).Iterate {
		connectionNames[connection.Name] = true
	}

	var additionalConnections []entities.Connection
	for {
		nodes := make([]int, len(allLabels))
		for index := range nodes {
			nodes[index] = index
		}
		bestIndexes, ok := geometry_helpers.FindClosestAcrossComponents(
			positions,
			adjacency.ConnectedComponents(nodes))
		if !ok {
			break
		}

		labelA, labelB := allLabels[bestIndexes.X], allLabels[bestIndexes.Y]
		if labelA > labelB {
			labelA, labelB = labelB, labelA
		}
		bridgeName := constants.GetBridgeConnectionNameFor(labelA, labelB)
		if connectionNames[bridgeName] {
			adjacency.Link(bestIndexes.X, bestIndexes.Y)
			continue
		}

		zoneFrom := this.zoneLabelProvider.CreateZoneName(allLabels[bestIndexes.X], playerLabels)
		zoneTo := this.zoneLabelProvider.CreateZoneName(allLabels[bestIndexes.Y], playerLabels)
		additionalConnections = append(additionalConnections, this.createBridgeConnection(
			bridgeName, zoneFrom, zoneTo, labelA, labelB, playerLabels, neutralZones, tuning))
		connectionNames[bridgeName] = true
		appendBridgeRoads(zones, zoneFrom, zoneTo, bridgeName)
		adjacency.Link(bestIndexes.X, bestIndexes.Y)
	}

	return additionalConnections
}

func (this *TopologyConnectionService) GetBorderGuardValue(
	labelA, labelB string,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
) int {
	higherQuality := max(
		labelQuality(labelA, playerLabels, neutralZones),
		labelQuality(labelB, playerLabels, neutralZones))
	return tuning.ScaleByBorderGuardStrength(higherQuality.GetGuardValue())
}

// labelQuality ranks a label for guarding purposes. QualityUnknown is -1, so a
// player label always loses the max against a real tier while still supplying
// the player-border guard value when both endpoints are players.
func labelQuality(
	label string,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
) neutral_zone.Quality {
	if zone_helpers.IsZoneNameHub(label) {
		return neutral_zone.QualityHighest
	}
	if slices.Contains(playerLabels, label) {
		return neutral_zone.QualityUnknown
	}
	return neutralZones.GetQuality(label)
}

func (this *TopologyConnectionService) createBridgeConnection(
	bridgeName, zoneFrom, zoneTo, labelA, labelB string,
	playerLabels []string,
	neutralZones neutral_zone.Plans,
	tuning models.GenerationTuning,
) entities.Connection {
	return variant_content.NewConnectionBuilder().
		WithName(bridgeName).
		WithFrom(zoneFrom).
		WithTo(zoneTo).
		WithConnectionTypeDirect().
		WithGuardZone(zoneFrom).
		WithSimTurnSquad().
		WithGuardValue(this.GetBorderGuardValue(labelA, labelB, playerLabels, neutralZones, tuning)).
		WithGuardWeeklyIncrement(common_connections.GetGuardWeeklyIncrements().Standard).
		WithGuardMatchGroup(fmt.Sprintf("bridge_guard_%s-%s", labelA, labelB)).
		Build()
}

func (this *TopologyConnectionService) buildZoneAdjacency(
	playerLabels, allLabels []string,
	connections []entities.Connection,
) data.Adjacency[int] {
	nodes := make([]int, len(allLabels))
	for index := range nodes {
		nodes[index] = index
	}
	adjacency := data.NewAdjacency(nodes)
	zoneNameToIndex := map[string]int{}
	for index, label := range allLabels {
		zoneNameToIndex[this.zoneLabelProvider.CreateZoneName(label, playerLabels)] = index
	}
	for connection := range linq.FromSlice(connections).
		Where(func(candidate entities.Connection) bool {
			connectionTypes := registry.GetConnectionTypeValues()
			return candidate.ConnectionType == connectionTypes.Direct ||
				candidate.ConnectionType == connectionTypes.Portal
		}).
		Where(func(candidate entities.Connection) bool {
			_, hasFrom := zoneNameToIndex[candidate.From]
			_, hasTo := zoneNameToIndex[candidate.To]
			return hasFrom && hasTo
		}).Iterate {
		adjacency.Link(zoneNameToIndex[connection.From], zoneNameToIndex[connection.To])
	}
	return adjacency
}

func spawnZoneHasConnection(zone entities.Zone, connectionNames map[string]bool) bool {
	connectionType := registry.GetRoadConnectionTypeValues().Connection
	for _, road := range zone.Roads {
		if road.To.Type == connectionType && len(road.To.Args) > 0 && connectionNames[road.To.Args[0]] {
			return true
		}
	}

	return false
}

func appendSpawnFallbackRoads(zones []entities.Zone, label, partner, fallbackName string) {
	for _, playerLabel := range []string{label, partner} {
		spawnZoneName := constants.GetPlayerZoneNameFor(playerLabel)
		zoneIndex := slices.IndexFunc(zones, func(candidate entities.Zone) bool {
			return candidate.Name == spawnZoneName
		})
		if zoneIndex >= 0 {
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads, variant_content.NewRoadBuilder().
				WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).
				WithTo(variant_content.NewRefBuilder().BuildConnectionType(fallbackName)).
				Build())
		}
	}
}

func findExistingConnectionName(zone entities.Zone) string {
	connectionType := registry.GetRoadConnectionTypeValues().Connection
	for _, road := range zone.Roads {
		if road.From.Type == connectionType && len(road.From.Args) > 0 {
			return road.From.Args[0]
		}

		if road.To.Type == connectionType && len(road.To.Args) > 0 {
			return road.To.Args[0]
		}
	}

	return ""
}

func appendBridgeRoads(zones []entities.Zone, zoneFrom, zoneTo, bridgeName string) {
	for _, zoneName := range []string{zoneFrom, zoneTo} {
		zoneIndex := slices.IndexFunc(zones, func(candidate entities.Zone) bool {
			return candidate.Name == zoneName
		})
		if zoneIndex < 0 {
			continue
		}
		roadBuilder := variant_content.NewRoadBuilder().WithTo(
			variant_content.NewRefBuilder().BuildConnectionType(bridgeName))
		switch {
		case len(zones[zoneIndex].MainObjects) > 0:
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads,
				roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildMainObjectType("0")).Build())
		case len(zones[zoneIndex].Roads) > 0:
			connectionName := findExistingConnectionName(zones[zoneIndex])
			if connectionName == "" {
				connectionName = bridgeName
			}
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads,
				roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(connectionName)).Build())
		default:
			zones[zoneIndex].Roads = append(zones[zoneIndex].Roads,
				roadBuilder.WithFrom(variant_content.NewRefBuilder().BuildConnectionType(bridgeName)).Build())
		}
	}
}

func buildNonAdjacentDerangement(count int) []int {
	for range 100 {
		if destinations, ok := tryRandomDerangement(count); ok {
			return destinations
		}
	}
	return buildShiftDerangement(count)
}

func tryRandomDerangement(count int) ([]int, bool) {
	candidates := make([]int, count)
	for index := range candidates {
		candidates[index] = index
	}
	rand.Shuffle(len(candidates), func(firstIndex, secondIndex int) {
		candidates[firstIndex], candidates[secondIndex] = candidates[secondIndex], candidates[firstIndex]
	})

	destinations := make([]int, count)
	used := make([]bool, count)
	for index := range count {
		foundIndex := pickDerangementTarget(candidates, used, index, count)
		if foundIndex < 0 {
			return nil, false
		}
		destinations[index] = candidates[foundIndex]
		used[candidates[foundIndex]] = true
	}
	return destinations, true
}

func pickDerangementTarget(candidates []int, used []bool, sourceIndex, count int) int {
	for index := range candidates {
		if used[candidates[index]] {
			continue
		}
		candidate := candidates[index]
		if candidate != sourceIndex &&
			candidate != (sourceIndex+1)%count &&
			candidate != (sourceIndex-1+count)%count {
			return index
		}
	}
	for index := range candidates {
		if !used[candidates[index]] && candidates[index] != sourceIndex {
			return index
		}
	}
	return -1
}

func buildShiftDerangement(count int) []int {
	destinations := make([]int, count)
	shift := max(1, count/2)
	for index := range count {
		destinations[index] = (index + shift) % count
	}
	return destinations
}
