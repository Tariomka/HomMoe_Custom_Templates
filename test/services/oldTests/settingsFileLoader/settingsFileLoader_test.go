package settingsFileLoader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// ── LoadSettingsFile ─────────────────────────────────────────────────

func TestLoadSettingsFile_MissingFile_ReturnsError(t *testing.T) {
	_, err := services.LoadSettingsFile(filepath.Join(t.TempDir(), "missing.gen.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadSettingsFile_InvalidJSON_ReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.gen.json")
	if err := os.WriteFile(path, []byte("{not json"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := services.LoadSettingsFile(path); err == nil {
		t.Fatal("expected JSON unmarshal error")
	}
}

func TestLoadSettingsFile_ValidJSON_LoadsFields(t *testing.T) {
	path := filepath.Join(t.TempDir(), "ok.gen.json")
	body := `{"templateName":"X","playerCount":4,"mapSize":192}`
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	state, err := services.LoadSettingsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if state.TemplateName != "X" {
		t.Errorf("TemplateName = %q, want X", state.TemplateName)
	}
	if state.PlayerCount != 4 {
		t.Errorf("PlayerCount = %d, want 4", state.PlayerCount)
	}
	if state.MapSize != 192 {
		t.Errorf("MapSize = %d, want 192", state.MapSize)
	}
}

// ── SaveSettingsFile ─────────────────────────────────────────────────

func TestSaveSettingsFile_WritesIndentedJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.gen.json")
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = "Saved"
	if err := services.SaveSettingsFile(path, &state); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Error("file is empty")
	}
}

func TestSaveSettingsFile_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "rt.gen.json")
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = "RT"
	state.PlayerCount = 6
	if err := services.SaveSettingsFile(path, &state); err != nil {
		t.Fatal(err)
	}
	loaded, err := services.LoadSettingsFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.TemplateName != "RT" || loaded.PlayerCount != 6 {
		t.Errorf("round-trip mismatch: %+v", loaded)
	}
}

func TestSaveSettingsFile_BadPath_ReturnsError(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	if err := services.SaveSettingsFile(filepath.Join(t.TempDir(), "missing_dir", "x.json"), &state); err == nil {
		t.Fatal("expected write error")
	}
}

// ── SettingsToGenerator ──────────────────────────────────────────────

func TestSettingsToGenerator_CopiesScalarFields(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.TemplateName = "T"
	state.GameMode = "SingleHero"
	state.PlayerCount = 5
	state.MapSize = 224
	state.Topology = config.TopologyChain
	state.GenerateRoads = true
	state.RandomPortals = true
	state.MaxPortalConnections = 7
	state.MinNeutralZonesBetweenPlayers = 2
	state.MatchPlayerCastleFactions = true
	state.BannedItems = "ItemA"
	state.BannedMagics = "MagicA"
	state.ValueOverridesText = "ov"
	state.FactionLawXpPercent = 120
	state.AstrologyXpPercent = 80

	gs := services.SettingsToGenerator(&state)
	if gs.TemplateName != "T" || gs.GameMode != "SingleHero" || gs.PlayerCount != 5 || gs.MapSize != 224 {
		t.Errorf("scalar mismatch: %+v", gs)
	}
	if gs.Topology != config.TopologyChain {
		t.Errorf("Topology = %v", gs.Topology)
	}
	if !gs.GenerateRoads || !gs.RandomPortals || gs.MaxPortalConnections != 7 {
		t.Errorf("flags not propagated: %+v", gs)
	}
	if gs.MinNeutralZonesBetweenPlayers != 2 || !gs.MatchPlayerCastleFactions {
		t.Errorf("extra flags not propagated: %+v", gs)
	}
	if gs.BannedItems != "ItemA" || gs.BannedMagics != "MagicA" || gs.ValueOverridesText != "ov" {
		t.Errorf("text fields not propagated: %+v", gs)
	}
	if gs.FactionLawsExpPercent != 120 || gs.AstrologyExpPercent != 80 {
		t.Errorf("xp percents not propagated")
	}
}

func TestSettingsToGenerator_ZoneConfigurationPopulated(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.NeutralZoneCount = 5
	state.PlayerZoneCastles = 2
	state.NeutralZoneCastles = 3
	state.ResourceDensityPercent = 150
	state.StructureDensityPercent = 80
	state.NeutralStackStrengthPercent = 110
	state.BorderGuardStrengthPercent = 90
	state.HubZoneSize = 1.5
	state.HubZoneCastles = 1
	state.AdvancedMode = true
	state.NeutralLowNoCastleCount = 1
	state.NeutralMediumCastleCount = 2
	state.NeutralHighCastleCount = 1
	state.PlayerZoneSize = 1.2
	state.NeutralZoneSize = 0.8
	state.GuardRandomization = 0.1

	gs := services.SettingsToGenerator(&state)
	zc := gs.ZoneConfiguration
	if zc.NeutralZoneCount != 5 || zc.PlayerZoneCastles != 2 || zc.NeutralZoneCastles != 3 {
		t.Errorf("base zone config mismatch: %+v", zc)
	}
	if !zc.Advanced.Enabled || zc.Advanced.NeutralLowNoCastleCount != 1 {
		t.Errorf("advanced not propagated: %+v", zc.Advanced)
	}
	if zc.Advanced.PlayerZoneSize != 1.2 || zc.Advanced.GuardRandomization != 0.1 {
		t.Errorf("advanced sizes not propagated")
	}
}

