package services_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// ── helpers ──────────────────────────────────────────────────────────

func defaultSettings() *models.GeneratorSettings {
	return models.NewGeneratorSettings()
}

func settingsWithTopology(topo models.MapTopology, players, neutrals int) *models.GeneratorSettings {
	s := defaultSettings()
	s.Topology = topo
	s.PlayerCount = players
	s.ZoneCfg.NeutralZoneCount = neutrals
	return s
}

// ── Generate: basic contract ─────────────────────────────────────────

func TestGenerate_EmptyName_ReturnsError(t *testing.T) {
	s := defaultSettings()
	s.TemplateName = ""
	_, err := services.Generate(s)
	if err == nil {
		t.Fatal("expected error for empty template name")
	}
}

func TestGenerate_DefaultSettings_Succeeds(t *testing.T) {
	s := defaultSettings()
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Name != s.TemplateName {
		t.Errorf("name = %q, want %q", tmpl.Name, s.TemplateName)
	}
	if tmpl.SizeX != s.MapSize || tmpl.SizeZ != s.MapSize {
		t.Errorf("size = %dx%d, want %dx%d", tmpl.SizeX, tmpl.SizeZ, s.MapSize, s.MapSize)
	}
	if tmpl.GameMode != "Classic" {
		t.Errorf("gameMode = %q, want Classic", tmpl.GameMode)
	}
	if len(tmpl.Variants) != 1 {
		t.Fatalf("len(variants) = %d, want 1", len(tmpl.Variants))
	}
}

func TestGenerate_MapSizePreserved(t *testing.T) {
	for _, size := range []int{96, 128, 160, 192, 224} {
		s := defaultSettings()
		s.MapSize = size
		tmpl, err := services.Generate(s)
		if err != nil {
			t.Fatalf("size %d: %v", size, err)
		}
		if tmpl.SizeX != size || tmpl.SizeZ != size {
			t.Errorf("size %d: got %dx%d", size, tmpl.SizeX, tmpl.SizeZ)
		}
	}
}

// ── Generate: zone counts ────────────────────────────────────────────

func TestGenerate_ZoneCount_MatchesPlayersPlusNeutrals(t *testing.T) {
	cases := []struct {
		name     string
		topo     models.MapTopology
		players  int
		neutrals int
		wantMin  int // minimum expected zones
	}{
		{"Ring 2p 2n", models.TopologyDefault, 2, 2, 4},
		{"Ring 4p 4n", models.TopologyDefault, 4, 4, 8},
		{"Chain 2p 1n", models.TopologyChain, 2, 1, 3},
		{"Hub 4p 0n", models.TopologyHubAndSpoke, 4, 0, 5}, // 4 player + 1 hub
		{"SharedWeb 2p 1n", models.TopologySharedWeb, 2, 1, 3},
		{"Random 3p 3n", models.TopologyRandom, 3, 3, 6},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s := settingsWithTopology(tc.topo, tc.players, tc.neutrals)
			tmpl, err := services.Generate(s)
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			got := len(tmpl.Variants[0].Zones)
			if got < tc.wantMin {
				t.Errorf("zones = %d, want >= %d", got, tc.wantMin)
			}
		})
	}
}

func TestGenerate_PlayerZoneNames(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 3, 0)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	spawns := 0
	for _, z := range tmpl.Variants[0].Zones {
		if strings.HasPrefix(z.Name, "Spawn-") {
			spawns++
		}
	}
	if spawns != 3 {
		t.Errorf("spawn zones = %d, want 3", spawns)
	}
}

func TestGenerate_NeutralZoneNames(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 3)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	neutrals := 0
	for _, z := range tmpl.Variants[0].Zones {
		if strings.HasPrefix(z.Name, "Neutral-") {
			neutrals++
		}
	}
	if neutrals != 3 {
		t.Errorf("neutral zones = %d, want 3", neutrals)
	}
}

// ── Generate: connections ────────────────────────────────────────────

func TestGenerate_RingTopology_HasConnections(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 4, 4)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	conns := tmpl.Variants[0].Connections
	if len(conns) == 0 {
		t.Fatal("expected connections in ring topology")
	}
	// Ring with 8 zones should have 8 connections
	if len(conns) != 8 {
		t.Errorf("connections = %d, want 8", len(conns))
	}
	for _, c := range conns {
		if c.ConnectionType != "Direct" && c.ConnectionType != "Portal" && c.ConnectionType != "Proximity" {
			t.Errorf("unexpected connection type %q", c.ConnectionType)
		}
	}
}

