package geometricHubTopology_test

import (
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

// buildGeoHubVariant runs the Geometric Hub topology service with the given
// players and plans using default generator options (no random portals).
func buildGeoHubVariant(playerLabels []string, plans neutral_zone.Plans) entities.Variant {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	configuration.PlayerCount = len(playerLabels)
	tuning := test_helpers.NewGenerationTuning(configuration, len(playerLabels)+len(plans)+1)
	return topology.NewGeometricHubTopologyService(test_helpers.NewZoneFactories()).
		CreateTopologyVariant(*configuration, playerLabels, plans, tuning, "")
}

// neighborsOf collects the names of every zone connected to zoneName.
func neighborsOf(variant entities.Variant, zoneName string) []string {
	var neighbors []string
	for _, connection := range variant.Connections {
		if connection.From == zoneName {
			neighbors = append(neighbors, connection.To)
		} else if connection.To == zoneName {
			neighbors = append(neighbors, connection.From)
		}
	}
	return neighbors
}

// hubConnections returns every connection with the Hub as an endpoint.
func hubConnections(variant entities.Variant) []entities.Connection {
	var connections []entities.Connection
	for _, connection := range variant.Connections {
		if connection.From == "Hub" || connection.To == "Hub" {
			connections = append(connections, connection)
		}
	}
	return connections
}

// hubPortalTargets returns the names of the zones connected to the Hub.
func hubPortalTargets(variant entities.Variant) []string {
	return neighborsOf(variant, "Hub")
}

// spawnNeighborsOf filters the neighbors of zoneName down to spawn zones.
func spawnNeighborsOf(variant entities.Variant, zoneName string) []string {
	var spawns []string
	for _, neighbor := range neighborsOf(variant, zoneName) {
		if strings.HasPrefix(neighbor, "Spawn-") {
			spawns = append(spawns, neighbor)
		}
	}
	return spawns
}

// mixedPlans builds the standard test plan set: mediums first, then lows,
// then highs, so tier-to-slot assignment is observable by label.
func mixedPlans(mediumLabels, lowLabels, highLabels []string) neutral_zone.Plans {
	plans := neutral_zone.Plans{}
	for _, label := range mediumLabels {
		plans.AddPlan(label, neutral_zone.QualityMedium, 1)
	}
	for _, label := range lowLabels {
		plans.AddPlan(label, neutral_zone.QualityLow, 0)
	}
	for _, label := range highLabels {
		plans.AddPlan(label, neutral_zone.QualityHigh, 1)
	}
	return plans
}

// positionOf returns the generator position of the named zone.
func positionOf(variant entities.Variant, zoneName string) [2]float64 {
	for _, zone := range variant.Zones {
		if zone.Name == zoneName && zone.GeneratorPosition != nil {
			return *zone.GeneratorPosition
		}
	}

	panic("zone not found: " + zoneName)
}

// distanceBetween returns the Euclidean distance between two zones' positions.
func distanceBetween(variant entities.Variant, firstZone, secondZone string) float64 {
	first := positionOf(variant, firstZone)
	second := positionOf(variant, secondZone)
	return math.Hypot(first[0]-second[0], first[1]-second[1])
}

// spreadOf returns the difference between the largest and smallest value.
func spreadOf(values []float64) float64 {
	lowest, highest := values[0], values[0]
	for _, value := range values[1:] {
		lowest = math.Min(lowest, value)
		highest = math.Max(highest, value)
	}
	return highest - lowest
}

// interiorAngleAt returns the interior angle (degrees) of the hexagon ring at
// zone `at`, formed by its ring neighbors `previous` and `next`.
func interiorAngleAt(variant entities.Variant, previous, at, next string) float64 {
	previousPosition := positionOf(variant, previous)
	atPosition := positionOf(variant, at)
	nextPosition := positionOf(variant, next)
	toPreviousX := previousPosition[0] - atPosition[0]
	toPreviousY := previousPosition[1] - atPosition[1]
	toNextX := nextPosition[0] - atPosition[0]
	toNextY := nextPosition[1] - atPosition[1]
	dot := toPreviousX*toNextX + toPreviousY*toNextY
	magnitudes := math.Hypot(toPreviousX, toPreviousY) * math.Hypot(toNextX, toNextY)
	return math.Acos(dot/magnitudes) * 180 / math.Pi
}

// perimeterFreeAngles returns the interior angles (degrees) at the five
// non-hub vertices of a hexagon perimeter listed hub-first.
func perimeterFreeAngles(variant entities.Variant, perimeter [6]string) []float64 {
	angles := make([]float64, 0, 5)
	for index := 1; index < len(perimeter); index++ {
		previous := perimeter[index-1]
		next := perimeter[(index+1)%len(perimeter)]
		angles = append(angles, interiorAngleAt(variant, previous, perimeter[index], next))
	}
	return angles
}
