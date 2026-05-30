package templateGenerator_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

func newCfg() *config.GeneratorConfig { return config.NewGeneratorConfig() }

func cfgWith(topo config.MapTopology, players, neutrals int) *config.GeneratorConfig {
	s := newCfg()
	s.Topology = topo
	s.PlayerCount = players
	s.ZoneConfiguration.NeutralZoneCount = neutrals
	return s
}

// ── Generate: input validation ───────────────────────────────────────

func TestGenerate_EmptyName_ReturnsError(t *testing.T) {
	s := newCfg()
	s.TemplateName = ""
	if _, err := services.Generate(s); err == nil {
		t.Fatal("expected error for empty template name")
	}
}

func TestGenerate_ValidConfig_NoError(t *testing.T) {
	if _, err := services.Generate(newCfg()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// ── Generate: scalar field propagation ───────────────────────────────

func TestGenerate_NamePropagated(t *testing.T) {
	s := newCfg()
	s.TemplateName = "My Template"
	tmpl, _ := services.Generate(s)
	if tmpl.Name != "My Template" {
		t.Errorf("Name = %q", tmpl.Name)
	}
}

func TestGenerate_MapSizePropagated(t *testing.T) {
	s := newCfg()
	s.MapSize = 224
	tmpl, _ := services.Generate(s)
	if tmpl.SizeX != 224 || tmpl.SizeZ != 224 {
		t.Errorf("size = %dx%d", tmpl.SizeX, tmpl.SizeZ)
	}
}

func TestGenerate_GameModeClassic(t *testing.T) {
	s := newCfg()
	s.GameMode = "Classic"
	tmpl, _ := services.Generate(s)
	if tmpl.GameMode != "Classic" {
		t.Errorf("gameMode = %q", tmpl.GameMode)
	}
	if tmpl.GameRules.HeroHireBan {
		t.Error("HeroHireBan should be false in Classic")
	}
}

func TestGenerate_GameModeSingleHero_ForcesHeroCounts(t *testing.T) {
	s := newCfg()
	s.GameMode = "SingleHero"
	s.HeroSettings = config.HeroSettings{HeroCountMin: 4, HeroCountMax: 8, HeroCountIncrement: 2}
	tmpl, _ := services.Generate(s)
	gr := tmpl.GameRules
	if gr.HeroCountMin != 1 || gr.HeroCountMax != 1 || gr.HeroCountIncrement != 1 {
		t.Errorf("SingleHero should force hero counts to 1, got %+v", gr)
	}
	if !gr.HeroHireBan {
		t.Error("HeroHireBan should be true in SingleHero")
	}
}

func TestGenerate_HeroSettingsPropagatedInClassic(t *testing.T) {
	s := newCfg()
	s.HeroSettings = config.HeroSettings{HeroCountMin: 3, HeroCountMax: 7, HeroCountIncrement: 2}
	tmpl, _ := services.Generate(s)
	gr := tmpl.GameRules
	if gr.HeroCountMin != 3 || gr.HeroCountMax != 7 || gr.HeroCountIncrement != 2 {
		t.Errorf("hero settings mismatch: %+v", gr)
	}
}

func TestGenerate_FactionLawsExpModifierClampedAndScaled(t *testing.T) {
	s := newCfg()
	s.FactionLawsExpPercent = 150
	tmpl, _ := services.Generate(s)
	if tmpl.GameRules.FactionLawsExpModifier != 1.5 {
		t.Errorf("got %v, want 1.5", tmpl.GameRules.FactionLawsExpModifier)
	}
}

func TestGenerate_AstrologyExpModifierScaled(t *testing.T) {
	s := newCfg()
	s.AstrologyExpPercent = 50
	tmpl, _ := services.Generate(s)
	if tmpl.GameRules.AstrologyExpModifier != 0.5 {
		t.Errorf("got %v, want 0.5", tmpl.GameRules.AstrologyExpModifier)
	}
}

// ── Generate: zone counts per topology ───────────────────────────────

func countZones(t *testing.T, s *config.GeneratorConfig) (spawn, neutral int) {
	t.Helper()
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	for _, z := range tmpl.Variants[0].Zones {
		if strings.HasPrefix(z.Name, "Spawn-") {
			spawn++
		}
		if strings.HasPrefix(z.Name, "Neutral-") {
			neutral++
		}
	}
	return
}

func TestGenerate_DefaultTopology_ZoneCounts(t *testing.T) {
	sp, ne := countZones(t, cfgWith(config.TopologyDefault, 3, 2))
	if sp != 3 || ne != 2 {
		t.Errorf("spawn=%d neutral=%d, want 3/2", sp, ne)
	}
}

func TestGenerate_ChainTopology_HasNMinusOneConnections(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyChain, 3, 2))
	if got := len(tmpl.Variants[0].Connections); got != 4 {
		t.Errorf("chain conns = %d, want 4", got)
	}
}

