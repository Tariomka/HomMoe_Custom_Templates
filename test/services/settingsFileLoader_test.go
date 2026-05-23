package services_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

// TestLoadSettingsFile_LegacyBalancedFlagUpgradesToTopology verifies the
// one-way migration of older .gen.json files that used the
// ExperimentalBalancedZonePlacement boolean.
func TestLoadSettingsFile_LegacyBalancedFlagUpgradesToTopology(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "legacy.gen.json")
	legacy := map[string]any{
		"templateName":                      "Legacy",
		"experimentalBalancedZonePlacement": true,
		"topology":                          string(generator.TopologyRandom),
	}
	data, err := json.Marshal(legacy)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	loaded, err := services.LoadSettingsFile(path)
	if err != nil {
		t.Fatalf("LoadSettingsFile: %v", err)
	}
	if loaded.Topology != generator.TopologyBalanced {
		t.Errorf("topology = %q, want %q (legacy upgrade)", loaded.Topology, generator.TopologyBalanced)
	}
	if loaded.ExperimentalBalancedZonePlacement {
		t.Error("legacy flag should be cleared after upgrade")
	}
}

// TestSettingsToGenerator_PopulatesNewFields covers the loader paths that
// translate persisted SettingsFile content rows and bonuses into the
// GeneratorSettings model used by the template generator.
func TestSettingsToGenerator_PopulatesNewFields(t *testing.T) {
	sf := models.NewSettingsFile()
	sf.BannedItems = "x"
	sf.BannedMagics = "y"
	sf.ValueOverridesText = "z"
	sf.BonusesJSON = "StartingWood|start_hero|7|"
	sf.PlayerZoneContentRows = []models.ZoneContentRowSave{
		{Sid: "mine_gold", Count: 2, IsMine: true},
	}
	sf.HighNeutralContentRows = []models.ZoneContentRowSave{
		{Sid: "pandora_box", Count: 1},
	}

	gs := services.SettingsToGenerator(sf)
	if gs.BannedItems != "x" || gs.BannedMagics != "y" || gs.ValueOverridesText != "z" {
		t.Errorf("banned/overrides not propagated: %+v", gs)
	}
	if len(gs.Bonuses) != 1 || gs.Bonuses[0].PresetType != generator.BonusStartingWood {
		t.Errorf("bonuses not parsed: %+v", gs.Bonuses)
	}
	if len(gs.PlayerZoneMandatoryContent) != 2 {
		t.Errorf("PlayerZoneMandatoryContent count = %d, want 2 (mine_gold count=2)", len(gs.PlayerZoneMandatoryContent))
	}
	if !gs.PlayerZoneMandatoryContent[0].IsMine {
		t.Error("IsMine flag lost on row → mandatory item conversion")
	}
	if len(gs.HighNeutralMandatoryContent) != 1 || gs.HighNeutralMandatoryContent[0].SID != "pandora_box" {
		t.Errorf("HighNeutralMandatoryContent not populated: %+v", gs.HighNeutralMandatoryContent)
	}
}

// TestRowsToMandatoryContent_RoadDistanceRule ensures the RoadDistance
// label on a row produces a corresponding placement rule.
func TestRowsToMandatoryContent_RoadDistanceRule(t *testing.T) {
	rows := []models.ZoneContentRowSave{
		{Sid: "watchtower", Count: 1, RoadDistance: "Near"},
	}
	items := services.RowsToMandatoryContent(rows)
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	if len(items[0].Rules) != 1 || items[0].Rules[0].Type != "Road" {
		t.Errorf("expected a single Road rule, got %+v", items[0].Rules)
	}
}

// TestRowsToMandatoryContent_GroupBecomesIncludeList verifies the IsGroup
// flag routes the SID into IncludeLists instead of SID.
func TestRowsToMandatoryContent_GroupBecomesIncludeList(t *testing.T) {
	rows := []models.ZoneContentRowSave{
		{Sid: "content_list_building_random_hires_low_tier", Count: 1, IsGroup: true},
	}
	items := services.RowsToMandatoryContent(rows)
	if len(items) != 1 {
		t.Fatalf("items = %d, want 1", len(items))
	}
	if items[0].SID != "" {
		t.Errorf("group row should leave SID empty, got %q", items[0].SID)
	}
	if len(items[0].IncludeLists) != 1 || items[0].IncludeLists[0] != "content_list_building_random_hires_low_tier" {
		t.Errorf("group row should populate IncludeLists, got %+v", items[0].IncludeLists)
	}
}
