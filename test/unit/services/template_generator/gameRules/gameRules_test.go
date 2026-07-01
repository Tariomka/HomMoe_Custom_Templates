package gameRules_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// bonusesFor runs the generator's bonus expansion (via CreateGameRules) for the
// given UI bonus entries and returns the produced raw bonuses.
func bonusesFor(entries ...config.BonusEntry) entities.BonusList {
	configuration := config.NewGeneratorConfig()
	configuration.Bonuses = entries
	return providers.NewGameRulesProvider().CreateGameRules(*configuration).Bonuses
}

func loadExampleTemplate(t *testing.T, name string) entities.RmgTemplate {
	t.Helper()
	path, err := filepath.Abs(filepath.Join("..", "..", "..", "..", "..", "data", "ExampleTemplates", name))
	require.NoError(t, err)
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var tpl entities.RmgTemplate
	require.NoError(t, json.Unmarshal(raw, &tpl))
	return tpl
}

// ── Bonuses ──────────────────────────────────────────────────────────────

func TestCreateBonuses_TownPortalFree_ExpandsToSpellPlusFreeCost(t *testing.T) {
	bonuses := bonusesFor(config.BonusEntry{PresetType: config.BonusTownPortalFree, ReceiverFilter: "start_hero"})
	assert.Equal(t, entities.BonusList{
		{SID: "add_bonus_hero_spell", ReceiverSide: -1, ReceiverFilter: "start_hero", Parameters: []string{"neutral_magic_town_portal"}},
		{SID: "add_bonus_hero_stat", ReceiverSide: -1, ReceiverFilter: "start_hero", Parameters: []string{"magicCostSidSet", "neutral_magic_town_portal", "-999", "0"}},
	}, bonuses)
}

func TestCreateBonuses_Spell_FreeAddsCostOverride(t *testing.T) {
	free := bonusesFor(config.BonusEntry{PresetType: config.BonusSpell, ReceiverFilter: "all_heroes", Param: "magic_fireball", Param2: "1"})
	require.Len(t, free, 2)
	assert.Equal(t, "add_bonus_hero_spell", free[0].SID)
	assert.Equal(t, []string{"magic_fireball"}, free[0].Parameters)
	assert.Equal(t, []string{"magicCostSidSet", "magic_fireball", "-999", "0"}, free[1].Parameters)

	paid := bonusesFor(config.BonusEntry{PresetType: config.BonusSpell, ReceiverFilter: "all_heroes", Param: "magic_fireball", Param2: "0"})
	assert.Len(t, paid, 1)
}

func TestCreateBonuses_StatAndResourcePresets(t *testing.T) {
	cases := []struct {
		preset   config.BonusPresetType
		param    string
		wantSID  string
		wantArgs []string
	}{
		{config.BonusUnitMultiplier, "2", "add_bonus_hero_unit_multipler", []string{"2"}},
		{config.BonusMovementBonus, "300", "add_bonus_hero_stat", []string{"movementBonus", "300"}},
		{config.BonusStartingItem, "some_item", "add_bonus_hero_item", []string{"some_item"}},
		{config.BonusStartingGold, "1000", "add_bonus_res", []string{"gold", "1000"}},
		{config.BonusStartingGems, "5", "add_bonus_res", []string{"gemstones", "5"}},
		{config.BonusStartingCrystals, "5", "add_bonus_res", []string{"crystals", "5"}},
		{config.BonusStartingMercury, "5", "add_bonus_res", []string{"mercury", "5"}},
		{config.BonusStartingWood, "10", "add_bonus_res", []string{"wood", "10"}},
		{config.BonusStartingOre, "10", "add_bonus_res", []string{"ore", "10"}},
	}
	for _, testCase := range cases {
		bonuses := bonusesFor(config.BonusEntry{PresetType: testCase.preset, ReceiverFilter: "start_hero", Param: testCase.param})
		if assert.Len(t, bonuses, 1, testCase.wantSID) {
			assert.Equal(t, testCase.wantSID, bonuses[0].SID)
			assert.Equal(t, -1, bonuses[0].ReceiverSide)
			assert.Equal(t, "start_hero", bonuses[0].ReceiverFilter)
			assert.Equal(t, testCase.wantArgs, bonuses[0].Parameters)
		}
	}
}