func TestGenerate_HubAndSpokeTopology_HasHubZone(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyHubAndSpoke, 3, 2))
	found := false
	for _, z := range tmpl.Variants[0].Zones {
		if z.Name == "Hub" {
			found = true
		}
	}
	if !found {
		t.Error("expected Hub zone")
	}
}

func TestGenerate_SharedWebTopology_ForcesNeutralWhenZero(t *testing.T) {
	_, ne := countZones(t, cfgWith(config.TopologySharedWeb, 2, 0))
	if ne < 1 {
		t.Error("expected SharedWeb to inject at least one neutral")
	}
}

func TestGenerate_RandomTopology_SetsGeneratorPosition(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyRandom, 3, 2))
	for _, z := range tmpl.Variants[0].Zones {
		if z.GeneratorPosition == nil {
			t.Errorf("zone %q missing GeneratorPosition", z.Name)
		}
	}
}

func TestGenerate_BalancedTopology_SetsGeneratorRing(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyBalanced, 3, 2))
	for _, z := range tmpl.Variants[0].Zones {
		if z.GeneratorRing == nil {
			t.Errorf("zone %q missing GeneratorRing", z.Name)
		}
	}
}

// ── Generate: connection-type behaviour ──────────────────────────────

func TestGenerate_RandomPortals_AddsPortalConnections(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 4, 4)
	s.RandomPortals = true
	s.MaxPortalConnections = 4
	tmpl, _ := services.Generate(s)
	portals := 0
	for _, c := range tmpl.Variants[0].Connections {
		if c.ConnectionType == "Portal" {
			portals++
		}
	}
	if portals == 0 {
		t.Error("expected portal connections")
	}
}

func TestGenerate_RandomPortalsDisabled_NoPortals(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 4, 4)
	s.RandomPortals = false
	tmpl, _ := services.Generate(s)
	for _, c := range tmpl.Variants[0].Connections {
		if c.ConnectionType == "Portal" {
			t.Error("portal connection found when RandomPortals=false")
		}
	}
}

func TestGenerate_NoDirectPlayerConnections_Enforced(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 2)
	s.NoDirectPlayerConnections = true
	tmpl, _ := services.Generate(s)
	for _, c := range tmpl.Variants[0].Connections {
		if c.ConnectionType == "Direct" &&
			strings.HasPrefix(c.From, "Spawn-") && strings.HasPrefix(c.To, "Spawn-") {
			t.Errorf("direct player-player connection: %s→%s", c.From, c.To)
		}
	}
}

// ── Generate: roads ──────────────────────────────────────────────────

func TestGenerate_RoadsEnabled_ProducesRoads(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 2)
	s.GenerateRoads = true
	tmpl, _ := services.Generate(s)
	any := false
	for _, z := range tmpl.Variants[0].Zones {
		if len(z.Roads) > 0 {
			any = true
			break
		}
	}
	if !any {
		t.Error("expected some roads")
	}
}

func TestGenerate_RoadsDisabled_NoRoads(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 2)
	s.GenerateRoads = false
	tmpl, _ := services.Generate(s)
	for _, z := range tmpl.Variants[0].Zones {
		if len(z.Roads) > 0 {
			t.Errorf("zone %q has roads", z.Name)
		}
	}
}

// ── Generate: castle factions ────────────────────────────────────────