func TestGenerate_ChainTopology_HasFewerConnections(t *testing.T) {
	s := settingsWithTopology(models.TopologyChain, 3, 2)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	conns := tmpl.Variants[0].Connections
	// Chain with 5 zones: 4 connections (N-1)
	if len(conns) != 4 {
		t.Errorf("chain connections = %d, want 4", len(conns))
	}
}

func TestGenerate_HubTopology_HasHubZone(t *testing.T) {
	s := settingsWithTopology(models.TopologyHubAndSpoke, 4, 2)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, z := range tmpl.Variants[0].Zones {
		if z.Name == "Hub" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected a zone named 'Hub' in HubAndSpoke topology")
	}
}

func TestGenerate_RandomPortals_AddsPortalConnections(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 4, 4)
	s.RandomPortals = true
	s.MaxPortalConnections = 4
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	portals := 0
	for _, c := range tmpl.Variants[0].Connections {
		if c.ConnectionType == "Portal" {
			portals++
		}
	}
	if portals == 0 {
		t.Error("expected portal connections when RandomPortals=true")
	}
}

// ── Generate: roads ──────────────────────────────────────────────────

func TestGenerate_RoadsDisabled_NoRoads(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 2)
	s.GenerateRoads = false
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, z := range tmpl.Variants[0].Zones {
		if len(z.Roads) > 0 {
			t.Errorf("zone %q has %d roads but GenerateRoads=false", z.Name, len(z.Roads))
		}
	}
}

func TestGenerate_RoadsEnabled_HasRoads(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 2)
	s.GenerateRoads = true
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	anyRoads := false
	for _, z := range tmpl.Variants[0].Zones {
		if len(z.Roads) > 0 {
			anyRoads = true
			break
		}
	}
	if !anyRoads {
		t.Error("expected some roads when GenerateRoads=true")
	}
}

// ── Generate: game rules ─────────────────────────────────────────────

func TestGenerate_GameRules_HeroSettings(t *testing.T) {
	s := defaultSettings()
	s.HeroSettings = &models.HeroSettings{HeroCountMin: 6, HeroCountMax: 10, HeroCountIncrement: 2}
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	gr := tmpl.GameRules
	if gr.HeroCountMin != 4 { // 6 - 2 = 4
		t.Errorf("heroCountMin = %d, want 4", gr.HeroCountMin)
	}
	if gr.HeroCountMax != 10 {
		t.Errorf("heroCountMax = %d, want 10", gr.HeroCountMax)
	}
	if gr.HeroCountIncrement != 2 {
		t.Errorf("heroCountIncrement = %d, want 2", gr.HeroCountIncrement)
	}
}

func TestGenerate_WinConditions_CityHold(t *testing.T) {
	s := defaultSettings()
	s.GameEndConditions = &models.GameEndConditions{
		VictoryCondition: "win_condition_5",
		CityHold:         true,
		CityHoldDays:     7,
		LostStartCityDay: 3,
	}
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	wc := tmpl.GameRules.WinConditions
	if !wc.CityHold {
		t.Error("expected cityHold=true")
	}
	if wc.CityHoldDays != 7 {
		t.Errorf("cityHoldDays = %d, want 7", wc.CityHoldDays)
	}
}

func TestGenerate_WinConditions_GladiatorArena(t *testing.T) {
	s := defaultSettings()
	s.GladiatorArenaRules = &models.GladiatorArenaRules{Enabled: true, DaysDelayStart: 15, CountDay: 5}
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	wc := tmpl.GameRules.WinConditions
	if !wc.GladiatorArena {
		t.Error("expected gladiatorArena=true")
	}
	if wc.GladiatorArenaDaysDelayStart != 15 {
		t.Errorf("gladiatorArenaDaysDelayStart = %d, want 15", wc.GladiatorArenaDaysDelayStart)
	}
	if wc.GladiatorArenaCountDay != 5 {
		t.Errorf("gladiatorArenaCountDay = %d, want 5", wc.GladiatorArenaCountDay)
	}
	if wc.ChampionSelectRule != "StartHero" {
		t.Errorf("championSelectRule = %q, want StartHero", wc.ChampionSelectRule)
	}
}