func TestCreateBonuses_MultipleEntriesConcatenate(t *testing.T) {
	bonuses := bonusesFor(
		config.BonusEntry{PresetType: config.BonusStartingGold, ReceiverFilter: "start_hero", Param: "1000"},
		config.BonusEntry{PresetType: config.BonusTownPortalFree, ReceiverFilter: "start_hero"},
	)
	assert.Len(t, bonuses, 3) // 1 resource + 2 town-portal
}

func TestCreateBonuses_NoEntriesProducesEmpty(t *testing.T) {
	assert.Empty(t, bonusesFor())
}

// Functional-equivalence check against a real game template: Blitz grants the
// free Town Portal, whose raw bonuses must match our expansion exactly.
func TestCreateBonuses_MatchesBlitzExampleTownPortal(t *testing.T) {
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")
	expansion := bonusesFor(config.BonusEntry{PresetType: config.BonusTownPortalFree, ReceiverFilter: "start_hero"})
	assert.Equal(t, blitz.GameRules.Bonuses, expansion)
}

// ── Value overrides ────────────────────────────────────────────────────────

func TestCreateValueOverrides_ParsesValidLinesAndSkipsJunk(t *testing.T) {
	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = "watchtower=25000\n\n  =5 \nbad_line\ngold_mine = 12000 \nnonnum=abc"
	overrides := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)
	assert.Equal(t, []entities.ValueOverride{
		{SID: "watchtower", Variant: -1, GuardValue: 25000},
		{SID: "gold_mine", Variant: -1, GuardValue: 12000},
	}, overrides)
}

func TestCreateValueOverrides_EmptyReturnsNil(t *testing.T) {
	assert.Nil(t, providers.NewGameRulesProvider().CreateValueOverrides(*config.NewGeneratorConfig()))
}

func TestCreateValueOverrides_ReproducesBlitzExampleSidAndGuard(t *testing.T) {
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")
	require.NotEmpty(t, blitz.ValueOverrides)
	want := blitz.ValueOverrides[0]

	configuration := config.NewGeneratorConfig()
	configuration.ValueOverridesText = fmt.Sprintf("%s=%d", want.SID, want.GuardValue)
	got := providers.NewGameRulesProvider().CreateValueOverrides(*configuration)

	require.Len(t, got, 1)
	assert.Equal(t, want.SID, got[0].SID)
	assert.Equal(t, want.GuardValue, got[0].GuardValue)
	assert.Equal(t, -1, got[0].Variant, "generator applies overrides to all variants")
}

// ── Global bans ─────────────────────────────────────────────────────────────

func TestCreateGlobalBans_ItemsAndMagics(t *testing.T) {
	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = "voodoosh_doll_artifact\nflag_of_truce_artifact"
	configuration.BannedMagics = "magic_armageddon"
	bans := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)
	require.NotNil(t, bans)
	assert.Equal(t, []string{"voodoosh_doll_artifact", "flag_of_truce_artifact"}, bans.Items)
	assert.Equal(t, []string{"magic_armageddon"}, bans.Magics)
}

func TestCreateGlobalBans_EmptyReturnsNil(t *testing.T) {
	assert.Nil(t, providers.NewGameRulesProvider().CreateGlobalBans(*config.NewGeneratorConfig()))
}

func TestCreateGlobalBans_ReproducesBlitzExampleItems(t *testing.T) {
	blitz := loadExampleTemplate(t, "Blitz.rmg.json")
	require.NotNil(t, blitz.GlobalBans)
	require.NotEmpty(t, blitz.GlobalBans.Items)

	configuration := config.NewGeneratorConfig()
	configuration.BannedItems = strings.Join(blitz.GlobalBans.Items, "\n")
	bans := providers.NewGameRulesProvider().CreateGlobalBans(*configuration)

	require.NotNil(t, bans)
	assert.Equal(t, blitz.GlobalBans.Items, bans.Items)
}
