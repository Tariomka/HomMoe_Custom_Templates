package services_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
)

// TestSettingsToGenerator_PopulatesNewFields covers the loader paths that
// translate persisted SettingsFile content rows and bonuses into the
// GeneratorSettings model used by the template generator.
func TestSettingsToGenerator_PopulatesNewFields(t *testing.T) {
	state := dtos.NewDefaultEditorStateDto()
	state.BannedItems = "x"
	state.BannedMagics = "y"
	state.ValueOverridesText = "z"
	state.BonusesJSON = "StartingWood|start_hero|7|"
	state.PlayerZoneContentRows = []models.ZoneContentRowSave{
		{Sid: "mine_gold", Count: 2, IsMine: true},
	}
	state.HighNeutralContentRows = []models.ZoneContentRowSave{
		{Sid: "pandora_box", Count: 1},
	}

	configuration := mappers.NewConfigMapper().FromEditorState(state)
	if configuration.BannedItems != "x" || configuration.BannedMagics != "y" || configuration.ValueOverridesText != "z" {
		t.Errorf("banned/overrides not propagated: %+v", configuration)
	}
	if len(configuration.Bonuses) != 1 || configuration.Bonuses[0].PresetType != config_inner.BonusStartingWood {
		t.Errorf("bonuses not parsed: %+v", configuration.Bonuses)
	}
	if len(configuration.PlayerZoneMandatoryContent) != 2 {
		t.Errorf("PlayerZoneMandatoryContent count = %d, want 2 (mine_gold count=2)", len(configuration.PlayerZoneMandatoryContent))
	}
	if !configuration.PlayerZoneMandatoryContent[0].IsMine {
		t.Error("IsMine flag lost on row → mandatory item conversion")
	}
	if len(configuration.HighNeutralMandatoryContent) != 1 || configuration.HighNeutralMandatoryContent[0].SID != "pandora_box" {
		t.Errorf("HighNeutralMandatoryContent not populated: %+v", configuration.HighNeutralMandatoryContent)
	}
}

// TestRowsToMandatoryContent_RoadDistanceRule ensures the RoadDistance
// label on a row produces a corresponding placement rule.
// func TestRowsToMandatoryContent_RoadDistanceRule(t *testing.T) {
// 	rows := []models.ZoneContentRowSave{
// 		{Sid: "watchtower", Count: 1, RoadDistance: "Near"},
// 	}
// 	contentProvider := providers.NewMandatoryContentProvider()
// 	items := contentProvider.CreateContentItemsFrom(rows)
// 	if len(items) != 1 {
// 		t.Fatalf("items = %d, want 1", len(items))
// 	}
// 	if len(items[0].Rules) != 1 || items[0].Rules[0].Type != "Road" {
// 		t.Errorf("expected a single Road rule, got %+v", items[0].Rules)
// 	}
// }

// TestRowsToMandatoryContent_GroupBecomesIncludeList verifies the IsGroup
// flag routes the SID into IncludeLists instead of SID.
func TestRowsToMandatoryContent_GroupBecomesIncludeList(t *testing.T) {
	rows := []models.ZoneContentRowSave{
		{Sid: "content_list_building_random_hires_low_tier", Count: 1, IsGroup: true},
	}
	contentProvider := providers.NewMandatoryContentProvider()
	items := contentProvider.CreateContentItemsFrom(rows)
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
