package previewLayout_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// ── helpers ──────────────────────────────────────────────────────────

func pos(x, y float64) *[2]float64 { p := [2]float64{x, y}; return &p }
func ringIdx(i int) *int           { return &i }

func zone(name string) entities.Zone { return entities.Zone{Name: name} }
func zonePos(name string, x, y float64) entities.Zone {
	z := zone(name)
	z.GeneratorPosition = pos(x, y)
	return z
}
func zoneRing(name string, ring int, x, y float64) entities.Zone {
	z := zonePos(name, x, y)
	z.GeneratorRing = ringIdx(ring)
	return z
}
func conn(from, to string) entities.Connection {
	return entities.Connection{From: from, To: to, ConnectionType: "Direct"}
}

func tmpl(zones []entities.Zone, conns []entities.Connection) *entities.RmgTemplate {
	return &entities.RmgTemplate{
		Variants: []entities.Variant{{Zones: zones, Connections: conns}},
	}
}

// ── BuildPreviewLayout: edge cases ──────────────────────────────────

func TestBuildPreviewLayout_NilTemplate(t *testing.T) {
	out := services.BuildPreviewLayout(nil, config.TopologyRing, 600)
	if len(out.Positions) != 0 || len(out.Zones) != 0 {
		t.Errorf("expected empty layout, got %+v", out)
	}
}

func TestBuildPreviewLayout_EmptyVariants(t *testing.T) {
	out := services.BuildPreviewLayout(&entities.RmgTemplate{}, config.TopologyRing, 600)
	if len(out.Positions) != 0 {
		t.Errorf("expected empty layout")
	}
}

func TestBuildPreviewLayout_EmptyZones(t *testing.T) {
	out := services.BuildPreviewLayout(tmpl(nil, nil), config.TopologyRing, 600)
	if len(out.Positions) != 0 {
		t.Errorf("expected empty layout")
	}
}

// ── BuildPreviewLayout: ring/default dispatch ────────────────────────

func TestBuildPreviewLayout_DefaultTopology_RingLayout(t *testing.T) {
	zones := []entities.Zone{zone("Spawn-A"), zone("Spawn-B"), zone("Neutral-C")}
	conns := []entities.Connection{conn("Spawn-A", "Spawn-B"), conn("Spawn-B", "Neutral-C")}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRing, 600)
	if len(out.Positions) != 3 {
		t.Errorf("got %d positions, want 3", len(out.Positions))
	}
	if out.ZoneRadius <= 0 {
		t.Error("expected positive zone radius")
	}
}

func TestBuildPreviewLayout_SingleZoneCentredOnCanvas(t *testing.T) {
	zones := []entities.Zone{zone("Spawn-A")}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyRing, 600)
	p, ok := out.Positions["Spawn-A"]
	if !ok {
		t.Fatal("missing position")
	}
	if p.X != 300 || p.Y != 300 {
		t.Errorf("center = (%d, %d), want (300,300)", p.X, p.Y)
	}
}

func TestBuildPreviewLayout_ExplicitHubGetsCentre(t *testing.T) {
	zones := []entities.Zone{zone("Hub"), zone("Spawn-A"), zone("Spawn-B")}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyRing, 600)
	if p, ok := out.Positions["Hub"]; !ok || p.X != 300 || p.Y != 300 {
		t.Errorf("Hub position = %+v, want (300,300)", p)
	}
}

func TestBuildPreviewLayout_ImplicitHubGetsCentre(t *testing.T) {
	// Neutral connected to both players, deg >= 2 → implicit hub.
	zones := []entities.Zone{zone("Spawn-A"), zone("Spawn-B"), zone("Neutral-H")}
	conns := []entities.Connection{conn("Neutral-H", "Spawn-A"), conn("Neutral-H", "Spawn-B")}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRing, 600)
	if p := out.Positions["Neutral-H"]; p.X != 300 || p.Y != 300 {
		t.Errorf("implicit hub position = %+v, want (300,300)", p)
	}
}

func TestBuildPreviewLayout_MultiHubFanLayout(t *testing.T) {
	zones := []entities.Zone{zone("Hub-A"), zone("Hub-B"), zone("Spawn-A"), zone("Spawn-B")}
	conns := []entities.Connection{
		conn("Hub-A", "Spawn-A"),
		conn("Hub-B", "Spawn-B"),
		conn("Hub-A", "Hub-B"),
	}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRing, 600)
	if len(out.Positions) != 4 {
		t.Errorf("expected 4 positions, got %d", len(out.Positions))
	}
}

// ── BuildPreviewLayout: scatter (Random) dispatch ────────────────────

func TestBuildPreviewLayout_RandomTopology_WithPositions(t *testing.T) {
	zones := []entities.Zone{
		zonePos("Spawn-A", 0.2, 0.2),
		zonePos("Spawn-B", 0.8, 0.8),
		zonePos("Neutral-C", 0.5, 0.5),
	}
	conns := []entities.Connection{conn("Spawn-A", "Neutral-C"), conn("Neutral-C", "Spawn-B")}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRandom, 600)
	if len(out.Positions) != 3 {
		t.Errorf("expected 3 positions, got %d", len(out.Positions))
	}
}

