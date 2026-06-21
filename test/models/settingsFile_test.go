package models_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config/config_inner"
	"github.com/stretchr/testify/assert"
)

func TestNewSettingsFile_DefaultsToCircles(t *testing.T) {
	s := dtos.NewDefaultEditorStateDto()
	assert.Equal(t, config.TopologyRandom, s.Topology)
}

// func TestSettingsFile_RoundTrip(t *testing.T) {
// 	original := dtos.NewDefaultEditorStateDto()
// 	original.TemplateName = "Round Trip"
// 	original.BannedItems = "bad_item_1,bad_item_2"
// 	original.BannedMagics = "bad_spell_1"
// 	original.ValueOverridesText = "watchtower=999"
// 	original.BonusesJSON = "StartingWood|start_hero|10|\nStartingOre|all_heroes|5|"
// 	original.PlayerZoneContentRows = []models.ZoneContentRowSave{
// 		{Sid: "watchtower", Count: 2, IsGuarded: true, NearCastle: true, RoadDistance: "Near"},
// 		{Sid: "mine_gold", Count: 1, IsMine: true, RoadDistance: "Any"},
// 	}
// 	original.HubZoneContentRows = []models.ZoneContentRowSave{
// 		{Sid: "pandora_box", Count: 3, IsGroup: false},
// 	}

// 	data, err := json.Marshal(original)
// 	assert.NoError(t, err, "marshal")
// 	round := &dtos.EditorStateDto{}
// 	err = json.Unmarshal(data, round)
// 	assert.NoError(t, err, "unmarshal")

// 	assert.Equal(t, original.BannedItems, round.BannedItems, "BannedItems round-trip mismatch")
// 	assert.Equal(t, original.BannedMagics, round.BannedMagics, "BannedMagics round-trip mismatch")
// 	assert.Equal(t, original.ValueOverridesText, round.ValueOverridesText, "ValueOverridesText round-trip mismatch")
// 	assert.Equal(t, original.BonusesJSON, round.BonusesJSON, "BonusesJson round-trip mismatch")
// 	assert.Equal(t, 2, len(round.PlayerZoneContentRows), "PlayerZoneContentRows lost")
// 	assert.Equal(t, "watchtower", round.PlayerZoneContentRows[0].Sid, "PlayerZoneContentRows[0] Sid mismatch")
// 	assert.Equal(t, 2, round.PlayerZoneContentRows[0].Count, "PlayerZoneContentRows[0] Count mismatch")
// 	assert.True(t, round.PlayerZoneContentRows[1].IsMine, "IsMine flag lost in round trip")
// 	assert.Equal(t, 1, len(round.HubZoneContentRows), "HubZoneContentRows lost")
// 	assert.Equal(t, "pandora_box", round.HubZoneContentRows[0].Sid, "HubZoneContentRows[0] Sid mismatch")
// }

func TestBonusEntry_RoundTrip(t *testing.T) {
	entries := []config.BonusEntry{
		{PresetType: config_inner.BonusStartingWood, ReceiverFilter: "start_hero", Param: "10", Param2: ""},
		{PresetType: config_inner.BonusStartingOre, ReceiverFilter: "all_heroes", Param: "5", Param2: ""},
		{PresetType: config_inner.BonusSpell, ReceiverFilter: "start_hero", Param: "spell_fireball", Param2: "1"},
	}
	encoded := config_inner.SerializeBonusEntries(entries)
	decoded := config_inner.DeserializeBonusEntries(encoded)
	assert.Equal(t, len(entries), len(decoded))
	for i, expected := range entries {
		assert.Equal(t, expected, decoded[i])
	}
}

func TestBonusEntry_AcceptsLegacyOrdinalForm(t *testing.T) {
	const legacy = "9|start_hero|10|"
	decoded := config_inner.DeserializeBonusEntries(legacy)
	assert.Equal(t, 1, len(decoded))
	assert.Equal(t, config_inner.BonusStartingWood, decoded[0].PresetType)
}
