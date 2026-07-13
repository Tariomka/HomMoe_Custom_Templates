package preview_service

import (
	"image"
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

// layoutRingOrHub renders the structured topologies (Default, HubAndSpoke,
// Chain, SharedWeb). Multi-hub "Hub-*" templates fan their spokes out from
// each cluster center; otherwise zones land on a single outer ring with an
// optional center hub.
//
// Hub detection: only an explicitly named "Hub" / "Hub-*" zone is treated as
// a hub. The preview is a faithful representation of the template data, so
// connectivity is never used to guess an implicit hub.
func (this *PreviewLayoutService) layoutRingOrHub(
	zones []entities.Zone,
	conns []entities.Connection,
	side float64) {
	metrics := newCanvasMetrics(side)
	if this.placeTrivial(zones, metrics) {
		return
	}

	// Multi-hub tournament layout: clusters fan out around the canvas.
	var hubIndices []int
	for i, zone := range zones {
		if strings.HasPrefix(zone.Name, "Hub-") {
			hubIndices = append(hubIndices, i)
		}
	}
	if len(hubIndices) >= 2 {
		this.layoutMultiHub(zones, conns, hubIndices, metrics)
		return
	}

	hubIdx := -1
	for i, zone := range zones {
		if zone.Name == "Hub" {
			hubIdx = i
			break
		}
	}
	var outer []int
	for i := range zones {
		if i != hubIdx {
			outer = append(outer, i)
		}
	}

	zoneRadius := ringZoneRadius(len(outer), hubIdx >= 0, metrics)
	this.layout.ZoneRadius = int(math.Round(zoneRadius))

	ringRadius0 := side/2.0 - metrics.margin
	ringRadius := math.Max(metrics.hubRadiusMin+zoneRadius+metrics.minGap,
		math.Min(ringRadius0, side/2.0-zoneRadius-metrics.margin))

	if hubIdx >= 0 {
		this.layout.Positions[zones[hubIdx].Name] = metrics.center()
	}
	for i, zoneIndex := range outer {
		angle := -math.Pi/2.0 + float64(i)*2.0*math.Pi/float64(len(outer))
		x := metrics.cx + math.Cos(angle)*ringRadius
		y := metrics.cy + math.Sin(angle)*ringRadius
		this.layout.Positions[zones[zoneIndex].Name] = image.Pt(int(math.Round(x)), int(math.Round(y)))
	}
}

// ringZoneRadius sizes the ring zones so neighbouring chords keep a visible
// connection gap; hub layouts skip the self-clearance term because the hub
// occupies the center.
func ringZoneRadius(outerCount int, hasHub bool, metrics canvasMetrics) float64 {
	outerN := max(outerCount, 1)
	ringRadius0 := metrics.side/2.0 - metrics.margin
	sinA := 1.0
	if outerN > 1 {
		sinA = math.Sin(math.Pi / float64(outerN))
	}
	var zoneRadius float64
	if hasHub {
		zoneRadius = (2.0*ringRadius0*sinA - metrics.connectionGap) / 2.0
	} else {
		zoneRadius = (2.0*ringRadius0*sinA - metrics.connectionGap) / (2.0 * (1.0 + sinA))
	}
	return math.Min(metrics.zoneRadiusMax, math.Max(zoneRadius, 4.0))
}

// layoutMultiHub fans each "Hub-*" cluster out around the canvas, placing the
// hubs on an inner ring and their direct spokes around each hub. Zones that
// spoke off no hub (e.g. cross-cluster zones) collapse to the canvas center.
func (this *PreviewLayoutService) layoutMultiHub(
	zones []entities.Zone,
	conns []entities.Connection,
	hubIndices []int,
	metrics canvasMetrics) {
	hubSpokes := buildHubSpokes(zones, conns, hubIndices)
	numHubs := len(hubIndices)
	maxSpokes := 1
	for _, spokes := range hubSpokes {
		maxSpokes = max(maxSpokes, len(spokes))
	}

	canvasHalf := metrics.side/2.0 - metrics.margin
	sinB := 0.0
	if numHubs > 1 {
		sinB = math.Sin(math.Pi / float64(numHubs))
	}
	sinA := 1.0
	if maxSpokes > 1 {
		sinA = math.Sin(math.Pi / float64(maxSpokes))
	}
	hubRing := 0.0
	if numHubs > 1 {
		hubRing = (canvasHalf + metrics.minGap/2.0) / (1.0 + sinB)
	}
	radialLeft := canvasHalf - hubRing
	zoneRadius := math.Min(metrics.zoneRadiusMax, (radialLeft*sinA-metrics.minGap/2.0)/(1.0+sinA))
	zoneRadius = math.Max(1.0, zoneRadius)
	spokeRing := math.Max(radialLeft-zoneRadius, metrics.hubRadiusMin+metrics.minGap+zoneRadius)
	this.layout.ZoneRadius = int(math.Round(zoneRadius))

	for h, hubIndex := range hubIndices {
		hubAngle := -math.Pi/2.0 + float64(h)*2.0*math.Pi/float64(numHubs)
		hx, hy := metrics.cx, metrics.cy
		if numHubs > 1 {
			hx = metrics.cx + math.Cos(hubAngle)*hubRing
			hy = metrics.cy + math.Sin(hubAngle)*hubRing
		}
		this.layout.Positions[zones[hubIndex].Name] = image.Pt(int(math.Round(hx)), int(math.Round(hy)))

		spokes := hubSpokes[zones[hubIndex].Name]
		if len(spokes) == 0 {
			continue
		}
		spokeBase := hubAngle
		if numHubs == 1 {
			spokeBase = -math.Pi / 2.0
		}
		for i, spokeIndex := range spokes {
			angle := spokeBase + float64(i)*2.0*math.Pi/float64(len(spokes))
			x := hx + math.Cos(angle)*spokeRing
			y := hy + math.Sin(angle)*spokeRing
			this.layout.Positions[zones[spokeIndex].Name] = image.Pt(int(math.Round(x)), int(math.Round(y)))
		}
	}
	// Stragglers (e.g. cross-cluster zones) collapse to canvas center.
	for _, zone := range zones {
		if _, ok := this.layout.Positions[zone.Name]; !ok {
			this.layout.Positions[zone.Name] = metrics.center()
		}
	}
}

// buildHubSpokes collects each hub's directly connected zone indices
// (structural connections only, deduplicated, in connection order).
func buildHubSpokes(
	zones []entities.Zone,
	conns []entities.Connection,
	hubIndices []int) map[string][]int {
	zoneIdx := make(map[string]int, len(zones))
	for i, zone := range zones {
		zoneIdx[zone.Name] = i
	}
	hubSpokes := make(map[string][]int, len(hubIndices))
	for _, h := range hubIndices {
		hub := zones[h].Name
		seen := map[int]bool{}
		for _, conn := range conns {
			if isStructuralIgnored(conn.ConnectionType) {
				continue
			}
			other := ""
			switch {
			case conn.From == hub:
				other = conn.To
			case conn.To == hub:
				other = conn.From
			}
			if other == "" {
				continue
			}
			otherIndex, ok := zoneIdx[other]
			if !ok || seen[otherIndex] {
				continue
			}
			seen[otherIndex] = true
			hubSpokes[hub] = append(hubSpokes[hub], otherIndex)
		}
	}
	return hubSpokes
}