func TestGenerate_WinConditions_Tournament(t *testing.T) {
	s := defaultSettings()
	s.TournamentRules = &models.TournamentRules{Enabled: true, FirstTournamentDay: 14, Interval: 7, PointsToWin: 3}
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	wc := tmpl.GameRules.WinConditions
	if !wc.Tournament {
		t.Error("expected tournament=true")
	}
	if wc.TournamentPointsToWin != 3 {
		t.Errorf("tournamentPointsToWin = %d, want 3", wc.TournamentPointsToWin)
	}
	// pointsToWin=3 → roundCount=5 → 5 announce days and 5 battle offsets
	if len(wc.TournamentAnnounceDays) != 5 {
		t.Errorf("len(announceDays) = %d, want 5", len(wc.TournamentAnnounceDays))
	}
	if len(wc.TournamentDays) != 5 {
		t.Errorf("len(tournamentDays) = %d, want 5", len(wc.TournamentDays))
	}
}

func TestGenerate_FactionLawsExpModifier(t *testing.T) {
	s := defaultSettings()
	s.FactionLawsExpPercent = 150
	s.AstrologyExpPercent = 75
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.GameRules.FactionLawsExpModifier != 1.5 {
		t.Errorf("factionLawsExp = %v, want 1.5", tmpl.GameRules.FactionLawsExpModifier)
	}
	if tmpl.GameRules.AstrologyExpModifier != 0.75 {
		t.Errorf("astrologyExp = %v, want 0.75", tmpl.GameRules.AstrologyExpModifier)
	}
}

// ── Generate: mandatory content ──────────────────────────────────────

func TestGenerate_MandatoryContent_PerZone(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 3)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	// 2 player + 3 neutral = 5 mandatory content groups
	if len(tmpl.MandatoryContent) != 5 {
		t.Errorf("mandatoryContent groups = %d, want 5", len(tmpl.MandatoryContent))
	}
	playerGroups := 0
	neutralGroups := 0
	for _, mc := range tmpl.MandatoryContent {
		if strings.HasPrefix(mc.Name, "mandatory_content_side_") {
			playerGroups++
		}
		if strings.HasPrefix(mc.Name, "mandatory_content_neutral_") {
			neutralGroups++
		}
	}
	if playerGroups != 2 {
		t.Errorf("player mandatory groups = %d, want 2", playerGroups)
	}
	if neutralGroups != 3 {
		t.Errorf("neutral mandatory groups = %d, want 3", neutralGroups)
	}
}

// ── Generate: zone layouts ───────────────────────────────────────────

func TestGenerate_ZoneLayouts_Has4(t *testing.T) {
	s := defaultSettings()
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	if len(tmpl.ZoneLayouts) != 4 {
		t.Errorf("zoneLayouts = %d, want 4", len(tmpl.ZoneLayouts))
	}
	names := map[string]bool{}
	for _, zl := range tmpl.ZoneLayouts {
		names[zl.Name] = true
	}
	for _, want := range []string{"zone_layout_spawns", "zone_layout_sides", "zone_layout_treasure_zone", "zone_layout_center"} {
		if !names[want] {
			t.Errorf("missing zone layout %q", want)
		}
	}
}

// ── Generate: content count limits ───────────────────────────────────

func TestGenerate_ContentCountLimits(t *testing.T) {
	s := defaultSettings()
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	// 2 base + 15 pair combos (C(5,2) for sides 1-5 against 2-6) = 17
	if len(tmpl.ContentCountLimits) != 17 {
		t.Errorf("contentCountLimits = %d, want 17", len(tmpl.ContentCountLimits))
	}
}

// ── Generate: orientation ────────────────────────────────────────────

func TestGenerate_Orientation_AngleStep(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 4, 4)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	orient := tmpl.Variants[0].Orientation
	// 8 zones → 360/8 = 45
	if orient.RandomAngleStep != 45 {
		t.Errorf("randomAngleStep = %d, want 45", orient.RandomAngleStep)
	}
	if orient.BaseAngleMin != 45 || orient.BaseAngleMax != 45 {
		t.Errorf("baseAngle = %d..%d, want 45..45", orient.BaseAngleMin, orient.BaseAngleMax)
	}
}

// ── Generate: all topologies succeed ─────────────────────────────────

