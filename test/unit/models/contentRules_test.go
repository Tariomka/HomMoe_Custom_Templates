package models_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func boolPtr(value bool) *bool { return &value }
func intPtr(value int) *int    { return &value }

// ── ContentRuleRowSave ───────────────────────────────────────────────

func TestContentRuleRowSave_EmptyOmitsAllFields(t *testing.T) {
	data, err := json.Marshal(models.ContentRuleRowSave{})
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestContentRuleRowSave_OnlySetFieldsAreSerialized(t *testing.T) {
	data, err := json.Marshal(models.ContentRuleRowSave{Name: "Distance to road", DistanceName: "Far"})
	assert.NoError(t, err)

	out := string(data)
	assert.Contains(t, out, "\"name\":\"Distance to road\"")
	assert.Contains(t, out, "\"distanceName\":\"Far\"")
	assert.False(t, strings.Contains(out, "isGuarded"))
	assert.False(t, strings.Contains(out, "variantId"))
	assert.False(t, strings.Contains(out, "isSoloEncounter"))
}

func TestContentRuleRowSave_FalsePointerStillSerialized(t *testing.T) {
	data, err := json.Marshal(models.ContentRuleRowSave{Name: "Guarded", IsGuarded: boolPtr(false)})
	assert.NoError(t, err)
	assert.Contains(t, string(data), "\"isGuarded\":false")
}

func TestContentRuleRowSave_RoundTrip(t *testing.T) {
	original := models.ContentRuleRowSave{
		Name:            "Variant",
		DistanceName:    "Near",
		IsGuarded:       boolPtr(true),
		IsSoloEncounter: boolPtr(false),
		VariantId:       intPtr(3),
	}
	data, err := json.Marshal(original)
	assert.NoError(t, err)

	var round models.ContentRuleRowSave
	assert.NoError(t, json.Unmarshal(data, &round))
	assert.Equal(t, original, round)
}

// ── VariantMapping ───────────────────────────────────────────────────

func TestVariantMapping_DisplayTextUsesLowestKey(t *testing.T) {
	mapping := models.NewVariantMapping(
		models.SidMapping{Sid: "x", Name: "Fallback"},
		map[int]string{2: "second", 0: "first", 1: "middle"},
	)
	assert.Equal(t, "first", mapping.DisplayText())
	assert.Equal(t, "first", mapping.String())
}

func TestVariantMapping_DisplayTextFallsBackToContentName(t *testing.T) {
	mapping := models.NewVariantMapping(models.SidMapping{Sid: "x", Name: "Fallback"}, map[int]string{})
	assert.Equal(t, "Fallback", mapping.DisplayText())
}

// ── ZoneContentRowSave: Rules round-trip ─────────────────────────────

func TestZoneContentRowSave_RulesRoundTrip(t *testing.T) {
	original := models.ZoneContentRowSave{
		Sid:   "dragon_utopia",
		Count: 2,
		Rules: []models.ContentRuleRowSave{
			{Name: "Guarded", IsGuarded: boolPtr(true)},
			{Name: "Distance to road", DistanceName: "Far"},
			{Name: "Variant", VariantId: intPtr(1)},
		},
	}
	data, err := json.Marshal(original)
	assert.NoError(t, err)

	var round models.ZoneContentRowSave
	assert.NoError(t, json.Unmarshal(data, &round))
	assert.Equal(t, original.Rules, round.Rules)
}

func TestZoneContentRowSave_NewFormatRowOmitsLegacyFields(t *testing.T) {
	data, err := json.Marshal(models.ZoneContentRowSave{
		Sid:   "x",
		Count: 1,
		Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: boolPtr(true)}},
	})
	assert.NoError(t, err)

	out := string(data)
	assert.False(t, strings.Contains(out, "nearCastle"))
	assert.False(t, strings.Contains(out, "roadDistance"))
}

// ── ZoneContentRowSave: Normalized ───────────────────────────────────

func TestZoneContentRowSave_NormalizedSeedsLegacyDefaultWhenNoRules(t *testing.T) {
	out := models.ZoneContentRowSave{Sid: "x"}.Normalized()
	assert.Equal(t, 1, out.Count)
	// assert.Equal(t, "Any", out.RoadDistance)
}

func TestZoneContentRowSave_NormalizedDoesNotSeedDefaultWhenRulesPresent(t *testing.T) {
	out := models.ZoneContentRowSave{
		Sid:   "x",
		Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: boolPtr(true)}},
	}.Normalized()
	assert.Equal(t, 1, out.Count)
	// assert.Equal(t, "", out.RoadDistance)
}
