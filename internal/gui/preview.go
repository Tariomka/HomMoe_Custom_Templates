package gui

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// Preview palette — ported from TemplatePreviewPngWriter colours.
var (
	prevBg         = color.NRGBA{R: 0x1C, G: 0x16, B: 0x10, A: 0xFF}
	prevFrame      = color.NRGBA{R: 0x8F, G: 0x73, B: 0x3F, A: 0xFF}
	prevBronzeFill = color.NRGBA{R: 0x65, G: 0x43, B: 0x21, A: 0xFF}
	prevBronzeEdge = color.NRGBA{R: 0xCD, G: 0x7F, B: 0x32, A: 0xFF}
	prevSilverFill = color.NRGBA{R: 0x48, G: 0x4C, B: 0x50, A: 0xFF}
	prevSilverEdge = color.NRGBA{R: 0xC0, G: 0xC0, B: 0xC0, A: 0xFF}
	prevGoldFill   = color.NRGBA{R: 0x78, G: 0x5A, B: 0x14, A: 0xFF}
	prevGoldEdge   = color.NRGBA{R: 0xFF, G: 0xD2, B: 0x32, A: 0xFF}
	prevSpawnFill  = color.NRGBA{R: 0x2A, G: 0x5A, B: 0x32, A: 0xFF}
	prevSpawnEdge  = color.NRGBA{R: 0x64, G: 0xC8, B: 0x78, A: 0xFF}
	prevHubFill    = color.NRGBA{R: 0x37, G: 0x50, B: 0x5F, A: 0xFF}
	prevHubEdge    = color.NRGBA{R: 0x82, G: 0xB4, B: 0xC8, A: 0xFF}
	prevDirectLine = color.NRGBA{R: 0xB4, G: 0x91, B: 0x3C, A: 0xFF}
	prevPortalLine = color.NRGBA{R: 0x5A, G: 0xAA, B: 0xD2, A: 0xB4}
	prevConnEdge   = color.NRGBA{R: 0x40, G: 0x35, B: 0x20, A: 0xFF}
)

// previewState holds the layout cache + buttons for the preview panel.
type previewState struct {
	btnSavePNG  widget.Clickable
	btnRefresh  widget.Clickable
	pngStatus   string
	pngStatusOK bool
}

// — Layout —

// layoutPreviewPanel renders the right-hand preview area. Returns empty
// dimensions when there's nothing to show (so the caller can omit it).
func (s *State) layoutPreviewPanel(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if s.btnSavePreview.Clicked(gtx) {
		s.savePreviewPNG()
	}
	tmpl := s.lastTemplate

	return borderedPanel(gtx, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				lbl := material.H6(th, "Preview")
				lbl.Color = colGold
				lbl.Font = font.Font{Weight: font.SemiBold}
				lbl.TextSize = unit.Sp(15)
				return lbl.Layout(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				name := "(no template generated yet)"
				if tmpl != nil {
					name = tmpl.Name
				}
				lbl := material.Body2(th, name)
				lbl.Color = colTextDim
				lbl.TextSize = unit.Sp(11)
				lbl.MaxLines = 1
				lbl.Truncator = "…"
				return layout.Inset{Top: unit.Dp(2), Bottom: unit.Dp(6)}.Layout(gtx, lbl.Layout)
			}),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				if tmpl == nil {
					return s.layoutPreviewEmpty(gtx, th)
				}
				return s.layoutPreviewCanvas(gtx, th, tmpl)
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return s.layoutPreviewLegend(gtx, th) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
						if s.preview.pngStatus == "" {
							return layout.Dimensions{}
						}
						col := colTextDim
						if !s.preview.pngStatusOK {
							col = colError
						}
						lbl := material.Body2(th, s.preview.pngStatus)
						lbl.Color = col
						lbl.TextSize = unit.Sp(11)
						lbl.MaxLines = 2
						return lbl.Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return toolbarButton{Text: "🖼  Save PNG", Click: &s.btnSavePreview, Disabled: tmpl == nil}.Layout(gtx, th)
					}),
				)
			}),
		)
	})
}

