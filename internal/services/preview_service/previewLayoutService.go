package preview_service

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type PreviewLayoutService struct {
	layout      *preview.Layout
	tierService zone_interfaces.IZoneTierService
}

func NewPreviewLayoutService(tierService zone_interfaces.IZoneTierService) IPreviewLayoutService {
	return &PreviewLayoutService{tierService: tierService}
}

// BuildPreviewLayout computes zone positions, radius and connections for a preview canvas of the given side length.
func (this *PreviewLayoutService) BuildPreviewLayout(
	template *template_model.Template,
	topology config.MapTopology,
	side float64) preview.Layout {
	this.layout = &preview.Layout{Positions: map[string]data.Vec2[float64]{}}
	if template == nil || len(template.Variants) == 0 {
		return *this.layout
	}

	variant := template.Variants[0]
	if len(variant.Zones) == 0 {
		return *this.layout
	}

	zones := orderZonesByZeroAngle(variant.Zones, variant.Orientation.ZeroAngleZone)
	this.dispatchClusterLayout(zones, variant.Connections, topology, side)

	this.buildPreviewZones(variant.Zones)
	this.layout.Connections = this.buildPreviewConnections(variant.Connections)

	layout := *this.layout
	this.layout = nil
	return layout
}

// dispatchClusterLayout writes positions for the given zones into the layout,
// picking the topology-specific renderer. Each path sets layout.Positions and
// layout.ZoneRadius.
func (this *PreviewLayoutService) dispatchClusterLayout(
	zones []template_model.Zone,
	connections []template_model.Connection,
	topology config.MapTopology,
	side float64) {
	capabilities := common_topologies.GetTopologyCapabilities(topology)
	switch {
	case allHaveManualPosition(zones):
		this.layoutManualPositions(zones, side)
	case capabilities.UsesGeneratorRing && allHaveRing(zones):
		this.layoutBalancedRings(zones, side)
	case capabilities.LayoutKind == models.TopologyLayoutFixedGeometry && allHavePosition(zones):
		this.layoutFixedPositions(zones, side, fixedGeometryEdgeInset(topology, zones))
	case capabilities.LayoutKind == models.TopologyLayoutScatter && allHavePosition(zones):
		this.layoutScatter(zones, connections, side)
	default:
		this.layoutRingOrHub(zones, connections, side)
	}
}

func (this *PreviewLayoutService) layoutManualPositions(zones []template_model.Zone, side float64) {
	metrics := newCanvasMetrics(side)

	var positions models.Positions
	for _, zone := range zones {
		positions.Add(zone.ManualPosition.MultiplyScalar(side))
	}
	radius := radiusFromClosestPair(positions, metrics.zoneRadiusMax, metrics.minGap)
	this.commitPositions(zones, positions, radius)
}

// buildPreviewZones turns every positioned zone into its drawable preview form.
func (this *PreviewLayoutService) buildPreviewZones(zones []template_model.Zone) {
	for _, zone := range zones {
		pos, ok := this.layout.Positions[zone.Name]
		if !ok {
			continue
		}

		previewZone := preview.Zone{
			Name:    zone.Name,
			Label:   helpers.GetZoneLabel(zone.Name),
			Center:  pos,
			Quality: this.tierService.ResolveQuality(zone),
			Type:    zone_helpers.GetZoneTypeFromName(zone.Name),
		}
		applyMainObjects(zone, &previewZone)
		this.layout.Zones = append(this.layout.Zones, previewZone)
	}
}

// applyMainObjects folds the zone's Spawn/City/GladiatorArena main objects into
// the preview zone's castle count, player-owner number and arena marker.
func applyMainObjects(zone template_model.Zone, previewZone *preview.Zone) {
	objectTypes := registry.GetMainObjectTypeValues()
	for _, mainObject := range zone.MainObjects {
		switch mainObject.Type {
		case objectTypes.Spawn:
			previewZone.Castles++
			if strings.HasPrefix(mainObject.Spawn, "Player") {
				for _, ch := range mainObject.Spawn[len("Player"):] {
					if ch >= '0' && ch <= '9' {
						previewZone.Owner = previewZone.Owner*10 + int(ch-'0')
					}
				}
			}
		case objectTypes.City:
			previewZone.Castles++
		case objectTypes.GladiatorArena:
			previewZone.Arena = true
		}
	}
}

func (this *PreviewLayoutService) buildPreviewConnections(
	connections []template_model.Connection) []preview.Connection {
	curves := preview.ConnectionCurveLayout{
		Positions:  this.layout.Positions,
		ZoneRadius: this.layout.ZoneRadius,
	}.Build(connections)
	result := make([]preview.Connection, 0, len(curves))
	for _, curve := range curves {
		connection := connections[curve.ConnectionIndex]
		result = append(
			result,
			preview.Connection{
				Start:          curve.Start,
				End:            curve.End,
				Ctrl:           curve.Control,
				Type:           getPreviewConnectionType(connection),
				HasRoad:        connection.HasRoad(),
				ExplicitPortal: connection.IsExplicitPortal(),
			})
	}
	return result
}

func getPreviewConnectionType(connection template_model.Connection) preview.ConnectionType {
	if connection.IsEffectivePortal() {
		return preview.ConnectionTypePortal
	}

	connectionTypes := registry.GetConnectionTypeValues()
	switch connection.ConnectionType {
	case connectionTypes.GladiatorArena:
		return preview.ConnectionTypeGladiatorArena
	case connectionTypes.Proximity:
		return preview.ConnectionTypeProximity
	default:
		return preview.ConnectionTypeDirect
	}
}
