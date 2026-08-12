package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

// NewZoneFactories builds the collaborators that every topology service
// constructor takes. The results match those parameter lists exactly, so it can
// be spread directly into a constructor call:
//
//	topology.NewRingTopologyService(test_helpers.NewZoneFactories())
func NewZoneFactories() (
	zone_interfaces.IZoneFactory,
	zone_interfaces.IRoadFactory,
	zone_interfaces.IZoneLabelProvider,
	base.ITopologyConnectionService) {
	roadFactory := zones.NewRoadFactory()
	zoneLabelProvider := zones.NewZoneLabelProvider()
	return zones.NewZoneFactory(zones.NewCastleFactory(), roadFactory),
		roadFactory,
		zoneLabelProvider,
		base.NewTopologyConnectionService(zoneLabelProvider)
}