func (s *State) layoutPreviewEmpty(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		lbl := material.Body2(th, "Press \"Generate Template\" to see the map layout.")
		lbl.Color = colTextDim
		lbl.TextSize = unit.Sp(12)
		lbl.Alignment = text.Middle
		return lbl.Layout(gtx)
	})
}

// layoutPreviewCanvas draws the actual map preview at the largest possible
// square fitting inside the available area.
func (s *State) layoutPreviewCanvas(gtx layout.Context, th *material.Theme, tmpl *models.RmgTemplate) layout.Dimensions {
	maxX := gtx.Constraints.Max.X
	maxY := gtx.Constraints.Max.Y
	side := maxX
	if maxY < side {
		side = maxY
	}
	if side < 80 {
		side = 80
	}
	sz := image.Pt(side, side)

	// Background.
	rect := image.Rectangle{Max: sz}
	paint.FillShape(gtx.Ops, prevBg, clip.Rect(rect).Op())

	// Frame.
	r := gtx.Dp(unit.Dp(6))
	frame := image.Rectangle{Min: image.Pt(4, 4), Max: image.Pt(side-4, side-4)}
	paint.FillShape(gtx.Ops, prevFrame, clip.Stroke{
		Path:  clip.UniformRRect(frame, r).Path(gtx.Ops),
		Width: 2,
	}.Op())

	topology := s.sf.Topology
	pl := buildPreviewLayout(tmpl, topology, float64(side))
	if len(pl.Positions) == 0 {
		drawCenteredText(gtx, th, sz, tmpl.Name, 18, colText)
		return layout.Dimensions{Size: sz}
	}

	// Connections beneath zones.
	for _, c := range pl.Connections {
		drawConnection(gtx, c, pl.ZoneRadius)
	}
	// Non-spawn zones first, then spawn zones on top.
	for _, z := range pl.Zones {
		if z.IsPlayer {
			continue
		}
		drawPreviewZone(gtx, th, z, pl.ZoneRadius)
	}
	for _, z := range pl.Zones {
		if !z.IsPlayer {
			continue
		}
		drawPreviewZone(gtx, th, z, pl.ZoneRadius)
	}

	return layout.Dimensions{Size: sz}
}

func (s *State) layoutPreviewLegend(gtx layout.Context, th *material.Theme) layout.Dimensions {
	type item struct {
		Color color.NRGBA
		Label string
	}
	items := []item{
		{prevSpawnEdge, "Player"},
		{prevBronzeEdge, "Bronze"},
		{prevSilverEdge, "Silver"},
		{prevGoldEdge, "Gold"},
		{prevHubEdge, "Hub"},
	}
	children := make([]layout.FlexChild, 0, len(items)*2)
	for i, it := range items {
		it := it
		if i > 0 {
			children = append(children, layout.Rigid(layout.Spacer{Width: unit.Dp(8)}.Layout))
		}
		children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					side := gtx.Dp(unit.Dp(10))
					rect := image.Rect(0, 0, side, side)
					paint.FillShape(gtx.Ops, it.Color, clip.UniformRRect(rect, side/2).Op(gtx.Ops))
					return layout.Dimensions{Size: rect.Max}
				}),
				layout.Rigid(layout.Spacer{Width: unit.Dp(4)}.Layout),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					lbl := material.Body2(th, it.Label)
					lbl.Color = colTextDim
					lbl.TextSize = unit.Sp(10)
					return lbl.Layout(gtx)
				}),
			)
		}))
	}
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
}

// — Geometry —