func TestGenerate_AllTopologies_DoNotError(t *testing.T) {
	topos := []models.MapTopology{
		models.TopologyDefault,
		models.TopologyChain,
		models.TopologyHubAndSpoke,
		models.TopologySharedWeb,
		models.TopologyRandom,
	}
	for _, topo := range topos {
		for _, players := range []int{2, 3, 4, 8} {
			for _, neutrals := range []int{0, 1, 4} {
				name := string(topo) + "_" + string(rune('0'+players)) + "p_" + string(rune('0'+neutrals)) + "n"
				t.Run(name, func(t *testing.T) {
					s := settingsWithTopology(topo, players, neutrals)
					tmpl, err := services.Generate(s)
					if err != nil {
						t.Fatalf("error: %v", err)
					}
					if len(tmpl.Variants) == 0 {
						t.Fatal("no variants produced")
					}
					if len(tmpl.Variants[0].Zones) == 0 {
						t.Fatal("no zones produced")
					}
				})
			}
		}
	}
}

// ── Generate: spawn zone content ─────────────────────────────────────

func TestGenerate_SpawnZone_HasMainObjectSpawn(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 0)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, z := range tmpl.Variants[0].Zones {
		if !strings.HasPrefix(z.Name, "Spawn-") {
			continue
		}
		if len(z.MainObjects) == 0 {
			t.Errorf("zone %q has no main objects", z.Name)
			continue
		}
		if z.MainObjects[0].Type != "Spawn" {
			t.Errorf("zone %q: first main object type = %q, want Spawn", z.Name, z.MainObjects[0].Type)
		}
		if z.MainObjects[0].Spawn == "" {
			t.Errorf("zone %q: spawn field is empty", z.Name)
		}
	}
}

func TestGenerate_SpawnZone_MultipleCastles(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 0)
	s.ZoneCfg.PlayerZoneCastles = 3
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, z := range tmpl.Variants[0].Zones {
		if !strings.HasPrefix(z.Name, "Spawn-") {
			continue
		}
		if len(z.MainObjects) != 3 {
			t.Errorf("zone %q: mainObjects = %d, want 3", z.Name, len(z.MainObjects))
		}
	}
}

// ── Generate: advanced mode neutral zones ────────────────────────────

func TestGenerate_AdvancedMode_MixedNeutralTiers(t *testing.T) {
	s := defaultSettings()
	s.Topology = models.TopologyDefault
	s.PlayerCount = 2
	s.ZoneCfg.Advanced.Enabled = true
	s.ZoneCfg.Advanced.NeutralLowNoCastleCount = 1
	s.ZoneCfg.Advanced.NeutralMediumCastleCount = 1
	s.ZoneCfg.Advanced.NeutralHighCastleCount = 1
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	zones := tmpl.Variants[0].Zones
	neutrals := 0
	for _, z := range zones {
		if strings.HasPrefix(z.Name, "Neutral-") {
			neutrals++
		}
	}
	if neutrals != 3 {
		t.Errorf("advanced mode neutral zones = %d, want 3", neutrals)
	}
}

// ── Generate: city hold ──────────────────────────────────────────────

func TestGenerate_CityHold_NeutralHasHoldCityFlag(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 2)
	s.GameEndConditions = &models.GameEndConditions{
		VictoryCondition: "win_condition_5",
		CityHold:         true,
		CityHoldDays:     6,
		LostStartCityDay: 3,
	}
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	holdCities := 0
	for _, z := range tmpl.Variants[0].Zones {
		for _, mo := range z.MainObjects {
			if mo.HoldCityWinCon {
				holdCities++
			}
		}
	}
	if holdCities == 0 {
		t.Error("expected at least one HoldCityWinCon=true main object")
	}
}

// ── Generate: description ────────────────────────────────────────────

func TestGenerate_Description_ContainsTopology(t *testing.T) {
	cases := []struct {
		topo models.MapTopology
		want string
	}{
		{models.TopologyDefault, "Ring"},
		{models.TopologyChain, "Chain"},
		{models.TopologyHubAndSpoke, "Hub"},
		{models.TopologySharedWeb, "Shared Web"},
		{models.TopologyRandom, "Random"},
	}
	for _, tc := range cases {
		s := settingsWithTopology(tc.topo, 2, 1)
		tmpl, err := services.Generate(s)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(tmpl.Description, tc.want) {
			t.Errorf("topology %s: description %q missing %q", tc.topo, tmpl.Description, tc.want)
		}
	}
}

func TestGenerate_Description_ContainsOptions(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 2)
	s.NoDirectPlayerConnections = true
	s.RandomPortals = true
	s.GenerateRoads = false
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"isolated player starts", "random portals", "roads disabled"} {
		if !strings.Contains(tmpl.Description, want) {
			t.Errorf("description %q missing option %q", tmpl.Description, want)
		}
	}
}