func TestSettingsToGenerator_HeroSettings(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.HeroCountMin = 2
	state.HeroCountMax = 9
	state.HeroCountIncrement = 3

	gs := services.SettingsToGenerator(&state)
	if gs.HeroSettings.HeroCountMin != 2 || gs.HeroSettings.HeroCountMax != 9 || gs.HeroSettings.HeroCountIncrement != 3 {
		t.Errorf("hero settings mismatch: %+v", gs.HeroSettings)
	}
}

func TestSettingsToGenerator_GameEndConditions_ManualCityHold(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.VictoryCondition = "win_condition_1"
	state.CityHold = true
	state.CityHoldDays = 10
	state.LostStartCity = true
	state.LostStartCityDay = 4
	state.LostStartHero = true

	gs := services.SettingsToGenerator(&state)
	if !gs.GameEndConditions.CityHold {
		t.Error("expected manual CityHold=true")
	}
	if gs.GameEndConditions.CityHoldDays != 10 {
		t.Errorf("CityHoldDays = %d", gs.GameEndConditions.CityHoldDays)
	}
	if !gs.GameEndConditions.LostStartCity || gs.GameEndConditions.LostStartCityDay != 4 || !gs.GameEndConditions.LostStartHero {
		t.Errorf("lost-start fields not propagated: %+v", gs.GameEndConditions)
	}
}

func TestSettingsToGenerator_GameEndConditions_AutoCityHoldFromWinCondition5(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.VictoryCondition = "win_condition_5"
	state.CityHold = false

	gs := services.SettingsToGenerator(&state)
	if !gs.GameEndConditions.CityHold {
		t.Error("win_condition_5 should force CityHold=true even when flag is false")
	}
}

func TestSettingsToGenerator_GladiatorArenaRules(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.GladiatorArena = true
	state.GladiatorArenaDaysDelayStart = 12
	state.GladiatorArenaCountDay = 4

	gs := services.SettingsToGenerator(&state)
	if !gs.GladiatorArenaRules.Enabled || gs.GladiatorArenaRules.DaysDelayStart != 12 || gs.GladiatorArenaRules.CountDay != 4 {
		t.Errorf("gladiator rules mismatch: %+v", gs.GladiatorArenaRules)
	}
}

func TestSettingsToGenerator_TournamentRules(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.Tournament = true
	state.TournamentFirstTournamentDay = 21
	state.TournamentInterval = 5
	state.TournamentPointsToWin = 4
	state.TournamentSaveArmy = true

	gs := services.SettingsToGenerator(&state)
	tr := gs.TournamentRules
	if !tr.Enabled || tr.FirstTournamentDay != 21 || tr.Interval != 5 || tr.PointsToWin != 4 || !tr.SaveArmy {
		t.Errorf("tournament rules mismatch: %+v", tr)
	}
}

func TestSettingsToGenerator_MandatoryContentRowsExpandedAcrossAllZones(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.PlayerZoneContentRows = []models.ZoneContentRowSave{{Sid: "a", Count: 1}}
	state.LowNeutralContentRows = []models.ZoneContentRowSave{{Sid: "b", Count: 1}}
	state.MediumNeutralContentRows = []models.ZoneContentRowSave{{Sid: "c", Count: 1}}
	state.HighNeutralContentRows = []models.ZoneContentRowSave{{Sid: "d", Count: 1}}
	state.HubZoneContentRows = []models.ZoneContentRowSave{{Sid: "e", Count: 1}}

	gs := services.SettingsToGenerator(&state)
	if len(gs.PlayerZoneMandatoryContent) != 1 ||
		len(gs.LowNeutralMandatoryContent) != 1 ||
		len(gs.MediumNeutralMandatoryContent) != 1 ||
		len(gs.HighNeutralMandatoryContent) != 1 ||
		len(gs.HubZoneMandatoryContent) != 1 {
		t.Errorf("content rows lost during conversion: %+v", gs)
	}
}

func TestSettingsToGenerator_BonusesParsedFromJSON(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.BonusesJSON = "StartingWood|start_hero|5|"

	gs := services.SettingsToGenerator(&state)
	if len(gs.Bonuses) != 1 {
		t.Errorf("Bonuses parsed = %d, want 1", len(gs.Bonuses))
	}
}