func TestGenerate_MatchPlayerCastleFactions_True(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 0)
	s.ZoneConfiguration.PlayerZoneCastles = 2
	s.MatchPlayerCastleFactions = true
	tmpl, _ := services.Generate(s)
	for _, z := range tmpl.Variants[0].Zones {
		if !strings.HasPrefix(z.Name, "Spawn-") || len(z.MainObjects) < 2 {
			continue
		}
		if z.MainObjects[1].Faction == nil || z.MainObjects[1].Faction.Type != "Match" {
			t.Errorf("expected Match faction, got %+v", z.MainObjects[1].Faction)
		}
	}
}

func TestGenerate_MatchPlayerCastleFactions_False(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 0)
	s.ZoneConfiguration.PlayerZoneCastles = 2
	s.MatchPlayerCastleFactions = false
	tmpl, _ := services.Generate(s)
	for _, z := range tmpl.Variants[0].Zones {
		if !strings.HasPrefix(z.Name, "Spawn-") || len(z.MainObjects) < 2 {
			continue
		}
		if z.MainObjects[1].Faction == nil || z.MainObjects[1].Faction.Type != "Random" {
			t.Errorf("expected Random faction, got %+v", z.MainObjects[1].Faction)
		}
	}
}

// ── Generate: city hold / lost city ──────────────────────────────────

func TestGenerate_CityHoldExplicit_SetsHoldCity(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 2)
	s.GameEndConditions = &config.GameEndConditions{
		VictoryCondition: "win_condition_1",
		CityHold:         true, CityHoldDays: 5, LostStartCityDay: 3,
	}
	tmpl, _ := services.Generate(s)
	if !tmpl.GameRules.WinConditions.CityHold {
		t.Error("expected CityHold true")
	}
	hits := 0
	for _, z := range tmpl.Variants[0].Zones {
		for _, mo := range z.MainObjects {
			if mo.HoldCityWinCon {
				hits++
			}
		}
	}
	if hits == 0 {
		t.Error("expected a HoldCityWinCon main object")
	}
}

func TestGenerate_CityHoldFromVictoryCondition5(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 2)
	s.GameEndConditions = &config.GameEndConditions{VictoryCondition: "win_condition_5", LostStartCityDay: 3, CityHoldDays: 6}
	tmpl, _ := services.Generate(s)
	if !tmpl.GameRules.WinConditions.CityHold {
		t.Error("win_condition_5 should set CityHold")
	}
}

func TestGenerate_HubAndSpokeCityHold_HubIsHoldCity(t *testing.T) {
	s := cfgWith(config.TopologyHubAndSpoke, 2, 0)
	s.GameEndConditions = &config.GameEndConditions{VictoryCondition: "win_condition_5", LostStartCityDay: 3, CityHoldDays: 6, CityHold: true}
	tmpl, _ := services.Generate(s)
	hubHold := false
	for _, z := range tmpl.Variants[0].Zones {
		if z.Name != "Hub" {
			continue
		}
		for _, mo := range z.MainObjects {
			if mo.HoldCityWinCon {
				hubHold = true
			}
		}
	}
	if !hubHold {
		t.Error("expected Hub to be the hold city")
	}
}

func TestGenerate_LostStartCityFromVictoryCondition3(t *testing.T) {
	s := newCfg()
	s.GameEndConditions = &config.GameEndConditions{VictoryCondition: "win_condition_3", LostStartCityDay: 5, CityHoldDays: 6}
	tmpl, _ := services.Generate(s)
	if !tmpl.GameRules.WinConditions.LostStartCity {
		t.Error("win_condition_3 should set LostStartCity")
	}
}

// ── Generate: gladiator arena ────────────────────────────────────────

func TestGenerate_GladiatorArena_Enabled(t *testing.T) {
	s := newCfg()
	s.GladiatorArenaRules = &config.GladiatorArenaRules{Enabled: true, DaysDelayStart: 10, CountDay: 4}
	tmpl, _ := services.Generate(s)
	wc := tmpl.GameRules.WinConditions
	if !wc.GladiatorArena || wc.GladiatorArenaDaysDelayStart != 10 || wc.GladiatorArenaCountDay != 4 {
		t.Errorf("gladiator config not propagated: %+v", wc)
	}
	if wc.ChampionSelectRule != "StartHero" {
		t.Errorf("ChampionSelectRule = %q", wc.ChampionSelectRule)
	}
}

