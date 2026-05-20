package models_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

// TestNewSettingsFile_DefaultsToBalanced ensures the default topology
// matches the C# editor v0.7 default (Balanced) instead of the legacy
// random placement that the Go port shipped with previously.
func TestNewSettingsFile_DefaultsToBalanced(t *testing.T) {
	s := models.NewSettingsFile()
	if s.Topology != generator.TopologyBalanced {
		t.Errorf("default topology = %q, want %q", s.Topology, generator.TopologyBalanced)
	}
}

// TestSettingsFile_RoundTrip verifies that all new C#-parity fields
// survive a JSON encode/decode cycle.
func TestSettingsFile_RoundTrip(t *testing.T) {
	original := models.NewSettingsFile()
	original.TemplateName = "Round Trip"
	original.BannedItems = "bad_item_1,bad_item_2"
	original.BannedMagics = "bad_spell_1"
	original.ValueOverridesText = "watchtower=999"
	original.BonusesJson = "StartingWood|start_hero|10|\nStartingOre|all_heroes|5|"
	original.PlayerZoneContentRows = []models.ZoneContentRowSave{
		{Sid: "watchtower", Count: 2, IsGuarded: true, NearCastle: true, RoadDistance: "Near"},
		{Sid: "mine_gold", Count: 1, IsMine: true, RoadDistance: "Any"},
	}
	original.HubZoneContentRows = []models.ZoneContentRowSave{
		{Sid: "pandora_box", Count: 3, IsGroup: false},
	}

	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	round := &models.SettingsFile{}
	if err := json.Unmarshal(data, round); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if round.BannedItems != original.BannedItems {
		t.Errorf("BannedItems round-trip mismatch: %q", round.BannedItems)
	}
	if round.BannedMagics != original.BannedMagics {
		t.Errorf("BannedMagics round-trip mismatch: %q", round.BannedMagics)
	}
	if round.ValueOverridesText != original.ValueOverridesText {
		t.Errorf("ValueOverridesText round-trip mismatch: %q", round.ValueOverridesText)
	}
	if round.BonusesJson != original.BonusesJson {
		t.Errorf("BonusesJson round-trip mismatch: %q", round.BonusesJson)
	}
	if len(round.PlayerZoneContentRows) != 2 {
		t.Fatalf("PlayerZoneContentRows lost: %d", len(round.PlayerZoneContentRows))
	}
	if round.PlayerZoneContentRows[0].Sid != "watchtower" || round.PlayerZoneContentRows[0].Count != 2 {
		t.Errorf("PlayerZoneContentRows[0] mismatch: %+v", round.PlayerZoneContentRows[0])
	}
	if !round.PlayerZoneContentRows[1].IsMine {
		t.Error("IsMine flag lost in round trip")
	}
	if len(round.HubZoneContentRows) != 1 || round.HubZoneContentRows[0].Sid != "pandora_box" {
		t.Errorf("HubZoneContentRows lost: %+v", round.HubZoneContentRows)
	}
}

// TestBonusEntry_RoundTrip verifies the C# pipe-separated bonus format.
func TestBonusEntry_RoundTrip(t *testing.T) {
	entries := []models.BonusEntry{
		{PresetType: models.BonusStartingWood, ReceiverFilter: "start_hero", Param: "10", Param2: ""},
		{PresetType: models.BonusStartingOre, ReceiverFilter: "all_heroes", Param: "5", Param2: ""},
		{PresetType: models.BonusSpell, ReceiverFilter: "start_hero", Param: "spell_fireball", Param2: "1"},
	}
	encoded := models.SerialiseBonuses(entries)
	decoded := models.ParseBonusesJson(encoded)
	if len(decoded) != len(entries) {
		t.Fatalf("decoded %d entries, want %d", len(decoded), len(entries))
	}
	for i, want := range entries {
		if decoded[i] != want {
			t.Errorf("entry[%d] = %+v, want %+v", i, decoded[i], want)
		}
	}
}

// TestBonusEntry_AcceptsLegacyOrdinalForm checks that BonusesJson values
// produced by older C# builds (numeric preset ordinal) still parse.
func TestBonusEntry_AcceptsLegacyOrdinalForm(t *testing.T) {
	const legacy = "9|start_hero|10|"
	decoded := models.ParseBonusesJson(legacy)
	if len(decoded) != 1 {
		t.Fatalf("legacy parse produced %d entries", len(decoded))
	}
	if decoded[0].PresetType != models.BonusStartingWood {
		t.Errorf("legacy ordinal 9 should map to StartingWood, got %v", decoded[0].PresetType)
	}
}
