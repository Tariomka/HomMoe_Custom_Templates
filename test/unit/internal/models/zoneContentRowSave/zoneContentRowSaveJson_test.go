package zoneContentRowSave_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenRowWithRulesIsSerialized_RoundTripsRules(t *testing.T) {
	t.Parallel()
	// Arrange
	original := models.ZoneContentRowSave{
		Sid:   "dragon_utopia",
		Count: 2,
		Rules: []models.ContentRuleRowSave{
			{Name: "Guarded", IsGuarded: new(true)},
			{Name: "Distance to road", DistanceName: "Far"},
			{Name: "Variant", VariantID: new(1)},
		},
	}
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Act
	var roundTripped models.ZoneContentRowSave
	require.NoError(t, json.Unmarshal(data, &roundTripped))

	// Assert
	assert.Equal(t, original.Rules, roundTripped.Rules)
}

func TestWhenRowIsSerialized_UsesRulesFormatWithoutLegacyFlatFields(t *testing.T) {
	t.Parallel()
	// Arrange
	row := models.ZoneContentRowSave{
		Sid:   "x",
		Count: 1,
		Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: new(true)}},
	}

	// Act
	data, err := json.Marshal(row)

	// Assert
	require.NoError(t, err)
	assert.JSONEq(
		t,
		`{"sid":"x","count":1,"isGroup":false,"rules":[{"name":"Guarded","isGuarded":true}]}`,
		string(data),
	)
}
