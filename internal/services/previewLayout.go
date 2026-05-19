package services

import (
	"image"
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// PreviewZone is one zone laid out on the preview canvas.
type PreviewZone struct {
	Name      string
	Letter    string
	Center    image.Point
	IsPlayer  bool
	IsHub     bool
	Tier      int // 0 unknown, 1 bronze, 2 silver, 3 gold
	Owner     int
	HasCastle bool
	Castles   int
}

// PreviewConnection is a drawn link between two zones on the preview canvas.
type PreviewConnection struct {
	A, B   image.Point
	Portal bool
}

// PreviewLayout is the full geometry of a preview rendered into a square
// canvas of the requested side length.
type PreviewLayout struct {
	Positions   map[string]image.Point
	Zones       []PreviewZone
	Connections []PreviewConnection
	ZoneRadius  int
}

// BuildPreviewLayout computes zone positions, radius and connections for a
// preview canvas of the given side length. The chosen layout (ring, hub or
// hub-and-spoke) is inferred from zone names and connection topology.
func BuildPreviewLayout(template *models.RmgTemplate, topology models.MapTopology, side float64) PreviewLayout {
	layout := PreviewLayout{Positions: map[string]image.Point{}}
	if template == nil || len(template.Variants) == 0 {
		return layout
	}
	variant := template.Variants[0]
	zones := variant.Zones
	if len(zones) == 0 {
		return layout
	}

	// Build name lookup. Zone names are like "Spawn-A", "Neutral-C", "Hub".
	zoneNameByKey := make(map[string]string, len(zones)*2)
	for _, zone := range zones {
		zoneNameByKey[zone.Name] = zone.Name
	}
	resolveKey := func(key string) (string, bool) {
		name, ok := zoneNameByKey[key]
		return name, ok
	}

	// Layout: hub-and-spoke if multiple "Hub-*" zones; ring otherwise (with optional center "Hub").
	cx := side / 2
	cy := side / 2
	margin := 24.0

	// Detect multi-hub (tournament): zones literally named "Hub-*".
	var hubs []int
	for i, zone := range zones {
		if strings.HasPrefix(zone.Name, "Hub") {
			hubs = append(hubs, i)
		}
	}

	// Identify the implicit hub zone — the single non-player zone that
	// every player zone connects to. Works for HubAndSpoke and any other
	// topology where one neutral acts as a central hub. Falls back to a
	// plain ring when no such zone exists.
	implicitHubIdx := -1
	if len(hubs) < 2 {
		// Collect player zone names.
		playerNames := map[string]bool{}
		for _, zone := range zones {
			if strings.HasPrefix(zone.Name, "Spawn-") {
				playerNames[zone.Name] = true
			}
		}
		// neighbours[zoneName] = set of connected zone names.
		neighbours := make(map[string]map[string]bool, len(zones))
		for _, conn := range variant.Connections {
			a, ok1 := resolveKey(conn.From)
			b, ok2 := resolveKey(conn.To)
			if !ok1 || !ok2 {
				continue
			}
			if neighbours[a] == nil {
				neighbours[a] = map[string]bool{}
			}
			if neighbours[b] == nil {
				neighbours[b] = map[string]bool{}
			}
			neighbours[a][b] = true
			neighbours[b][a] = true
		}
		bestDeg := -1
		for i, zone := range zones {
			if strings.HasPrefix(zone.Name, "Spawn-") {
				continue
			}
			zoneNeighbours := neighbours[zone.Name]
			if len(zoneNeighbours) < 2 {
				continue
			}
			// Must connect to every player zone.
			connectsAllPlayers := len(playerNames) > 0
			for playerName := range playerNames {
				if !zoneNeighbours[playerName] {
					connectsAllPlayers = false
					break
				}
			}
			if !connectsAllPlayers {
				continue
			}
			if len(zoneNeighbours) > bestDeg {
				bestDeg = len(zoneNeighbours)
				implicitHubIdx = i
			}
		}
	}
	if len(hubs) >= 2 {
		// Build spoke lists.
		zoneIdx := map[string]int{}
		for i, zone := range zones {
			zoneIdx[zone.Name] = i
		}
		hubSpokes := make(map[string][]int, len(hubs))
		for _, hubIndex := range hubs {
			hub := zones[hubIndex].Name
			seen := map[int]bool{}
			for _, conn := range variant.Connections {
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
				otherName, ok := resolveKey(other)
				if !ok {
					continue
				}
				if otherIdx, ok := zoneIdx[otherName]; ok && !seen[otherIdx] {
					seen[otherIdx] = true
					hubSpokes[hub] = append(hubSpokes[hub], otherIdx)
				}
			}
		}
		maxSpokes := 1
		for _, spokes := range hubSpokes {
			if len(spokes) > maxSpokes {
				maxSpokes = len(spokes)
			}
		}
		canvasHalf := side/2 - margin
		sinB := math.Sin(math.Pi / float64(len(hubs)))
		sinA := math.Sin(math.Pi / float64(maxIntPair(maxSpokes, 2)))
		hubRing := (canvasHalf + 3) / (1 + sinB)
		radialLeft := canvasHalf - hubRing
		zoneRadius := math.Min(38, (radialLeft*sinA-3)/(1+sinA))
		if zoneRadius < 8 {
			zoneRadius = 8
		}
		spokeRing := math.Max(radialLeft-zoneRadius, 28+zoneRadius+3)
		layout.ZoneRadius = int(math.Round(zoneRadius))

		for hubPos, hubIndex := range hubs {
			ang := -math.Pi/2 + float64(hubPos)*2*math.Pi/float64(len(hubs))
			hx := cx + math.Cos(ang)*hubRing
			hy := cy + math.Sin(ang)*hubRing
			layout.Positions[zones[hubIndex].Name] = image.Pt(int(hx), int(hy))
			spokes := hubSpokes[zones[hubIndex].Name]
			for i, idx := range spokes {
				sa := ang + float64(i)*2*math.Pi/float64(len(spokes))
				sx := hx + math.Cos(sa)*spokeRing
				sy := hy + math.Sin(sa)*spokeRing
				layout.Positions[zones[idx].Name] = image.Pt(int(sx), int(sy))
			}
		}
		// Place orphan zones at centre.
		for _, zone := range zones {
			if _, ok := layout.Positions[zone.Name]; !ok {
				layout.Positions[zone.Name] = image.Pt(int(cx), int(cy))
			}
		}
	} else {
		// Single ring with optional centre Hub.
		hubIdx := implicitHubIdx
		if hubIdx < 0 {
			for i, zone := range zones {
				if zone.Name == "Hub" {
					hubIdx = i
					break
				}
			}
		}
		var outer []int
		for i := range zones {
			if i != hubIdx {
				outer = append(outer, i)
			}
		}
		n := len(outer)
		if n == 0 {
			n = 1
		}
		ringR0 := side/2 - margin
		chord := 2 * ringR0 * math.Sin(math.Pi/float64(maxIntPair(n, 1)))
		zoneRadius := math.Min(38, (chord-6)/2)
		if zoneRadius < 8 {
			zoneRadius = 8
		}
		ringR := math.Min(ringR0, side/2-zoneRadius-margin)
		if hubIdx >= 0 {
			ringR = math.Max(ringR, 28+zoneRadius+6)
		}
		layout.ZoneRadius = int(math.Round(zoneRadius))

		if hubIdx >= 0 {
			layout.Positions[zones[hubIdx].Name] = image.Pt(int(cx), int(cy))
		}
		if len(zones) == 1 {
			layout.Positions[zones[0].Name] = image.Pt(int(cx), int(cy))
		} else {
			for i, idx := range outer {
				ang := -math.Pi/2 + float64(i)*2*math.Pi/float64(n)
				x := cx + math.Cos(ang)*ringR
				y := cy + math.Sin(ang)*ringR
				layout.Positions[zones[idx].Name] = image.Pt(int(x), int(y))
			}
		}
	}

	// Build PreviewZones.
	for _, zone := range zones {
		pos, ok := layout.Positions[zone.Name]
		if !ok {
			continue
		}
		isHub := strings.EqualFold(zone.Name, "Hub") || strings.HasPrefix(zone.Name, "Hub-")
		if implicitHubIdx >= 0 && zone.Name == zones[implicitHubIdx].Name {
			isHub = true
		}
		previewZone := PreviewZone{
			Name:   zone.Name,
			Letter: ExtractZoneLetter(zone.Name),
			Center: pos,
			Tier:   ClassifyZoneTier(zone),
			IsHub:  isHub,
		}
		previewZone.IsPlayer = strings.HasPrefix(zone.Name, "Spawn-")
		for _, mainObject := range zone.MainObjects {
			if strings.EqualFold(mainObject.Type, "Spawn") {
				previewZone.HasCastle = true
				previewZone.Castles++
				// Extract player number from Spawn field (e.g. "Player1" → 1).
				if strings.HasPrefix(mainObject.Spawn, "Player") {
					digits := mainObject.Spawn[len("Player"):]
					for _, char := range digits {
						if char >= '0' && char <= '9' {
							previewZone.Owner = previewZone.Owner*10 + int(char-'0')
						}
					}
				}
			} else if strings.EqualFold(mainObject.Type, "City") {
				previewZone.HasCastle = true
				previewZone.Castles++
			}
		}
		layout.Zones = append(layout.Zones, previewZone)
	}
	// Connections — endpoints may use either Zone.Name or Zone.Letter,
	// so resolve through zoneNameByKey before looking up positions.
	for _, conn := range variant.Connections {
		aName, ok1 := resolveKey(conn.From)
		bName, ok2 := resolveKey(conn.To)
		if !ok1 || !ok2 {
			continue
		}
		a, okA := layout.Positions[aName]
		b, okB := layout.Positions[bName]
		if !okA || !okB {
			continue
		}
		isPortal := len(conn.PortalPlacementRulesFrom) > 0 ||
			len(conn.PortalPlacementRulesTo) > 0 ||
			conn.ConnectionType == "Portal"
		layout.Connections = append(layout.Connections, PreviewConnection{A: a, B: b, Portal: isPortal})
	}
	return layout
}

// ExtractZoneLetter returns the trailing letter portion of a zone name like
// "Spawn-A" → "A" or "Neutral-C" → "C". Plain names (e.g. "Hub") pass through.
func ExtractZoneLetter(zoneName string) string {
	if strings.HasPrefix(zoneName, "Spawn-") {
		return strings.TrimPrefix(zoneName, "Spawn-")
	}
	if strings.HasPrefix(zoneName, "Neutral-") {
		return strings.TrimPrefix(zoneName, "Neutral-")
	}
	return zoneName
}

// ClassifyZoneTier guesses a neutral zone's tier from its layout name and
// falls back to keyword matching on the zone name.
// (Mirrors C# SideLayoutName / TreasureLayoutName / CenterLayoutName.)
func ClassifyZoneTier(zone models.RmgZone) int {
	if strings.HasPrefix(zone.Name, "Spawn-") {
		return 0
	}
	layout := strings.ToLower(zone.Layout)
	switch {
	case strings.Contains(layout, "sides"):
		return 1
	case strings.Contains(layout, "treasure"):
		return 2
	case strings.Contains(layout, "center"):
		return 3
	}
	// Heuristic by name.
	name := strings.ToLower(zone.Name)
	switch {
	case strings.Contains(name, "low") || strings.Contains(name, "side"):
		return 1
	case strings.Contains(name, "med") || strings.Contains(name, "treasure"):
		return 2
	case strings.Contains(name, "high") || strings.Contains(name, "center") || strings.Contains(name, "core"):
		return 3
	}
	return 1
}

func maxIntPair(a, b int) int {
	if a > b {
		return a
	}
	return b
}
