package geometricTopology_test

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// zoneNameSet returns the set of zone names present in the variant.
func zoneNameSet(variant template_model.Variant) map[string]bool {
	names := make(map[string]bool, len(variant.Zones))
	for _, zone := range variant.Zones {
		names[zone.Name] = true
	}
	return names
}

// danglingConnectionNames returns the names of connections whose endpoints do
// not both exist as zones in the variant.
func danglingConnectionNames(variant template_model.Variant) []string {
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
func countPortalConnections(variant template_model.Variant) int {
	count := 0
	for _, connection := range variant.Connections {
		if connection.ConnectionType == "Portal" {
			count++
		}
	}
	return count
}

// directSpawnToSpawnNames returns the names of Direct connections with the
// given name prefix that join two spawn zones.
func directSpawnToSpawnNames(variant template_model.Variant, prefix string) []string {
	var names []string
	for _, connection := range variant.Connections {
		if connection.ConnectionType != "Direct" || !strings.HasPrefix(connection.Name, prefix) {
			continue
		}
		if strings.HasPrefix(connection.From, "Spawn-") && strings.HasPrefix(connection.To, "Spawn-") {
			names = append(names, connection.Name)
		}
	}
	return names
}
