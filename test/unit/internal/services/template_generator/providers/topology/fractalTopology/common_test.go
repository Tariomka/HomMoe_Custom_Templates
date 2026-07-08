package fractalTopology_test

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// zoneNameSet returns the set of zone names present in the variant.
func zoneNameSet(variant entities.Variant) map[string]bool {
	names := make(map[string]bool, len(variant.Zones))
	for _, zone := range variant.Zones {
		names[zone.Name] = true
	}
	return names
}

// danglingConnectionNames returns the names of connections whose endpoints do
// not both exist as zones in the variant.
func danglingConnectionNames(variant entities.Variant) []string {
	names := zoneNameSet(variant)
	var dangling []string
	for _, connection := range variant.Connections {
		if !names[connection.From] || !names[connection.To] {
			dangling = append(dangling, connection.Name)
		}
	}
	return dangling
}

// zonesWithoutValidPosition returns the names of zones whose generator position
// is missing or falls outside the normalized [0,1]x[0,1] layout square.
func zonesWithoutValidPosition(variant entities.Variant) []string {
	var invalid []string
	for _, zone := range variant.Zones {
		if zone.GeneratorPosition == nil {
			invalid = append(invalid, zone.Name)
			continue
		}
		positionX, positionY := zone.GeneratorPosition[0], zone.GeneratorPosition[1]
		if positionX < 0 || positionX > 1 || positionY < 0 || positionY > 1 {
			invalid = append(invalid, zone.Name)
		}
	}
	return invalid
}

// countPortalConnections counts the connections of type Portal.
func countPortalConnections(variant entities.Variant) int {
	count := 0
	for _, connection := range variant.Connections {
		if connection.ConnectionType == "Portal" {
			count++
		}
	}
	return count
}

// randomSpawnToSpawnNames returns the names of triangulation-driven ("Rnd-")
// connections that join two spawn zones.
func randomSpawnToSpawnNames(variant entities.Variant) []string {
	var names []string
	for _, connection := range variant.Connections {
		if !strings.HasPrefix(connection.Name, "Rnd-") {
			continue
		}
		if strings.HasPrefix(connection.From, "Spawn-") && strings.HasPrefix(connection.To, "Spawn-") {
			names = append(names, connection.Name)
		}
	}
	return names
}
