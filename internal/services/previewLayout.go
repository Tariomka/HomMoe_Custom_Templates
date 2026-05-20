package services

import (
	"image"
	"math"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
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

	// Phase 7 — Balanced topology: concentric rings by zone tier, with
	// per-cluster canvas split when isolated clusters are detected
	// (Tournament + Balanced). Skip when the variant collapses to a
	// degenerate single-zone or hub-only layout, falling through to the
	// existing ring path.
	implicitHubIdx := -1
	if topology == models.MapTopology(generator.TopologyBalanced) && len(zones) > 1 {
		layoutBalanced(&layout, variant, zones, resolveKey, cx, cy, side, margin)
	} else {

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

// layoutBalanced lays out a Balanced-topology variant as concentric rings
// by zone tier (player outermost, then low / medium / high neutrals), with
// per-cluster canvas splitting for tournament-balanced (two disconnected
// clusters). It writes positions directly into layout.Positions and sets
// layout.ZoneRadius. Connections are routed straight (renderer handles
// drawing).
func layoutBalanced(layout *PreviewLayout, variant models.RmgVariant, zones []models.RmgZone, resolveKey func(string) (string, bool), cx, cy, side, margin float64) {
	// Floodfill clusters via connections so tournament-balanced (2 isolated
	// clusters) is auto-detected without needing to look at topology again.
	adj := map[string][]string{}
	for _, conn := range variant.Connections {
		a, ok1 := resolveKey(conn.From)
		b, ok2 := resolveKey(conn.To)
		if !ok1 || !ok2 {
			continue
		}
		adj[a] = append(adj[a], b)
		adj[b] = append(adj[b], a)
	}
	clusterOf := map[string]int{}
	clusters := [][]int{}
	for i, z := range zones {
		if _, seen := clusterOf[z.Name]; seen {
			continue
		}
		var members []int
		stack := []string{z.Name}
		idxByName := map[string]int{}
		for j, zz := range zones {
			idxByName[zz.Name] = j
		}
		for len(stack) > 0 {
			n := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if _, ok := clusterOf[n]; ok {
				continue
			}
			clusterOf[n] = len(clusters)
			if idx, ok := idxByName[n]; ok {
				members = append(members, idx)
			}
			stack = append(stack, adj[n]...)
		}
		if len(members) == 0 {
			members = append(members, i)
			clusterOf[z.Name] = len(clusters)
		}
		clusters = append(clusters, members)
	}

	// Choose a per-cluster centre & radius. For 2 clusters, split canvas
	// vertically; for 1, use the whole canvas; for >2 fall back to the
	// largest cluster centred.
	type clusterRegion struct {
		members []int
		cx, cy  float64
		maxR    float64
	}
	var regions []clusterRegion
	switch {
	case len(clusters) == 1:
		regions = []clusterRegion{{members: clusters[0], cx: cx, cy: cy, maxR: math.Min(cx, cy) - margin}}
	case len(clusters) == 2:
		halfW := side / 2
		regions = []clusterRegion{
			{members: clusters[0], cx: halfW / 2, cy: cy, maxR: math.Min(halfW/2, cy) - margin},
			{members: clusters[1], cx: halfW + halfW/2, cy: cy, maxR: math.Min(halfW/2, cy) - margin},
		}
	default:
		// Many clusters: stack each cluster on the largest-area square.
		biggest := 0
		for i, c := range clusters {
			if len(c) > len(clusters[biggest]) {
				biggest = i
			}
		}
		regions = []clusterRegion{{members: clusters[biggest], cx: cx, cy: cy, maxR: math.Min(cx, cy) - margin}}
		for i, c := range clusters {
			if i != biggest {
				regions = append(regions, clusterRegion{members: c, cx: cx, cy: cy, maxR: math.Min(cx, cy) - margin})
			}
		}
	}

	// Concentric rings: tier 0 (player) outermost; tier 1/2/3 inward. The
	// player ring sits at the cluster's max radius; subsequent rings shrink
	// by ringStep. The zone radius is sized off the busiest ring's chord.
	zoneRadius := 38.0
	for _, region := range regions {
		// Bucket cluster members by tier.
		byTier := map[int][]int{}
		for _, idx := range region.members {
			byTier[ClassifyZoneTier(zones[idx])] = append(byTier[ClassifyZoneTier(zones[idx])], idx)
		}
		// Override: player zones always tier-0 (outer ring).
		for _, idx := range region.members {
			if strings.HasPrefix(zones[idx].Name, "Spawn-") {
				// Remove from whatever tier it landed in, force tier 0.
				oldT := ClassifyZoneTier(zones[idx])
				if oldT != 0 {
					rebucket := byTier[oldT][:0]
					for _, ii := range byTier[oldT] {
						if ii != idx {
							rebucket = append(rebucket, ii)
						}
					}
					byTier[oldT] = rebucket
					byTier[0] = append(byTier[0], idx)
				}
			}
		}
		// Determine present tiers in order 0,1,2,3.
		var presentTiers []int
		for _, t := range []int{0, 1, 2, 3} {
			if len(byTier[t]) > 0 {
				presentTiers = append(presentTiers, t)
			}
		}
		if len(presentTiers) == 0 {
			continue
		}
		// Per-ring radius: evenly spaced from region.maxR down to ~30% of it.
		ringR := make(map[int]float64, len(presentTiers))
		for i, t := range presentTiers {
			frac := 1.0
			if len(presentTiers) > 1 {
				frac = 1.0 - 0.65*float64(i)/float64(len(presentTiers)-1)
			}
			ringR[t] = region.maxR * frac
		}
		// Pick zoneRadius based on tightest ring.
		minChord := math.MaxFloat64
		for _, t := range presentTiers {
			n := len(byTier[t])
			if n < 2 {
				continue
			}
			chord := 2 * ringR[t] * math.Sin(math.Pi/float64(n))
			if chord < minChord {
				minChord = chord
			}
		}
		if minChord != math.MaxFloat64 {
			r := math.Min(38, (minChord-6)/2)
			if r < 8 {
				r = 8
			}
			if r < zoneRadius {
				zoneRadius = r
			}
		}
		// Place each tier's zones at equal angular intervals.
		for _, t := range presentTiers {
			members := byTier[t]
			n := len(members)
			// Deterministic ordering by name for stability.
			sortByName(zones, members)
			phase := -math.Pi / 2
			if t%2 == 1 {
				phase += math.Pi / float64(maxIntPair(n, 1))
			}
			for i, idx := range members {
				ang := phase + float64(i)*2*math.Pi/float64(maxIntPair(n, 1))
				if n == 1 {
					layout.Positions[zones[idx].Name] = image.Pt(int(region.cx), int(region.cy))
					continue
				}
				x := region.cx + math.Cos(ang)*ringR[t]
				y := region.cy + math.Sin(ang)*ringR[t]
				layout.Positions[zones[idx].Name] = image.Pt(int(x), int(y))
			}
		}
	}
	layout.ZoneRadius = int(math.Round(zoneRadius))
}

// sortByName sorts a slice of zone indices in-place by the underlying
// zone name (ascending) — used for deterministic Balanced layout output.
func sortByName(zones []models.RmgZone, indices []int) {
	for i := 1; i < len(indices); i++ {
		for j := i; j > 0 && zones[indices[j-1]].Name > zones[indices[j]].Name; j-- {
			indices[j-1], indices[j] = indices[j], indices[j-1]
		}
	}
}