// ── Generate: isolation mode ─────────────────────────────────────────

func TestGenerate_Isolation_NoDirectPlayerConnections(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 2)
	s.NoDirectPlayerConnections = true
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range tmpl.Variants[0].Connections {
		if c.ConnectionType == "Direct" {
			isFromPlayer := strings.HasPrefix(c.From, "Spawn-")
			isToPlayer := strings.HasPrefix(c.To, "Spawn-")
			if isFromPlayer && isToPlayer {
				t.Errorf("direct player-player connection found: %s → %s", c.From, c.To)
			}
		}
	}
}

// ── ZoneContentManager tests ─────────────────────────────────────────

func TestBuildPlayerZoneMandatoryContent_HasContent(t *testing.T) {
	s := defaultSettings()
	content := services.BuildPlayerZoneMandatoryContent(s)
	if len(content) == 0 {
		t.Fatal("expected non-empty player zone mandatory content")
	}
}

func TestBuildPlayerZoneMandatoryContent_WithFoothold(t *testing.T) {
	s := defaultSettings()
	s.SpawnRemoteFootholds = true
	content := services.BuildPlayerZoneMandatoryContent(s)
	found := false
	for _, c := range content {
		if c.Name == "name_remote_foothold_1" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected remote foothold when SpawnRemoteFootholds=true")
	}
}

func TestBuildPlayerZoneMandatoryContent_WithoutFoothold(t *testing.T) {
	s := defaultSettings()
	s.SpawnRemoteFootholds = false
	content := services.BuildPlayerZoneMandatoryContent(s)
	for _, c := range content {
		if c.Name == "name_remote_foothold_1" {
			t.Error("unexpected remote foothold when SpawnRemoteFootholds=false")
		}
	}
}

