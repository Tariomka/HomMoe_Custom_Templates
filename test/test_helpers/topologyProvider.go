package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
)

// NewTopologyProvider builds a provider over the full topology lookup, the way
// internal/composition wires it for production.
func NewTopologyProvider() provider_interfaces.ITopologyProvider {
	return providers.NewTopologyProvider(NewTopologyServiceLookup(NewZoneFactories()))
}
