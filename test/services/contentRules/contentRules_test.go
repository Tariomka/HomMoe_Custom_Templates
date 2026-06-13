package contentRules_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func intPtr(value int) *int    { return &value }
func boolPtr(value bool) *bool { return &value }

// ── Distance presets ─────────────────────────────────────────────────

func TestDistancePresets_MatchCSharpValues(t *testing.T) {
	assert.Equal(t, content_rules.DistanceVariation{Name: "Next To", Min: 0.05, Max: 0.1}, content_rules.DistanceNextTo)
	assert.Equal(t, content_rules.DistanceVariation{Name: "Near", Min: 0.1, Max: 0.25}, content_rules.DistanceNear)
	assert.Equal(t, content_rules.DistanceVariation{Name: "Medium", Min: 0.25, Max: 0.5}, content_rules.DistanceMedium)
	assert.Equal(t, content_rules.DistanceVariation{Name: "Far", Min: 0.5, Max: 0.75}, content_rules.DistanceFar)
	assert.Equal(t, content_rules.DistanceVariation{Name: "Very Far", Min: 0.75, Max: 0.9}, content_rules.DistanceVeryFar)
}

func TestGetDistanceDisplayNames_Order(t *testing.T) {
	assert.Equal(t, []string{"Next To", "Near", "Medium", "Far", "Very Far"}, content_rules.GetDistanceDisplayNames())
}

func TestGetDistanceVariationByName_KnownAndUnknown(t *testing.T) {
	variation, ok := content_rules.GetDistanceVariationByName("Medium")
	assert.True(t, ok)
	assert.Equal(t, content_rules.DistanceMedium, variation)

	_, ok = content_rules.GetDistanceVariationByName("Whatever")
	assert.False(t, ok)
}

// ── Rule presets ─────────────────────────────────────────────────────

func TestRulePresets_RoadDistance(t *testing.T) {
	rule := content_rules.RoadDistance(content_rules.DistanceNear, 1)
	assert.Equal(t, "Road", rule.Type)
	assert.Equal(t, 0.1, rule.TargetMin)
	assert.Equal(t, 0.25, rule.TargetMax)
	assert.Equal(t, 1, rule.Weight)
}

func TestRulePresets_TownDistance(t *testing.T) {
	rule := content_rules.TownDistance(content_rules.DistanceNear, 1)
	assert.Equal(t, "MainObject", rule.Type)
	assert.Equal(t, []any{"0"}, rule.Args)
	assert.Equal(t, 0.1, rule.TargetMin)
	assert.Equal(t, 0.25, rule.TargetMax)
}

func TestRulePresets_CrossroadsDistance(t *testing.T) {
	rule := content_rules.CrossroadsDistance(content_rules.DistanceMedium, 2)
	assert.Equal(t, "Crossroads", rule.Type)
	assert.Equal(t, 0.25, rule.TargetMin)
	assert.Equal(t, 0.5, rule.TargetMax)
	assert.Equal(t, 2, rule.Weight)
}

// ── Individual rules: metadata, markers, serialization ───────────────

func TestRuleDistanceToRoad_Metadata(t *testing.T) {
	rule := content_rules.NewRuleDistanceToRoad(&content_rules.DistanceFar)
	assert.Equal(t, "Distance to road", rule.Name())
	assert.Equal(t, "R", rule.Marker())
	assert.Equal(t, "Distance to road: Far", rule.DisplayText())

	saved := rule.SerializeToRowSave()
	assert.Equal(t, "Distance to road", saved.Name)
	assert.Equal(t, "Far", saved.DistanceName)
}

func TestRuleDistanceToRoad_DefaultsToMedium(t *testing.T) {
	rule := content_rules.NewRuleDistanceToRoad(nil)
	assert.Equal(t, content_rules.DistanceMedium, rule.Distance)
}

