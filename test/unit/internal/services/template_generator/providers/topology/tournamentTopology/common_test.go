package tournamentTopology_test

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// connectedComponentCount returns the number of connected components in the
// variant's zone graph, where every connection links its From and To zones.
func connectedComponentCount(variant template_model.Variant) int {
	parent := map[string]string{}
	var find func(name string) string
	find = func(name string) string {
		if parent[name] != name {
			parent[name] = find(parent[name])
		}
		return parent[name]
	}
	union := func(first, second string) {
		rootFirst, rootSecond := find(first), find(second)
		if rootFirst != rootSecond {
			parent[rootFirst] = rootSecond
		}
	}

	for _, zone := range variant.Zones {
		parent[zone.Name] = zone.Name
	}
	for _, connection := range variant.Connections {
		if _, ok := parent[connection.From]; !ok {
			continue
		}
		if _, ok := parent[connection.To]; !ok {
			continue
		}
		union(connection.From, connection.To)
	}

	roots := map[string]bool{}
	for _, zone := range variant.Zones {
		roots[find(zone.Name)] = true
	}
	return len(roots)
}

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

// countConnectionsWithPrefix counts the connections whose name starts with the
// given prefix.
func countConnectionsWithPrefix(variant template_model.Variant, prefix string) int {
	count := 0
	for _, connection := range variant.Connections {
		if strings.HasPrefix(connection.Name, prefix) {
			count++
		}
	}
	return count
}
