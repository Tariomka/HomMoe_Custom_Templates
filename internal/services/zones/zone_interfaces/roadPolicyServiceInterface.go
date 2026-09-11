package zone_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type IRoadPolicyService interface {
	// Reconcile stamps the road setting onto the connection graph and brings
	// every zone's roads in line with the final graph and content.
	Reconcile(request models.RoadReconciliationRequest)

	RebuildZoneConnectionRoads(zones []template_model.Zone, connections []template_model.Connection)

	// RebuildCastleRoads reconciles only the zone's castle<->castle roads with
	// its current main objects, regardless of the road setting.
	RebuildCastleRoads(zone *template_model.Zone)

	// StampConnectionRoad applies the road setting to a single connection,
	// leaving an explicit Portal's own flag untouched.
	StampConnectionRoad(connection *template_model.Connection, generateRoads bool)
}