func TestGenerate_GladiatorArenaFromVictoryCondition4(t *testing.T) {
	s := newCfg()
	s.GameEndConditions = &config.GameEndConditions{VictoryCondition: "win_condition_4", LostStartCityDay: 3, CityHoldDays: 6}
	tmpl, _ := services.Generate(s)
	if !tmpl.GameRules.WinConditions.GladiatorArena {
		t.Error("win_condition_4 should enable GladiatorArena")
	}
}

// ── Generate: tournament ─────────────────────────────────────────────

func TestGenerate_TournamentEnabled_FillsRoundSchedule(t *testing.T) {
	s := newCfg()
	s.TournamentRules = &config.TournamentRules{Enabled: true, FirstTournamentDay: 10, Interval: 5, PointsToWin: 3, SaveArmy: true}
	tmpl, _ := services.Generate(s)
	wc := tmpl.GameRules.WinConditions
	if !wc.Tournament {
		t.Fatal("expected Tournament=true")
	}
	if wc.TournamentPointsToWin != 3 {
		t.Errorf("pointsToWin = %d", wc.TournamentPointsToWin)
	}
	// pointsToWin=3 → roundCount=5
	if len(wc.TournamentAnnounceDays) != 5 || len(wc.TournamentDays) != 5 {
		t.Errorf("expected 5 announce/battle slots, got %d/%d", len(wc.TournamentAnnounceDays), len(wc.TournamentDays))
	}
}

func TestGenerate_TournamentFromVictoryCondition6(t *testing.T) {
	s := newCfg()
	s.GameEndConditions = &config.GameEndConditions{VictoryCondition: "win_condition_6", LostStartCityDay: 3, CityHoldDays: 6}
	tmpl, _ := services.Generate(s)
	if !tmpl.GameRules.WinConditions.Tournament {
		t.Error("win_condition_6 should enable Tournament")
	}
}

func TestGenerate_Tournament2Players_BuildsSplitClusters(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 4)
	s.TournamentRules = &config.TournamentRules{Enabled: true, FirstTournamentDay: 14, Interval: 7, PointsToWin: 2}
	tmpl, _ := services.Generate(s)
	hasRingGuard := false
	for _, c := range tmpl.Variants[0].Connections {
		if strings.HasPrefix(c.GuardMatchGroup, "tourney_ring_guard_") {
			hasRingGuard = true
			break
		}
	}
	if !hasRingGuard {
		t.Error("expected tourney_ring_guard groups in 2p tournament default")
	}
}

func TestGenerate_TournamentHubAndSpoke_CreatesPerPlayerHubs(t *testing.T) {
	s := cfgWith(config.TopologyHubAndSpoke, 2, 2)
	s.TournamentRules = &config.TournamentRules{Enabled: true, FirstTournamentDay: 14, Interval: 7, PointsToWin: 2}
	tmpl, _ := services.Generate(s)
	hubCount := 0
	for _, z := range tmpl.Variants[0].Zones {
		if strings.HasPrefix(z.Name, "Hub-") {
			hubCount++
		}
	}
	if hubCount != 2 {
		t.Errorf("expected 2 per-player hubs, got %d", hubCount)
	}
}

func TestGenerate_TournamentChainTopology(t *testing.T) {
	s := cfgWith(config.TopologyChain, 2, 2)
	s.TournamentRules = &config.TournamentRules{Enabled: true, FirstTournamentDay: 14, Interval: 7, PointsToWin: 2}
	tmpl, _ := services.Generate(s)
	hits := 0
	for _, c := range tmpl.Variants[0].Connections {
		if strings.HasPrefix(c.GuardMatchGroup, "tourney_guard_") {
			hits++
		}
	}
	if hits == 0 {
		t.Error("expected tourney_guard_* in chain tournament")
	}
}