func TestRuleDistanceToTown_Metadata(t *testing.T) {
	rule := content_rules.NewRuleDistanceToTown(&content_rules.DistanceNear)
	assert.Equal(t, "Distance to town", rule.Name())
	assert.Equal(t, "T", rule.Marker())

	saved := rule.SerializeToRowSave()
	assert.Equal(t, "Distance to town", saved.Name)
	assert.Equal(t, "Near", saved.DistanceName)
}

func TestRuleGuarded_MarkerAndSerialization(t *testing.T) {
	guarded := content_rules.NewRuleGuarded(true)
	assert.Equal(t, "Guarded", guarded.Name())
	assert.Equal(t, "G", guarded.Marker())
	savedGuarded := guarded.SerializeToRowSave()
	assert.Equal(t, "Guarded", savedGuarded.Name)
	assert.NotNil(t, savedGuarded.IsGuarded)
	assert.True(t, *savedGuarded.IsGuarded)

	unguarded := content_rules.NewRuleGuarded(false)
	assert.Equal(t, "!G", unguarded.Marker())
	savedUnguarded := unguarded.SerializeToRowSave()
	assert.NotNil(t, savedUnguarded.IsGuarded)
	assert.False(t, *savedUnguarded.IsGuarded)
}

func TestRuleSoloEncounter_MarkerAndSerialization(t *testing.T) {
	solo := content_rules.NewRuleSoloEncounter(true)
	assert.Equal(t, "Solo Encounter", solo.Name())
	assert.Equal(t, "S", solo.Marker())
	saved := solo.SerializeToRowSave()
	assert.Equal(t, "Solo Encounter", saved.Name)
	assert.NotNil(t, saved.IsSoloEncounter)
	assert.True(t, *saved.IsSoloEncounter)

	notSolo := content_rules.NewRuleSoloEncounter(false)
	assert.Equal(t, "!S", notSolo.Marker())
}

func TestRuleVariant_MetadataAndSerialization(t *testing.T) {
	rule, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, intPtr(2))
	assert.NoError(t, err)
	assert.Equal(t, "Variant", rule.Name())
	assert.Equal(t, "", rule.Marker())
	assert.Equal(t, "Variant: Large Guard", rule.DisplayText())

	saved := rule.SerializeToRowSave()
	assert.Equal(t, "Variant", saved.Name)
	assert.NotNil(t, saved.VariantId)
	assert.Equal(t, 2, *saved.VariantId)
}

func TestRuleVariant_InvalidIdReturnsError(t *testing.T) {
	_, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, intPtr(99))
	assert.Error(t, err)
}

// ── ApplyRulesToItem ─────────────────────────────────────────────────

func TestApplyRulesToItem_DistanceToRoadAddsPlacementRule(t *testing.T) {
	item := entities.MandatoryContentItem{SID: "x"}
	near := content_rules.DistanceNear
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleDistanceToRoad(&near),
	})
	assert.Equal(t, 1, len(item.Rules))
	assert.Equal(t, "Road", item.Rules[0].Type)
	assert.Equal(t, 0.1, item.Rules[0].TargetMin)
	assert.Equal(t, 0.25, item.Rules[0].TargetMax)
}

func TestApplyRulesToItem_DistanceToTownAddsMainObjectRule(t *testing.T) {
	item := entities.MandatoryContentItem{SID: "x"}
	near := content_rules.DistanceNear
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleDistanceToTown(&near),
	})
	assert.Equal(t, 1, len(item.Rules))
	assert.Equal(t, "MainObject", item.Rules[0].Type)
	assert.Equal(t, []any{"0"}, item.Rules[0].Args)
}

func TestApplyRulesToItem_GuardedSetsField(t *testing.T) {
	item := entities.MandatoryContentItem{SID: "x"}
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleGuarded(true),
	})
	assert.True(t, item.IsGuarded)
	assert.Empty(t, item.Rules)
}