type previewZone struct {
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

type previewConn struct {
	A, B   image.Point
	Portal bool
}

type previewLayout struct {
	Positions   map[string]image.Point
	Zones       []previewZone
	Connections []previewConn
	ZoneRadius  int
}

func buildPreviewLayout(t *models.RmgTemplate, topology models.MapTopology, side float64) previewLayout {
	pl := previewLayout{Positions: map[string]image.Point{}}
	if t == nil || len(t.Variants) == 0 {
		return pl
	}
	v := t.Variants[0]
	zones := v.Zones
	if len(zones) == 0 {
		return pl
	}

	// Build name lookup. Zone names are like "Spawn-A", "Neutral-C", "Hub".
	zoneNameByKey := make(map[string]string, len(zones)*2)
	for _, z := range zones {
		zoneNameByKey[z.Name] = z.Name
	}
	resolveKey := func(key string) (string, bool) {
		n, ok := zoneNameByKey[key]
		return n, ok
	}

	// Layout: hub-and-spoke if multiple "Hub-*" zones; ring otherwise (with optional center "Hub").
	cx := side / 2
	cy := side / 2
	margin := 24.0

	// Detect multi-hub (tournament): zones literally named "Hub-*".
	var hubs []int
	for i, z := range zones {
		if strings.HasPrefix(z.Name, "Hub") {
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
		for _, z := range zones {
			if strings.HasPrefix(z.Name, "Spawn-") {
				playerNames[z.Name] = true
			}
		}
		// neighbours[zoneName] = set of connected zone names.
		neighbours := make(map[string]map[string]bool, len(zones))
		for _, c := range v.Connections {
			a, ok1 := resolveKey(c.From)
			b, ok2 := resolveKey(c.To)
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
		for i, z := range zones {
			if strings.HasPrefix(z.Name, "Spawn-") {
				continue
			}
			nbrs := neighbours[z.Name]
			if len(nbrs) < 2 {
				continue
			}
			// Must connect to every player zone.
			connectsAllPlayers := len(playerNames) > 0
			for pn := range playerNames {
				if !nbrs[pn] {
					connectsAllPlayers = false
					break
				}
			}
			if !connectsAllPlayers {
				continue
			}
			if len(nbrs) > bestDeg {
				bestDeg = len(nbrs)
				implicitHubIdx = i
			}
		}
	}
	if len(hubs) >= 2 {
		// Build spoke lists.
		zoneIdx := map[string]int{}
		for i, z := range zones {
			zoneIdx[z.Name] = i
		}
		hubSpokes := make(map[string][]int, len(hubs))
		for _, hi := range hubs {
			hub := zones[hi].Name
			seen := map[int]bool{}
			for _, c := range v.Connections {
				other := ""
				switch {
				case c.From == hub:
					other = c.To
				case c.To == hub:
					other = c.From
				}
				if other == "" {
					continue
				}
				otherName, ok := resolveKey(other)
				if !ok {
					continue
				}
				if oi, ok := zoneIdx[otherName]; ok && !seen[oi] {
					seen[oi] = true
					hubSpokes[hub] = append(hubSpokes[hub], oi)
				}
			}
		}
		maxSpokes := 1
		for _, sp := range hubSpokes {
			if len(sp) > maxSpokes {
				maxSpokes = len(sp)
			}
		}
		canvasHalf := side/2 - margin
		sinB := math.Sin(math.Pi / float64(len(hubs)))
		sinA := math.Sin(math.Pi / float64(max2(maxSpokes, 2)))
		hubRing := (canvasHalf + 3) / (1 + sinB)
		radialLeft := canvasHalf - hubRing
		zoneRadius := math.Min(38, (radialLeft*sinA-3)/(1+sinA))
		if zoneRadius < 8 {
			zoneRadius = 8
		}
		spokeRing := math.Max(radialLeft-zoneRadius, 28+zoneRadius+3)
		pl.ZoneRadius = int(math.Round(zoneRadius))

		for h, hi := range hubs {
			ang := -math.Pi/2 + float64(h)*2*math.Pi/float64(len(hubs))
			hx := cx + math.Cos(ang)*hubRing
			hy := cy + math.Sin(ang)*hubRing
			pl.Positions[zones[hi].Name] = image.Pt(int(hx), int(hy))
			sp := hubSpokes[zones[hi].Name]
			for i, idx := range sp {
				sa := ang + float64(i)*2*math.Pi/float64(len(sp))
				sx := hx + math.Cos(sa)*spokeRing
				sy := hy + math.Sin(sa)*spokeRing
				pl.Positions[zones[idx].Name] = image.Pt(int(sx), int(sy))
			}
		}
		// Place orphan zones at centre.
		for _, z := range zones {
			if _, ok := pl.Positions[z.Name]; !ok {
				pl.Positions[z.Name] = image.Pt(int(cx), int(cy))
			}
		}
	} else {
		// Single ring with optional centre Hub.
		hubIdx := implicitHubIdx
		if hubIdx < 0 {
			for i, z := range zones {
				if z.Name == "Hub" {
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
		chord := 2 * ringR0 * math.Sin(math.Pi/float64(max2(n, 1)))
		zoneRadius := math.Min(38, (chord-6)/2)
		if zoneRadius < 8 {
			zoneRadius = 8
		}
		ringR := math.Min(ringR0, side/2-zoneRadius-margin)
		if hubIdx >= 0 {
			ringR = math.Max(ringR, 28+zoneRadius+6)
		}
		pl.ZoneRadius = int(math.Round(zoneRadius))

		if hubIdx >= 0 {
			pl.Positions[zones[hubIdx].Name] = image.Pt(int(cx), int(cy))
		}
		if len(zones) == 1 {
			pl.Positions[zones[0].Name] = image.Pt(int(cx), int(cy))
		} else {
			for i, idx := range outer {
				ang := -math.Pi/2 + float64(i)*2*math.Pi/float64(n)
				x := cx + math.Cos(ang)*ringR
				y := cy + math.Sin(ang)*ringR
				pl.Positions[zones[idx].Name] = image.Pt(int(x), int(y))
			}
		}
	}

	// Build previewZones.
	for _, z := range zones {
		pos, ok := pl.Positions[z.Name]
		if !ok {
			continue
		}
		isHub := strings.EqualFold(z.Name, "Hub") || strings.HasPrefix(z.Name, "Hub-")
		if implicitHubIdx >= 0 && z.Name == zones[implicitHubIdx].Name {
			isHub = true
		}
		letter := extractLetter(z.Name)
		pz := previewZone{
			Name:   z.Name,
			Letter: letter,
			Center: pos,
			Tier:   classifyTier(z),
			IsHub:  isHub,
		}
		pz.IsPlayer = strings.HasPrefix(z.Name, "Spawn-")
		for _, mo := range z.MainObjects {
			if strings.EqualFold(mo.Type, "Spawn") || strings.EqualFold(mo.Type, "City") {
				pz.HasCastle = true
				pz.Castles++
			}
		}
		pl.Zones = append(pl.Zones, pz)
	}
	// Connections — endpoints may use either Zone.Name or Zone.Letter,
	// so resolve through zoneNameByKey before looking up positions.
	for _, c := range v.Connections {
		aName, ok1 := resolveKey(c.From)
		bName, ok2 := resolveKey(c.To)
		if !ok1 || !ok2 {
			continue
		}
		a, okA := pl.Positions[aName]
		b, okB := pl.Positions[bName]
		if !okA || !okB {
			continue
		}
		isPortal := len(c.PortalPlacementRulesFrom) > 0 || len(c.PortalPlacementRulesTo) > 0 || c.ConnectionType == "Portal"
		pl.Connections = append(pl.Connections, previewConn{A: a, B: b, Portal: isPortal})
	}
	return pl
}

func max2(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// classifyTier guesses a neutral zone's tier from its layout template name
// (matches C# SideLayoutName / TreasureLayoutName / CenterLayoutName).
func extractLetter(zoneName string) string {
	if strings.HasPrefix(zoneName, "Spawn-") {
		return strings.TrimPrefix(zoneName, "Spawn-")
	}
	if strings.HasPrefix(zoneName, "Neutral-") {
		return strings.TrimPrefix(zoneName, "Neutral-")
	}
	return zoneName
}

func classifyTier(z models.RmgZone) int {
	if strings.HasPrefix(z.Name, "Spawn-") {
		return 0
	}
	t := strings.ToLower(z.Layout)
	switch {
	case strings.Contains(t, "sides"):
		return 1
	case strings.Contains(t, "treasure"):
		return 2
	case strings.Contains(t, "center"):
		return 3
	}
	// Heuristic by name.
	n := strings.ToLower(z.Name)
	switch {
	case strings.Contains(n, "low") || strings.Contains(n, "side"):
		return 1
	case strings.Contains(n, "med") || strings.Contains(n, "treasure"):
		return 2
	case strings.Contains(n, "high") || strings.Contains(n, "center") || strings.Contains(n, "core"):
		return 3
	}
	return 1
}

// — Drawing primitives —

func drawConnection(gtx layout.Context, c previewConn, zoneR int) {
	// Trim line to circle edges so it doesn't overlap zone fill.
	dx := float64(c.B.X - c.A.X)
	dy := float64(c.B.Y - c.A.Y)
	d := math.Hypot(dx, dy)
	if d < 1 {
		return
	}
	ux := dx / d
	uy := dy / d
	r := float64(zoneR)
	ax := float64(c.A.X) + ux*r
	ay := float64(c.A.Y) + uy*r
	bx := float64(c.B.X) - ux*r
	by := float64(c.B.Y) - uy*r

	col := prevDirectLine
	w := float32(gtx.Dp(unit.Dp(2.0)))
	if c.Portal {
		col = prevPortalLine
		w = float32(gtx.Dp(unit.Dp(1.5)))
	}
	drawLine(gtx, image.Pt(int(ax), int(ay)), image.Pt(int(bx), int(by)), w, col)
}

func drawLine(gtx layout.Context, a, b image.Point, width float32, col color.NRGBA) {
	var path clip.Path
	path.Begin(gtx.Ops)
	path.MoveTo(f32Pt(a))
	path.LineTo(f32Pt(b))
	paint.FillShape(gtx.Ops, col, clip.Stroke{
		Path:  path.End(),
		Width: width,
	}.Op())
}

func f32Pt(p image.Point) (out f32Point) {
	out.X = float32(p.X)
	out.Y = float32(p.Y)
	return
}

// f32Point is a tiny shim so we don't import gioui.org/f32 just for two fields.
type f32Point = struct {
	X, Y float32
}

func drawPreviewZone(gtx layout.Context, th *material.Theme, z previewZone, zoneR int) {
	r := zoneR
	if z.IsHub && r < 28 {
		r = 28
	}
	cx, cy := z.Center.X, z.Center.Y
	rect := image.Rect(cx-r, cy-r, cx+r, cy+r)

	fill, edge := zoneColors(z)
	circle := clip.UniformRRect(rect, r).Op(gtx.Ops)
	paint.FillShape(gtx.Ops, fill, circle)
	paint.FillShape(gtx.Ops, edge, clip.Stroke{
		Path:  clip.UniformRRect(rect, r).Path(gtx.Ops),
		Width: float32(gtx.Dp(unit.Dp(2))),
	}.Op())

	label := zoneLabel(z)
	if label != "" {
		drawCenteredText(gtx, th, image.Pt(cx, cy), label, 12, color.NRGBA{R: 0xF8, G: 0xE8, B: 0xC0, A: 0xFF})
	}
	if z.HasCastle && z.Castles > 0 {
		// Small badge in lower right.
		bx := cx + r/2
		by := cy + r/2
		drawCenteredText(gtx, th, image.Pt(bx, by), fmt.Sprintf("⌂%d", z.Castles), 10, color.NRGBA{R: 0xFF, G: 0xE8, B: 0x90, A: 0xFF})
	}
}

func zoneColors(z previewZone) (fill, edge color.NRGBA) {
	switch {
	case z.IsPlayer:
		return prevSpawnFill, prevSpawnEdge
	case z.IsHub:
		return prevHubFill, prevHubEdge
	}
	switch z.Tier {
	case 3:
		return prevGoldFill, prevGoldEdge
	case 2:
		return prevSilverFill, prevSilverEdge
	default:
		return prevBronzeFill, prevBronzeEdge
	}
}

func zoneLabel(z previewZone) string {
	if z.IsPlayer {
		if z.Owner > 0 {
			return fmt.Sprintf("P%d", z.Owner)
		}
		// Spawn-1 / Spawn-2 → "P1"…
		if strings.HasPrefix(z.Name, "Spawn-") {
			return "P" + z.Name[len("Spawn-"):]
		}
		return z.Letter
	}
	if z.IsHub {
		return "Hub"
	}
	switch z.Tier {
	case 3:
		return "G"
	case 2:
		return "S"
	default:
		return "B"
	}
}

// drawCenteredText draws text centered on the given canvas point.
func drawCenteredText(gtx layout.Context, th *material.Theme, ctx interface{}, txt string, sizeSp int, col color.NRGBA) {
	// ctx may be image.Point (centre) or image.Point used as full-canvas size.
	centreOnPoint, isPoint := ctx.(image.Point)
	canvasSize, isSize := ctx.(image.Point)
	_ = canvasSize

	macro := op.Record(gtx.Ops)
	dims := func() layout.Dimensions {
		gtx2 := gtx
		gtx2.Constraints.Min = image.Point{}
		gtx2.Constraints.Max = image.Pt(1<<14, 1<<14)
		lbl := material.Label(th, unit.Sp(float32(sizeSp)), txt)
		lbl.Color = col
		lbl.Font = font.Font{Weight: font.SemiBold}
		return lbl.Layout(gtx2)
	}()
	call := macro.Stop()

	if !isPoint && !isSize {
		call.Add(gtx.Ops)
		return
	}
	tx := centreOnPoint.X - dims.Size.X/2
	ty := centreOnPoint.Y - dims.Size.Y/2
	stack := op.Offset(image.Pt(tx, ty)).Push(gtx.Ops)
	call.Add(gtx.Ops)
	stack.Pop()
}

// — PNG export —

// savePreviewPNG renders the current template into a software bitmap and
// writes it next to the .rmg.json output (matching the C# behaviour).
func (s *State) savePreviewPNG() {
	tmpl := s.lastTemplate
	if tmpl == nil {
		s.preview.pngStatus = "Generate a template first."
		s.preview.pngStatusOK = false
		return
	}
	dir := strings.TrimSpace(s.outputPath.Text())
	if dir == "" {
		s.preview.pngStatus = "Output directory is empty."
		s.preview.pngStatusOK = false
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		s.preview.pngStatus = "Cannot create directory: " + err.Error()
		s.preview.pngStatusOK = false
		return
	}
	name := sanitizeFilename(tmpl.Name)
	if name == "" {
		name = "Generated_Template"
	}
	out := filepath.Join(dir, name+".png")
	img := renderPreviewToImage(tmpl, s.sf.Topology, 700)
	f, err := os.Create(out)
	if err != nil {
		s.preview.pngStatus = "Create failed: " + err.Error()
		s.preview.pngStatusOK = false
		return
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		s.preview.pngStatus = "Encode failed: " + err.Error()
		s.preview.pngStatusOK = false
		return
	}
	s.preview.pngStatus = "Saved " + out
	s.preview.pngStatusOK = true
}

// renderPreviewToImage rasterises the layout into an *image.RGBA. It uses
// only the standard library — circles are filled scanline-by-scanline, lines
// via a simple Bresenham/DDA with a tiny brush.
func renderPreviewToImage(t *models.RmgTemplate, topology models.MapTopology, side int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, side, side))
	fillRect(img, img.Bounds(), prevBg)

	pl := buildPreviewLayout(t, topology, float64(side))
	if len(pl.Positions) == 0 {
		return img
	}
	r := pl.ZoneRadius

	// Frame.
	strokeRect(img, image.Rect(4, 4, side-4, side-4), 2, prevFrame)

	// Connections.
	for _, c := range pl.Connections {
		dx := float64(c.B.X - c.A.X)
		dy := float64(c.B.Y - c.A.Y)
		d := math.Hypot(dx, dy)
		if d < 1 {
			continue
		}
		ux := dx / d
		uy := dy / d
		ax := image.Pt(int(float64(c.A.X)+ux*float64(r)), int(float64(c.A.Y)+uy*float64(r)))
		bx := image.Pt(int(float64(c.B.X)-ux*float64(r)), int(float64(c.B.Y)-uy*float64(r)))
		col := prevDirectLine
		w := 3
		if c.Portal {
			col = prevPortalLine
			w = 2
		}
		drawThickLine(img, ax, bx, w, col)
	}
	// Zones.
	for _, z := range pl.Zones {
		if z.IsPlayer {
			continue
		}
		fill, edge := zoneColors(z)
		zr := r
		if z.IsHub && zr < 28 {
			zr = 28
		}
		fillCircle(img, z.Center, zr, fill)
		strokeCircle(img, z.Center, zr, 2, edge)
	}
	for _, z := range pl.Zones {
		if !z.IsPlayer {
			continue
		}
		fill, edge := zoneColors(z)
		fillCircle(img, z.Center, r, fill)
		strokeCircle(img, z.Center, r, 2, edge)
	}
	return img
}

// fillRect fills a rectangle with a single colour.
func fillRect(img *image.RGBA, r image.Rectangle, c color.NRGBA) {
	for y := r.Min.Y; y < r.Max.Y; y++ {
		for x := r.Min.X; x < r.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{R: c.R, G: c.G, B: c.B, A: c.A})
		}
	}
}

func strokeRect(img *image.RGBA, r image.Rectangle, w int, c color.NRGBA) {
	fillRect(img, image.Rect(r.Min.X, r.Min.Y, r.Max.X, r.Min.Y+w), c)
	fillRect(img, image.Rect(r.Min.X, r.Max.Y-w, r.Max.X, r.Max.Y), c)
	fillRect(img, image.Rect(r.Min.X, r.Min.Y, r.Min.X+w, r.Max.Y), c)
	fillRect(img, image.Rect(r.Max.X-w, r.Min.Y, r.Max.X, r.Max.Y), c)
}

// fillCircle paints a filled disc.
func fillCircle(img *image.RGBA, c image.Point, r int, col color.NRGBA) {
	r2 := r * r
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r2 {
				x := c.X + dx
				y := c.Y + dy
				if x < 0 || y < 0 || x >= img.Rect.Max.X || y >= img.Rect.Max.Y {
					continue
				}
				img.SetRGBA(x, y, color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A})
			}
		}
	}
}