func TestGenerate_TournamentBalancedTopology(t *testing.T) {
	s := cfgWith(config.TopologyBalanced, 2, 4)
	s.TournamentRules = &config.TournamentRules{Enabled: true, FirstTournamentDay: 14, Interval: 7, PointsToWin: 2}
	tmpl, _ := services.Generate(s)
	hits := 0
	for _, c := range tmpl.Variants[0].Connections {
		if strings.HasPrefix(c.GuardMatchGroup, "tourney_bal_guard_") {
			hits++
		}
	}
	if hits == 0 {
		t.Error("expected tourney_bal_guard_* in balanced tournament")
	}
}

// ── Generate: advanced neutral mix ───────────────────────────────────

func TestGenerate_AdvancedMode_MixedNeutralCounts(t *testing.T) {
	s := newCfg()
	s.Topology = config.TopologyDefault
	s.PlayerCount = 2
	s.ZoneConfiguration.Advanced.Enabled = true
	s.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = 1
	s.ZoneConfiguration.Advanced.NeutralMediumCastleCount = 1
	s.ZoneConfiguration.Advanced.NeutralHighCastleCount = 1
	_, ne := countZones(t, s)
	if ne != 3 {
		t.Errorf("expected 3 neutrals, got %d", ne)
	}
}

func TestGenerate_AdvancedMode_GuardRandomizationClamped(t *testing.T) {
	s := newCfg()
	s.ZoneConfiguration.Advanced.Enabled = true
	s.ZoneConfiguration.Advanced.GuardRandomization = 5.0 // way over 0.5
	tmpl, _ := services.Generate(s)
	for _, z := range tmpl.Variants[0].Zones {
		if z.GuardRandomization > 0.5 {
			t.Errorf("zone %q: GuardRandomization = %f, want <= 0.5", z.Name, z.GuardRandomization)
		}
	}
}

// ── Generate: structural template fields ─────────────────────────────

func TestGenerate_AlwaysHasOneVariant(t *testing.T) {
	tmpl, _ := services.Generate(newCfg())
	if len(tmpl.Variants) != 1 {
		t.Errorf("expected 1 variant, got %d", len(tmpl.Variants))
	}
}

func TestGenerate_ZoneLayouts_Has4Named(t *testing.T) {
	tmpl, _ := services.Generate(newCfg())
	if len(tmpl.ZoneLayouts) != 4 {
		t.Errorf("expected 4 zone layouts, got %d", len(tmpl.ZoneLayouts))
	}
}

func TestGenerate_ContentCountLimits_HasGroups(t *testing.T) {
	tmpl, _ := services.Generate(newCfg())
	if len(tmpl.ContentCountLimits) == 0 {
		t.Error("expected content count limits")
	}
}

func TestGenerate_Description_ContainsTopologyName(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyChain, 2, 2))
	if !strings.Contains(tmpl.Description, "Chain") {
		t.Errorf("description %q missing Chain", tmpl.Description)
	}
}

func TestGenerate_Description_OptionsAppended(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 2)
	s.NoDirectPlayerConnections = true
	s.RandomPortals = true
	s.SpawnRemoteFootholds = false
	s.GenerateRoads = false
	tmpl, _ := services.Generate(s)
	for _, want := range []string{"isolated player starts", "random portals", "no remote footholds", "roads disabled"} {
		if !strings.Contains(tmpl.Description, want) {
			t.Errorf("description %q missing %q", tmpl.Description, want)
		}
	}
}

func TestGenerate_DisplayWinConditionPropagated(t *testing.T) {
	s := newCfg()
	s.GameEndConditions = &config.GameEndConditions{VictoryCondition: "win_condition_2", LostStartCityDay: 3, CityHoldDays: 6}
	tmpl, _ := services.Generate(s)
	if tmpl.DisplayWinCondition != "win_condition_2" {
		t.Errorf("got %q", tmpl.DisplayWinCondition)
	}
}

func TestGenerate_NilGameEndConditions_UsesDefaults(t *testing.T) {
	s := newCfg()
	s.GameEndConditions = nil
	tmpl, err := services.Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	if tmpl.DisplayWinCondition != "win_condition_1" {
		t.Errorf("expected default win_condition_1, got %q", tmpl.DisplayWinCondition)
	}
}

