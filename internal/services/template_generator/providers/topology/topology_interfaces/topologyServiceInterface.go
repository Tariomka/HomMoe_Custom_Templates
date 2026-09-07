package topology_interfaces

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ITopologyService is the contract every topology service satisfies: it turns a
// generator configuration and a set of zone labels into a finished variant.
//
// The twelve implementations are deliberately still constructed as concrete
// types, because internal/composition resolves them through wire and wire keys
// providers by their output type — twelve providers returning this one interface
// would be indistinguishable. The interface therefore serves the compile-time
// assertions in topologyServiceAssertions.go and any test that wants to treat
// the services uniformly.
type ITopologyService interface {
	CreateTopologyVariant(
		configuration config.GeneratorConfig,
		playerLabels []string,
		neutralZones neutral_zone.Plans,
		tuning models.GenerationTuning,
		holdCityNeutralLabel string) template_model.Variant
}
