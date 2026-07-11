package hubTopology_test

import (
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

// connectionNames returns the names of all connections in the variant.
func connectionNames(variant entities.Variant) []string {
	names := make([]string, 0, len(variant.Connections))
	for _, connection := range variant.Connections {
		names = append(names, connection.Name)
	}
	return names
}