func TestApplyRulesToItem_SoloEncounterSetsField(t *testing.T) {
	item := entities.MandatoryContentItem{SID: "x"}
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{
		content_rules.NewRuleSoloEncounter(true),
	})
	assert.True(t, item.SoloEncounter)
}

func TestApplyRulesToItem_VariantSetsField(t *testing.T) {
	item := entities.MandatoryContentItem{SID: "dragon_utopia"}
	variant, err := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, intPtr(3))
	assert.NoError(t, err)
	content_rules.ApplyRulesToItem(&item, []content_rules.ContentRule{variant})
	assert.NotNil(t, item.Variant)
	assert.Equal(t, 3, *item.Variant)
}

// ── CreateRuleFromSavedRule round-trips ──────────────────────────────

func TestCreateRuleFromSavedRule_RoundTrips(t *testing.T) {
	utopia := constants.ContentIds.DragonUtopia

	cases := []content_rules.ContentRule{
		content_rules.NewRuleDistanceToRoad(&content_rules.DistanceFar),
		content_rules.NewRuleDistanceToTown(&content_rules.DistanceNear),
		content_rules.NewRuleGuarded(true),
		content_rules.NewRuleGuarded(false),
		content_rules.NewRuleSoloEncounter(true),
	}
	for _, original := range cases {
		saved := original.SerializeToRowSave()
		restored := content_rules.CreateRuleFromSavedRule(saved, models.SidMapping{Sid: "x"})
		assert.NotNil(t, restored, original.Name())
		assert.Equal(t, original.SerializeToRowSave(), restored.SerializeToRowSave())
	}

	variant, _ := content_rules.NewRuleVariant(&content_rules.UtopiaVariants, intPtr(1))
	restoredVariant := content_rules.CreateRuleFromSavedRule(variant.SerializeToRowSave(), utopia)
	assert.NotNil(t, restoredVariant)
	assert.Equal(t, variant.SerializeToRowSave(), restoredVariant.SerializeToRowSave())
}

func TestCreateRuleFromSavedRule_UnknownNameReturnsNil(t *testing.T) {
	restored := content_rules.CreateRuleFromSavedRule(models.ContentRuleRowSave{Name: "Nope"}, models.SidMapping{Sid: "x"})
	assert.Nil(t, restored)
}

func TestCreateRuleFromSavedRule_GuardedWithoutValueReturnsNil(t *testing.T) {
	restored := content_rules.CreateRuleFromSavedRule(models.ContentRuleRowSave{Name: "Guarded"}, models.SidMapping{Sid: "x"})
	assert.Nil(t, restored)
}

// ── RestoreRulesFromRow ──────────────────────────────────────────────

func TestRestoreRulesFromRow_NewFormatUsesSerializedRules(t *testing.T) {
	row := models.ZoneContentRowSave{
		Sid: "x",
		Rules: []models.ContentRuleRowSave{
			{Name: "Guarded", IsGuarded: boolPtr(true)},
			{Name: "Distance to road", DistanceName: "Far"},
		},
	}
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})
	assert.Equal(t, 2, len(rules))
	assert.Equal(t, "Guarded", rules[0].Name())
	assert.Equal(t, "Distance to road", rules[1].Name())
}

func TestRestoreRulesFromRow_LegacyGuardedTrue(t *testing.T) {
	row := models.ZoneContentRowSave{Sid: "x", IsGuarded: true}
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, "Guarded", rules[0].Name())
	assert.Equal(t, "G", rules[0].Marker())
}

func TestRestoreRulesFromRow_LegacyPlainRowAddsUnguarded(t *testing.T) {
	row := models.ZoneContentRowSave{Sid: "x"}
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, "!G", rules[0].Marker())
}

func TestRestoreRulesFromRow_LegacyRoadDistanceAndNearCastle(t *testing.T) {
	row := models.ZoneContentRowSave{Sid: "x", IsGuarded: true, RoadDistance: "Medium", NearCastle: true}
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})
	assert.Equal(t, 3, len(rules))
	assert.Equal(t, "Guarded", rules[0].Name())
	assert.Equal(t, "Distance to road", rules[1].Name())
	assert.Equal(t, "Distance to town", rules[2].Name())
}

