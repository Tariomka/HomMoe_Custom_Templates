package previewLayoutService_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
)

// position returns a pointer to a normalized [0,1] coordinate pair.
func position(x, y float64) *[2]float64 {
	point := [2]float64{x, y}
	return &point
}

// namedZone builds a minimal zone carrying only a name.
func namedZone(name string) template_model.Zone { return template_model.Zone{Name: name} }

// positionedZone builds a zone with a generator position stamp.
func positionedZone(name string, x, y float64) template_model.Zone {
	zone := namedZone(name)
	zone.GeneratorPosition = position(x, y)
	return zone
}

// manualZone builds a zone with a manual editor position stamp.
func manualZone(name string, x, y float64) template_model.Zone {
	zone := namedZone(name)
	zone.ManualPosition = position(x, y)
	return zone
}

// ringedZone builds a zone with generator position and ring stamps.
func ringedZone(name string, ring int, x, y float64) template_model.Zone {
	zone := positionedZone(name, x, y)
	zone.GeneratorRing = new(ring)
	return zone
}

// directConnection builds a plain direct connection between two zones.
func directConnection(from, to string) template_model.Connection {
	return template_model.Connection{From: from, To: to, ConnectionType: "Direct"}
}

// lowTierZone builds a zone whose content pools infer as QualityLow but which
// records no tier of its own.
func lowTierZone(name string) template_model.Zone {
	zone := namedZone(name)
	zone.Layout = registry.GetLayoutValues().Sides
	zone.GuardedContentPool = []string{"pool_t2_a"}
	return zone
}

// templateWith wraps zones and connections into a single-variant template.
func templateWith(zones []template_model.Zone, connections []template_model.Connection) *template_model.Template {
	return &template_model.Template{
		Variants: []template_model.Variant{{Zones: zones, Connections: connections}},
	}
}
