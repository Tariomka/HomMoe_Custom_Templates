package previewLayoutService_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// position returns a pointer to a normalized [0,1] coordinate pair.
func position(x, y float64) *[2]float64 {
	point := [2]float64{x, y}
	return &point
}

// namedZone builds a minimal zone carrying only a name.
func namedZone(name string) entities.Zone { return entities.Zone{Name: name} }

// positionedZone builds a zone with a generator position stamp.
func positionedZone(name string, x, y float64) entities.Zone {
	zone := namedZone(name)
	zone.GeneratorPosition = position(x, y)
	return zone
}

// manualZone builds a zone with a manual editor position stamp.
func manualZone(name string, x, y float64) entities.Zone {
	zone := namedZone(name)
	zone.ManualPosition = position(x, y)
	return zone
}

// directConnection builds a plain direct connection between two zones.
func directConnection(from, to string) entities.Connection {
	return entities.Connection{From: from, To: to, ConnectionType: "Direct"}
}

// templateWith wraps zones and connections into a single-variant template.
func templateWith(zones []entities.Zone, connections []entities.Connection) *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Variants: []entities.Variant{{Zones: zones, Connections: connections}},
	}
}