func TestBuildPlayerZoneMandatoryContent_AppendsUserContent(t *testing.T) {
	s := defaultSettings()
	s.PlayerZoneMandatoryContent = []template.MandatoryContentItem{
		{SID: "custom_item_1"},
		{SID: "custom_item_2"},
	}
	content := services.BuildPlayerZoneMandatoryContent(s)
	found := 0
	for _, c := range content {
		if c.SID == "custom_item_1" || c.SID == "custom_item_2" {
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 user items, found %d", found)
	}
}

func TestBuildLowNeutralMandatoryContent_HasContent(t *testing.T) {
	content := services.BuildLowNeutralMandatoryContent(1, true)
	if len(content) == 0 {
		t.Fatal("expected non-empty low neutral content")
	}
}

func TestBuildLowNeutralMandatoryContent_FootholdToggle(t *testing.T) {
	withFoot := services.BuildLowNeutralMandatoryContent(1, true)
	withoutFoot := services.BuildLowNeutralMandatoryContent(1, false)
	if len(withFoot) <= len(withoutFoot) {
		t.Error("foothold=true should produce more items")
	}
}

func TestBuildMediumNeutralMandatoryContent_HasMines(t *testing.T) {
	content := services.BuildMediumNeutralMandatoryContent(1, false)
	hasMine := false
	for _, c := range content {
		if c.IsMine {
			hasMine = true
			break
		}
	}
	if !hasMine {
		t.Error("medium neutral content should include mines")
	}
}

func TestBuildHighNeutralMandatoryContent_HasMoreThanMedium(t *testing.T) {
	medium := services.BuildMediumNeutralMandatoryContent(1, false)
	high := services.BuildHighNeutralMandatoryContent(1, false)
	if len(high) <= len(medium) {
		t.Errorf("high (%d items) should have more items than medium (%d)", len(high), len(medium))
	}
}

func TestBuildAllContentCountLimits_Count(t *testing.T) {
	s := defaultSettings()
	limits := services.BuildAllContentCountLimits(s)
	// 2 base ("content_limits_side", "content_limits_side_0_0") + C(5,2)=15 pair combos = 17
	if len(limits) != 17 {
		t.Errorf("content count limits = %d, want 17", len(limits))
	}
}

func TestBuildAllContentCountLimits_UserContentLiftsLimit(t *testing.T) {
	s := defaultSettings()
	// Add 5 pandora boxes — should lift the default limit (4)
	for i := 0; i < 5; i++ {
		s.PlayerZoneMandatoryContent = append(s.PlayerZoneMandatoryContent, template.MandatoryContentItem{SID: constants.ContentIds.PandoraBox.Sid})
	}
	limits := services.BuildAllContentCountLimits(s)
	for _, group := range limits {
		for _, l := range group.Limits {
			if strings.EqualFold(l.SID, constants.ContentIds.PandoraBox.Sid) {
				if l.MaxCount < 5 {
					t.Errorf("pandora box limit = %d, want >= 5 after user added 5", l.MaxCount)
				}
				return
			}
		}
	}
	t.Error("pandora box limit not found in content count limits")
}

// ── SharedWeb topology ───────────────────────────────────────────────

func TestGenerate_SharedWeb_ForcesAtLeastOneNeutral(t *testing.T) {
	s := settingsWithTopology(models.TopologySharedWeb, 2, 0)
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	neutrals := 0
	for _, z := range tmpl.Variants[0].Zones {
		if strings.HasPrefix(z.Name, "Neutral-") {
			neutrals++
		}
	}
	if neutrals < 1 {
		t.Error("SharedWeb with 0 neutrals should still create at least 1 neutral zone")
	}
}

// ── Tournament topology ──────────────────────────────────────────────

func TestGenerate_TournamentMode_2Players(t *testing.T) {
	s := settingsWithTopology(models.TopologyDefault, 2, 4)
	s.TournamentRules = &models.TournamentRules{Enabled: true, FirstTournamentDay: 14, Interval: 7, PointsToWin: 2}
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	zones := tmpl.Variants[0].Zones
	spawns := 0
	for _, z := range zones {
		if strings.HasPrefix(z.Name, "Spawn-") {
			spawns++
		}
	}
	if spawns != 2 {
		t.Errorf("tournament 2p: spawn zones = %d, want 2", spawns)
	}
}

// ── Comprehensive topology JSON validity ─────────────────────────────

func TestGenerate_AllZones_HaveRequiredFields(t *testing.T) {
	topos := []models.MapTopology{
		models.TopologyDefault, models.TopologyChain,
		models.TopologyHubAndSpoke, models.TopologySharedWeb,
		models.TopologyRandom,
	}
	for _, topo := range topos {
		s := settingsWithTopology(topo, 3, 3)
		tmpl, err := services.Generate(s)
		if err != nil {
			t.Fatalf("%s: %v", topo, err)
		}
		for _, z := range tmpl.Variants[0].Zones {
			if z.Name == "" {
				t.Errorf("%s: zone with empty name", topo)
			}
			if z.Layout == "" {
				t.Errorf("%s: zone %q has empty layout", topo, z.Name)
			}
			if z.Size <= 0 || z.Size > 2.0 {
				t.Errorf("%s: zone %q size = %f out of range", topo, z.Name, z.Size)
			}
			if z.GuardCutoffValue <= 0 {
				t.Errorf("%s: zone %q guardCutoff = %d", topo, z.Name, z.GuardCutoffValue)
			}
			if len(z.GuardedContentPool) == 0 {
				t.Errorf("%s: zone %q has empty guarded content pool", topo, z.Name)
			}
			if len(z.UnguardedContentPool) == 0 {
				t.Errorf("%s: zone %q has empty unguarded content pool", topo, z.Name)
			}
			if len(z.ResourcesContentPool) == 0 {
				t.Errorf("%s: zone %q has empty resources content pool", topo, z.Name)
			}
		}
	}
}

func TestGenerate_AllConnections_ReferenceValidZones(t *testing.T) {
	topos := []models.MapTopology{
		models.TopologyDefault, models.TopologyChain,
		models.TopologyHubAndSpoke, models.TopologySharedWeb,
		models.TopologyRandom,
	}
	for _, topo := range topos {
		s := settingsWithTopology(topo, 3, 2)
		tmpl, err := services.Generate(s)
		if err != nil {
			t.Fatalf("%s: %v", topo, err)
		}
		zoneNames := map[string]bool{}
		for _, z := range tmpl.Variants[0].Zones {
			zoneNames[z.Name] = true
		}
		for _, c := range tmpl.Variants[0].Connections {
			if !zoneNames[c.From] {
				t.Errorf("%s: connection %q references unknown From zone %q", topo, c.Name, c.From)
			}
			if !zoneNames[c.To] {
				t.Errorf("%s: connection %q references unknown To zone %q", topo, c.Name, c.To)
			}
		}
	}
}