func TestBuildPreviewLayout_RandomTopology_NoConnections(t *testing.T) {
	// 3 zones with no edges → 3 components → cluster-strip skipped.
	zones := []entities.Zone{
		zonePos("Spawn-A", 0.1, 0.1),
		zonePos("Spawn-B", 0.9, 0.9),
		zonePos("Neutral-C", 0.5, 0.5),
	}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyRandom, 600)
	if len(out.Positions) != 3 {
		t.Errorf("expected 3 positions, got %d", len(out.Positions))
	}
}

func TestBuildPreviewLayout_RandomTopology_FallsBackWhenNoPositions(t *testing.T) {
	// No GeneratorPosition → falls through to ring path; 3 zones avoid strip.
	zones := []entities.Zone{zone("Spawn-A"), zone("Spawn-B"), zone("Neutral-C")}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyRandom, 600)
	if len(out.Positions) != 3 {
		t.Errorf("expected ring fallback, got %d positions", len(out.Positions))
	}
}

// ── BuildPreviewLayout: circles rings dispatch ──────────────────────

func TestBuildPreviewLayout_CirclesTopology_MultipleRings(t *testing.T) {
	zones := []entities.Zone{
		zoneRing("Spawn-A", 0, 0.1, 0.1),
		zoneRing("Spawn-B", 0, 0.9, 0.1),
		zoneRing("Neutral-C", 1, 0.5, 0.5),
	}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyCircles, 600)
	if len(out.Positions) != 3 {
		t.Errorf("expected 3 positions, got %d", len(out.Positions))
	}
}

func TestBuildPreviewLayout_CirclesTopology_SingleRingFallsBack(t *testing.T) {
	zones := []entities.Zone{
		zoneRing("Spawn-A", 0, 0.2, 0.2),
		zoneRing("Spawn-B", 0, 0.8, 0.8),
		zoneRing("Spawn-C", 0, 0.5, 0.5),
	}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyCircles, 600)
	if len(out.Positions) != 3 {
		t.Errorf("expected 3 positions, got %d", len(out.Positions))
	}
}

func TestBuildPreviewLayout_CirclesTopology_SingleZone(t *testing.T) {
	zones := []entities.Zone{zoneRing("Spawn-A", 0, 0.5, 0.5)}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyCircles, 600)
	if p := out.Positions["Spawn-A"]; p.X != 300 || p.Y != 300 {
		t.Errorf("expected centre, got %+v", p)
	}
}

// ── BuildPreviewLayout: connection rendering ─────────────────────────

func TestBuildPreviewLayout_DirectConnectionsAreCollected(t *testing.T) {
	zones := []entities.Zone{zone("Spawn-A"), zone("Spawn-B")}
	conns := []entities.Connection{conn("Spawn-A", "Spawn-B")}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRing, 600)
	if len(out.Connections) != 1 || out.Connections[0].Portal {
		t.Errorf("expected one direct connection, got %+v", out.Connections)
	}
}

func TestBuildPreviewLayout_PortalConnectionsFlagged(t *testing.T) {
	// Portal is structurally ignored, so we add a Direct backbone to keep all
	// zones connected (avoiding 2-component cluster-strip).
	zones := []entities.Zone{zone("Spawn-A"), zone("Spawn-B"), zone("Neutral-C")}
	conns := []entities.Connection{
		conn("Spawn-A", "Neutral-C"),
		conn("Neutral-C", "Spawn-B"),
		{From: "Spawn-A", To: "Spawn-B", ConnectionType: "Portal"},
	}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRing, 600)
	portals := 0
	for _, c := range out.Connections {
		if c.Portal {
			portals++
		}
	}
	if portals != 1 {
		t.Errorf("expected 1 portal, got %d (%+v)", portals, out.Connections)
	}
}

func TestBuildPreviewLayout_ConnectionWithUnknownZoneSkipped(t *testing.T) {
	zones := []entities.Zone{zone("Spawn-A"), zone("Spawn-B")}
	conns := []entities.Connection{conn("Spawn-A", "Missing-X")}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRing, 600)
	if len(out.Connections) != 0 {
		t.Errorf("expected no rendered connection, got %+v", out.Connections)
	}
}

// ── BuildPreviewLayout: two-cluster stripping (tournament) ───────────

func TestBuildPreviewLayout_TwoComponentsKeepsOnlyFirstCluster(t *testing.T) {
	// Two disjoint clusters → render only the first.
	zones := []entities.Zone{
		zone("Spawn-A"), zone("Neutral-X"),
		zone("Spawn-B"), zone("Neutral-Y"),
	}
	conns := []entities.Connection{
		conn("Spawn-A", "Neutral-X"),
		conn("Spawn-B", "Neutral-Y"),
	}
	out := services.BuildPreviewLayout(tmpl(zones, conns), config.TopologyRing, 600)
	if len(out.Positions) != 2 {
		t.Errorf("expected 2 positions after strip, got %d", len(out.Positions))
	}
}

