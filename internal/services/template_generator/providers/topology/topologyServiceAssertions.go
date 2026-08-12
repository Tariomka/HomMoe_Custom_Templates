package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/topology_interfaces"
)

// Every topology service must satisfy topology_interfaces.ITopologyService.
// The constructors still return the concrete types so that wire can tell the
// twelve providers apart, so these assertions are what keeps the shared
// contract honest.
var (
	_ topology_interfaces.ITopologyService = (*TournamentTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*RingTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*HubTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*GeometricHubTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*ChainTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*SharedWebTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*RandomTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*CirclesTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*SquareTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*GeometricTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*CrossTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*FractalTopologyService)(nil)
)
