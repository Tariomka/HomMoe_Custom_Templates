package geometricHubTopology_test

import (
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
)

// buildGeoHubVariant runs the Geometric Hub topology service with the given
// players and plans using default generator options (no random portals).
func buildGeoHubVariant(playerLabels []string, plans neutralZone.Plans) entities.Variant {
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyGeometricHub
	configuration.PlayerCount = len(playerLabels)
	tuning := models.NewGenerationTuning(configuration, len(playerLabels)+len(plans)+1)
	return topology.NewGeometricHubTopologyService().
		CreateTopologyVariant(*configuration, playerLabels, plans, tuning, false)
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
func mixedPlans(mediumLabels, lowLabels, highLabels []string) neutralZone.Plans {
	plans := neutralZone.Plans{}
	for _, label := range mediumLabels {
		plans.AddPlan(label, neutralZone.QualityMedium, 1)
	}
	for _, label := range lowLabels {
		plans.AddPlan(label, neutralZone.QualityLow, 0)
	}
	for _, label := range highLabels {
		plans.AddPlan(label, neutralZone.QualityHigh, 1)
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
	return [2]float64{}
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
