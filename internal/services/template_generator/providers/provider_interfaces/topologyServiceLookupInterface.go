package provider_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type ITopologyServiceLookup interface {
	Tournament() TopologyVariantCreator
	Resolve(mapTopology config.MapTopology) TopologyVariantCreator
}
