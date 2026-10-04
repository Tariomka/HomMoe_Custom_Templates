package topology

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/topology_interfaces"
)

// Every topology service must satisfy topology_interfaces.ITopologyService.
var (
	_ topology_interfaces.ITopologyService = (*TournamentTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*GeometricHubTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*RandomTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*CirclesTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*SquareTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*GeometricTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*CrossTopologyService)(nil)
	_ topology_interfaces.ITopologyService = (*FractalTopologyService)(nil)
)