func TestRestoreRulesFromRow_LegacyUnknownRoadDistanceSkipped(t *testing.T) {
	row := models.ZoneContentRowSave{Sid: "x", RoadDistance: "Whatever"}
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})
	// Only the unguarded rule; the unknown road label is skipped.
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, "Guarded", rules[0].Name())
}

func TestRestoreRulesFromRow_LegacyAnyRoadDistanceSkipped(t *testing.T) {
	row := models.ZoneContentRowSave{Sid: "x", RoadDistance: "Any"}
	rules := content_rules.RestoreRulesFromRow(row, models.SidMapping{Sid: "x"})
	assert.Equal(t, 1, len(rules))
	assert.Equal(t, "Guarded", rules[0].Name())
}

// ── MigrateLegacyRow ─────────────────────────────────────────────────

func TestMigrateLegacyRow_ConvertsLegacyAndClearsFlatFields(t *testing.T) {
	row := models.ZoneContentRowSave{Sid: "x", Count: 2, IsGuarded: true, NearCastle: true, RoadDistance: "Near"}
	migrated := content_rules.MigrateLegacyRow(row, models.SidMapping{Sid: "x"})

	assert.Equal(t, 2, migrated.Count)
	assert.False(t, migrated.IsGuarded)
	assert.False(t, migrated.NearCastle)
	assert.Equal(t, "", migrated.RoadDistance)
	assert.Equal(t, 3, len(migrated.Rules))
	assert.Equal(t, "Guarded", migrated.Rules[0].Name)
	assert.Equal(t, "Distance to road", migrated.Rules[1].Name)
	assert.Equal(t, "Distance to town", migrated.Rules[2].Name)
}

func TestMigrateLegacyRow_NewFormatRowUnchanged(t *testing.T) {
	row := models.ZoneContentRowSave{
		Sid:   "x",
		Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: boolPtr(true)}},
	}
	migrated := content_rules.MigrateLegacyRow(row, models.SidMapping{Sid: "x"})
	assert.Equal(t, row.Rules, migrated.Rules)
}

// ── Variant mappings ─────────────────────────────────────────────────

func TestGetAllVariantMappings_HasThree(t *testing.T) {
	assert.Equal(t, 3, len(content_rules.GetAllVariantMappings()))
}

func TestGetVariantsForContent_Counts(t *testing.T) {
	assert.Equal(t, 4, len(content_rules.GetVariantsForContent(constants.ContentIds.DragonUtopia)))
	assert.Equal(t, 28, len(content_rules.GetVariantsForContent(constants.ContentIds.PandoraBox)))
	assert.Equal(t, 4, len(content_rules.GetVariantsForContent(constants.ContentIds.MontyHall)))
	assert.Empty(t, content_rules.GetVariantsForContent(constants.ContentIds.Watchtower))
}

func TestGetVariantForContentById(t *testing.T) {
	mapping, ok := content_rules.GetVariantForContentById(constants.ContentIds.DragonUtopia, 2)
	assert.True(t, ok)
	assert.Equal(t, "Large Guard", mapping.Variants[2])

	_, ok = content_rules.GetVariantForContentById(constants.ContentIds.DragonUtopia, 99)
	assert.False(t, ok)
}

func TestGetRules_ReturnsFiveRuleTypes(t *testing.T) {
	rules := content_rules.GetRules()
	assert.Equal(t, 5, len(rules))
	names := make(map[string]bool)
	for _, rule := range rules {
		names[rule.Name()] = true
	}
	assert.True(t, names["Distance to road"])
	assert.True(t, names["Distance to town"])
	assert.True(t, names["Guarded"])
	assert.True(t, names["Variant"])
	assert.True(t, names["Solo Encounter"])
}
