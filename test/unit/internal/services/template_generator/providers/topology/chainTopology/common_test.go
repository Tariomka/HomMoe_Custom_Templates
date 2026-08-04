package chainTopology_test

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
)

func newChainTopologyService() *topology.ChainTopologyService {
	return topology.NewChainTopologyService(test_helpers.NewZoneFactories())
}

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

// countConnectionsWithPrefix counts the connections whose name starts with the
// given prefix.
func countConnectionsWithPrefix(variant entities.Variant, prefix string) int {
	count := 0
	for _, connection := range variant.Connections {
		if strings.HasPrefix(connection.Name, prefix) {
			count++
		}
	}
	return count
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
