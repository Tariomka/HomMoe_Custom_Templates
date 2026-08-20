package preview_service

import (
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type PreviewLayoutService struct {
	layout         *preview.Layout
	zoneClassifier zone_interfaces.IZoneClassifier
}

func NewPreviewLayoutService() IPreviewLayoutService {
	return &PreviewLayoutService{zoneClassifier: zone_services.NewZoneClassifier()}
}

// BuildPreviewLayout computes zone positions, radius and connections for a
// preview canvas of the given side length. The layout strategy is picked to
// match the in-game generator: Circles uses concentric rings keyed off the
// GeneratorRing stamps; Square, Geometric, Cross and Fractal are placed
// verbatim from their GeneratorPosition stamps (centered and scaled to fit) so
// the exact geometric figure is preserved; Random scatters zones using the
// GeneratorPosition stamps with hard-floor and edge-clearance correction
// passes; all other topologies fall back to the classic ring / hub-and-spoke
// renderer.
func (this *PreviewLayoutService) BuildPreviewLayout(
	template *entities.RmgTemplate,
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

	// Apply the optional ZeroAngleZone rotation so the first ring slot lines
	// up with the template author's chosen anchor, then lay out every zone with
	// the topology-specific renderer. Tournament templates are not special-
	// cased here: both player clusters are laid out together at full canvas
	// size (the generator seeds the two halves with mirrored positions and, for
	// hub topologies, layoutMultiHub fans the clusters out), so the preview and
	// the zone editor share one consistent, fully reversible coordinate system.
	zones := orderZonesByZeroAngle(variant.Zones, variant.Orientation.ZeroAngleZone)
	this.dispatchClusterLayout(zones, variant.Connections, topology, side)

	this.buildPreviewZones(variant.Zones)
	this.layout.Connections = this.buildPreviewConnections(variant.Connections, this.layout.Positions)

	layout := *this.layout
	this.layout = nil
	return layout
}

// dispatchClusterLayout writes positions for the given zones into the layout,
// picking the topology-specific renderer. Each path sets layout.Positions and
// layout.ZoneRadius.
func (this *PreviewLayoutService) dispatchClusterLayout(
	zones []entities.Zone,
	connections []entities.Connection,
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

// layoutManualPositions places zones exactly where the manual zone editor put
// them: canvas = normalized position × side. The mapping must stay trivially
// invertible (p = pos / side) so dragging in the editor is exact. The zone
// radius shrinks just enough to keep the closest pair of zones from
// overlapping.
func (this *PreviewLayoutService) layoutManualPositions(zones []entities.Zone, side float64) {
	metrics := newCanvasMetrics(side)

	var positions models.Positions
	for _, zone := range zones {
		p := *zone.ManualPosition // Is this required to be copied? can't it be used directly safely?
		positions.Add(data.NewVec2(p[0], p[1]).MultiplyScalar(side))
	}
	radius := radiusFromClosestPair(positions, metrics.zoneRadiusMax, metrics.minGap)
	this.commitPositions(zones, positions, radius)
}

// buildPreviewZones turns every positioned zone into its drawable preview
// form: letter, tier, hub/player flags and castle markers.
//
// A zone is only drawn as a hub when the template actually contains a hub
// zone (named "Hub" or "Hub-*"). Connectivity-based guesses are not used
// here: in topologies like Random or Circles an ordinary neutral can happen
// to touch every spawn without being a hub, which previously made the hub
// marker appear (and flicker) on non-hub zones.
func (this *PreviewLayoutService) buildPreviewZones(zones []entities.Zone) {
	for _, zone := range zones {
		pos, ok := this.layout.Positions[zone.Name]
		if !ok {
			continue
		}

		previewZone := preview.Zone{
			Name:    zone.Name,
			Label:   helpers.GetZoneLabel(zone.Name),
			Center:  pos,
			Quality: this.zoneClassifier.GetQuality(zone),
			Type:    zone_helpers.GetZoneTypeFromName(zone.Name),
		}
		applyMainObjects(zone, &previewZone)
		this.layout.Zones = append(this.layout.Zones, previewZone)
	}
}

// applyMainObjects folds the zone's Spawn/City/GladiatorArena main objects into
// the preview zone's castle count, player-owner number and arena marker.
func applyMainObjects(zone entities.Zone, previewZone *preview.Zone) {
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

// buildPreviewConnections turns the variant's connections into drawable preview
// edges - only those whose endpoints survived the layout. Connections sharing
// the same unordered endpoint pair are grouped and each is given a
// perpendicular bulge so they do not collapse onto a single overlapping line;
// a lone edge keeps its control point on the midpoint and therefore renders
// straight.
func (this *PreviewLayoutService) buildPreviewConnections(
	connections []entities.Connection,
	positions map[string]data.Vec2[float64]) []preview.Connection {
	type pairKey struct{ start, end string }
	sortedKey := func(connection entities.Connection) pairKey {
		if connection.From > connection.To {
			return pairKey{connection.To, connection.From}
		}

		return pairKey{connection.From, connection.To}
	}

	visible := func(connection entities.Connection) bool {
		_, okFrom := positions[connection.From]
		_, okTo := positions[connection.To]
		return okFrom && okTo
	}

	counts := make(map[pairKey]int)
	for _, connection := range connections {
		if visible(connection) {
			counts[sortedKey(connection)]++
		}
	}

	result := make([]preview.Connection, 0, len(connections))
	indexInPair := make(map[pairKey]int)

	const spacingBetweenEdges = 21.0
	for _, connection := range connections {
		if !visible(connection) {
			continue
		}

		key := sortedKey(connection)
		index := indexInPair[key]
		indexInPair[key]++

		startPoint := positions[key.start]
		endPoint := positions[key.end]
		delta := endPoint.Subtract(startPoint)
		distance := math.Max(math.Hypot(delta.X, delta.Y), 1)
		spread := (float64(index) - float64(counts[key]-1)/2.0) * spacingBetweenEdges
		// Ctrl offset is 2× the desired bulge: a quadratic Bézier's midpoint
		// sits halfway between the chord midpoint and the control point.
		ctrl := startPoint.Add(endPoint).MultiplyScalar(0.5).
			// ( x, y ) → ( y, -x ) rotates it 90°
			Add(data.NewVec2(delta.Y, -delta.X).MultiplyScalar(2.0 * spread / distance))
		result = append(
			result,
			preview.Connection{
				Start: startPoint,
				End:   endPoint,
				Ctrl:  ctrl,
				Type:  getPreviewConnectionType(connection),
			})
	}
	return result
}

// getPreviewConnectionType maps a template connection onto the drawable preview
// type. A connection also counts as a portal when it merely carries portal
// placement rules, because the in-game generator treats it as one.
func getPreviewConnectionType(connection entities.Connection) preview.ConnectionType {
	connectionTypes := registry.GetConnectionTypeValues()
	if connection.ConnectionType == connectionTypes.Portal ||
		len(connection.PortalPlacementRulesFrom) > 0 ||
		len(connection.PortalPlacementRulesTo) > 0 {
		return preview.ConnectionTypePortal
	}

	switch connection.ConnectionType {
	case connectionTypes.GladiatorArena:
		return preview.ConnectionTypeGladiatorArena
	case connectionTypes.Proximity:
		return preview.ConnectionTypeProximity
	default:
		return preview.ConnectionTypeDirect
	}
}
