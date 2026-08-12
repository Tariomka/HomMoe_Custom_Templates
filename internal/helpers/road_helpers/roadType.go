package road_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

func IsRoadTypeConnection(road entities.Road) bool {
	connectionType := registry.GetRoadConnectionTypeValues().Connection
	return road.From.Type == connectionType || road.To.Type == connectionType
}

// IsRoadTypeCastle reports whether a road connects two of the zone's own main
// objects (the stone castle<->castle roads). These must be regenerated whenever
// the zone's main-object count changes, otherwise added castles end up with no
// road and removed castles leave dangling roads.
func IsRoadTypeCastle(road entities.Road) bool {
	connectionType := registry.GetRoadConnectionTypeValues().MainObject
	return road.From.Type == connectionType && road.To.Type == connectionType
}
