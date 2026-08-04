package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
)

// NewTopologyProvider builds a provider over the full topology lookup, the way
// internal/composition wires it for production.
func NewTopologyProvider() *providers.TopologyProvider {
	return providers.NewTopologyProvider(NewTopologyServiceLookup(NewZoneFactories()))
}
