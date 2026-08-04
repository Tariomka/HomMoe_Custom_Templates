package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// NewZoneFactories builds the zone and road factories that every topology
// service constructor takes. The results match those parameter lists exactly,
// so it can be spread directly into a constructor call:
//
//	topology.NewRingTopologyService(test_helpers.NewZoneFactories())
func NewZoneFactories() (*zones.ZoneFactory, *zones.RoadFactory) {
	roadFactory := zones.NewRoadFactory()
	return zones.NewZoneFactory(zones.NewCastleFactory(), roadFactory), roadFactory
}