// ── Generate: mandatory content groups ───────────────────────────────

func TestGenerate_MandatoryContentGroupsMatchZones(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyDefault, 3, 2))
	playerGroups, neutralGroups := 0, 0
	for _, mc := range tmpl.MandatoryContent {
		switch {
		case strings.HasPrefix(mc.Name, "mandatory_content_side_"):
			playerGroups++
		case strings.HasPrefix(mc.Name, "mandatory_content_neutral_"):
			neutralGroups++
		}
	}
	if playerGroups != 3 || neutralGroups != 2 {
		t.Errorf("groups = %d player, %d neutral; want 3/2", playerGroups, neutralGroups)
	}
}

// ── Generate: spawn zone main objects ────────────────────────────────

func TestGenerate_SpawnZoneHasSpawnMainObject(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyDefault, 2, 0))
	for _, z := range tmpl.Variants[0].Zones {
		if !strings.HasPrefix(z.Name, "Spawn-") {
			continue
		}
		if len(z.MainObjects) == 0 || z.MainObjects[0].Type != "Spawn" {
			t.Errorf("zone %q: missing Spawn main object", z.Name)
		}
	}
}

func TestGenerate_MultipleCastlesAddedToSpawn(t *testing.T) {
	s := cfgWith(config.TopologyDefault, 2, 0)
	s.ZoneConfiguration.PlayerZoneCastles = 3
	tmpl, _ := services.Generate(s)
	for _, z := range tmpl.Variants[0].Zones {
		if !strings.HasPrefix(z.Name, "Spawn-") {
			continue
		}
		if len(z.MainObjects) != 3 {
			t.Errorf("zone %q: %d main objects, want 3", z.Name, len(z.MainObjects))
		}
	}
}

// ── Generate: border guards scale with quality ───────────────────────

func TestGenerate_BorderGuardsHigherForHighQuality(t *testing.T) {
	mk := func(low, high int) int {
		s := cfgWith(config.TopologyDefault, 2, 0)
		s.ZoneConfiguration.Advanced.Enabled = true
		s.ZoneConfiguration.Advanced.NeutralLowNoCastleCount = low
		s.ZoneConfiguration.Advanced.NeutralHighNoCastleCount = high
		tmpl, _ := services.Generate(s)
		total := 0
		for _, c := range tmpl.Variants[0].Connections {
			total += c.GuardValue
		}
		return total
	}
	if mk(0, 4) <= mk(4, 0) {
		t.Error("expected high-quality border guards > low-quality")
	}
}

// ── Generate: comprehensive topology smoke ───────────────────────────

func TestGenerate_AllTopologiesCompleteWithoutError(t *testing.T) {
	for _, topo := range []config.MapTopology{
		config.TopologyDefault, config.TopologyChain, config.TopologyHubAndSpoke,
		config.TopologySharedWeb, config.TopologyRandom, config.TopologyBalanced,
	} {
		t.Run(string(topo), func(t *testing.T) {
			tmpl, err := services.Generate(cfgWith(topo, 4, 4))
			if err != nil {
				t.Fatal(err)
			}
			if len(tmpl.Variants[0].Zones) == 0 {
				t.Fatal("no zones produced")
			}
		})
	}
}

// ── Generate: connection endpoints resolve ───────────────────────────

func TestGenerate_AllConnectionEndpointsReferenceKnownZones(t *testing.T) {
	tmpl, _ := services.Generate(cfgWith(config.TopologyDefault, 3, 2))
	names := map[string]bool{}
	for _, z := range tmpl.Variants[0].Zones {
		names[z.Name] = true
	}
	for _, c := range tmpl.Variants[0].Connections {
		if !names[c.From] || !names[c.To] {
			t.Errorf("connection %q references unknown zone(s): %s→%s", c.Name, c.From, c.To)
		}
	}
}

// ── Generate: ZoneLetters exported constant ──────────────────────────

func TestZoneLetters_HasAtLeast32Entries(t *testing.T) {
	if len(services.ZoneLetters) < 32 {
		t.Errorf("ZoneLetters len = %d, want >= 32", len(services.ZoneLetters))
	}
}
