package contentRuleRow_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenAllFieldsAreEmpty_SerializesToEmptyObject(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{}

	// Act
	data, err := json.Marshal(rule)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "{}", string(data))
}

func TestWhenOnlyNameAndDistanceAreSet_SerializesOnlyThoseFields(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{Name: "Distance to road", DistanceName: "Far"}

	// Act
	data, err := json.Marshal(rule)

	// Assert
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"Distance to road","distanceName":"Far"}`, string(data))
}

func TestWhenPointerFieldIsSetToFalse_StillSerializesField(t *testing.T) {
	t.Parallel()
	// Arrange
	rule := editor_state.ContentRuleRow{Name: "Guarded", IsGuarded: new(false)}

	// Act
	data, err := json.Marshal(rule)

	// Assert
	require.NoError(t, err)
	assert.Contains(t, string(data), `"isGuarded":false`)
}

func TestWhenSerializedRuleIsDeserialized_RoundTripsAllFields(t *testing.T) {
	t.Parallel()
	// Arrange
	original := editor_state.ContentRuleRow{
		Name:            "Variant",
		DistanceName:    "Near",
		IsGuarded:       new(true),
		IsSoloEncounter: new(false),
		VariantID:       new(3),
	}
	data, err := json.Marshal(original)
	require.NoError(t, err)

	// Act
	var roundTripped editor_state.ContentRuleRow
	require.NoError(t, json.Unmarshal(data, &roundTripped))

	// Assert
	assert.Equal(t, original, roundTripped)
}