// ── BuildPreviewLayout: zone classification side-effects ─────────────

func TestBuildPreviewLayout_PlayerZoneFlagged(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Spawn-A", MainObjects: []entities.MainObject{{Type: "Spawn", Spawn: "Player1"}}},
	}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyRing, 600)
	if len(out.Zones) != 1 || !out.Zones[0].IsPlayer {
		t.Errorf("expected player zone, got %+v", out.Zones)
	}
	if !out.Zones[0].HasCastle || out.Zones[0].Owner != 1 {
		t.Errorf("expected castle+owner=1, got %+v", out.Zones[0])
	}
}

func TestBuildPreviewLayout_CityMainObjectCounts(t *testing.T) {
	zones := []entities.Zone{
		{Name: "Neutral-Z", MainObjects: []entities.MainObject{{Type: "City"}, {Type: "City"}}},
	}
	out := services.BuildPreviewLayout(tmpl(zones, nil), config.TopologyRing, 600)
	if out.Zones[0].Castles != 2 {
		t.Errorf("expected 2 castles, got %d", out.Zones[0].Castles)
	}
}

// ── ExtractZoneLetter ────────────────────────────────────────────────

func TestExtractZoneLetter_SpawnPrefix(t *testing.T) {
	if got := services.ExtractZoneLetter("Spawn-AB"); got != "AB" {
		t.Errorf("got %q", got)
	}
}

func TestExtractZoneLetter_NeutralPrefix(t *testing.T) {
	if got := services.ExtractZoneLetter("Neutral-C"); got != "C" {
		t.Errorf("got %q", got)
	}
}

func TestExtractZoneLetter_PlainName(t *testing.T) {
	if got := services.ExtractZoneLetter("Hub"); got != "Hub" {
		t.Errorf("got %q", got)
	}
}

// ── ClassifyZoneTier ─────────────────────────────────────────────────

func TestClassifyZoneTier_SpawnZoneIsZero(t *testing.T) {
	if got := services.ClassifyZoneTier(entities.Zone{Name: "Spawn-A"}); got != 0 {
		t.Errorf("got %d, want 0", got)
	}
}

func TestClassifyZoneTier_GuardedPoolT5IsGold(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t5_item"}}
	if got := services.ClassifyZoneTier(z); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestClassifyZoneTier_GuardedPoolT4IsGold(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t4_item"}}
	if got := services.ClassifyZoneTier(z); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestClassifyZoneTier_GuardedPoolT3IsSilver(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t3_item"}}
	if got := services.ClassifyZoneTier(z); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestClassifyZoneTier_GuardedPoolT2IsBronze(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", GuardedContentPool: []string{"classic_template_pool_random_t2_item"}}
	if got := services.ClassifyZoneTier(z); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestClassifyZoneTier_UnguardedPoolFallback(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", UnguardedContentPool: []string{"classic_template_pool_random_unguarded_t3_item"}}
	if got := services.ClassifyZoneTier(z); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestClassifyZoneTier_LayoutSidesIsBronze(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", Layout: "zone_layout_sides"}
	if got := services.ClassifyZoneTier(z); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestClassifyZoneTier_LayoutTreasureIsSilver(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", Layout: "zone_layout_treasure_zone"}
	if got := services.ClassifyZoneTier(z); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestClassifyZoneTier_LayoutCenterIsGold(t *testing.T) {
	z := entities.Zone{Name: "Neutral-A", Layout: "zone_layout_center"}
	if got := services.ClassifyZoneTier(z); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestClassifyZoneTier_NameLowFallback(t *testing.T) {
	z := entities.Zone{Name: "Neutral-low"}
	if got := services.ClassifyZoneTier(z); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}

func TestClassifyZoneTier_NameMedFallback(t *testing.T) {
	z := entities.Zone{Name: "Neutral-med"}
	if got := services.ClassifyZoneTier(z); got != 2 {
		t.Errorf("got %d, want 2", got)
	}
}

func TestClassifyZoneTier_NameHighFallback(t *testing.T) {
	z := entities.Zone{Name: "Neutral-high"}
	if got := services.ClassifyZoneTier(z); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestClassifyZoneTier_NameCoreFallback(t *testing.T) {
	z := entities.Zone{Name: "Neutral-core"}
	if got := services.ClassifyZoneTier(z); got != 3 {
		t.Errorf("got %d, want 3", got)
	}
}

func TestClassifyZoneTier_UnknownDefaultsToBronze(t *testing.T) {
	z := entities.Zone{Name: "Neutral-Z"}
	if got := services.ClassifyZoneTier(z); got != 1 {
		t.Errorf("got %d, want 1", got)
	}
}