// strokeCircle paints a circular outline of the given thickness.
func strokeCircle(img *image.RGBA, c image.Point, r, w int, col color.NRGBA) {
	rOuter2 := r * r
	rInner := r - w
	rInner2 := rInner * rInner
	if rInner2 < 0 {
		rInner2 = 0
	}
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			d2 := dx*dx + dy*dy
			if d2 <= rOuter2 && d2 >= rInner2 {
				x := c.X + dx
				y := c.Y + dy
				if x < 0 || y < 0 || x >= img.Rect.Max.X || y >= img.Rect.Max.Y {
					continue
				}
				img.SetRGBA(x, y, color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A})
			}
		}
	}
}

// drawThickLine draws a width-w line via DDA with a small square brush.
func drawThickLine(img *image.RGBA, a, b image.Point, w int, col color.NRGBA) {
	dx := b.X - a.X
	dy := b.Y - a.Y
	steps := dx
	if -dx > steps {
		steps = -dx
	}
	if dy > steps {
		steps = dy
	}
	if -dy > steps {
		steps = -dy
	}
	if steps <= 0 {
		fillCircle(img, a, w, col)
		return
	}
	xinc := float64(dx) / float64(steps)
	yinc := float64(dy) / float64(steps)
	x := float64(a.X)
	y := float64(a.Y)
	half := w / 2
	for i := 0; i <= steps; i++ {
		px := int(math.Round(x))
		py := int(math.Round(y))
		for oy := -half; oy <= half; oy++ {
			for ox := -half; ox <= half; ox++ {
				xx := px + ox
				yy := py + oy
				if xx < 0 || yy < 0 || xx >= img.Rect.Max.X || yy >= img.Rect.Max.Y {
					continue
				}
				img.SetRGBA(xx, yy, color.RGBA{R: col.R, G: col.G, B: col.B, A: col.A})
			}
		}
		x += xinc
		y += yinc
	}
}
