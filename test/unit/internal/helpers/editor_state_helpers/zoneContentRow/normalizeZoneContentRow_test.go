package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/editor_state_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenCountIsUnset_NormalizesCountToOne(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state.ZoneContentRow{Sid: "x"}

	// Act
	normalized := editor_state_helpers.NormalizeZoneContentRow(row)

	// Assert
	assert.Equal(t, 1, normalized.Count)
}

func TestWhenRulesArePresentAndCountIsUnset_StillNormalizesCountToOne(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state.ZoneContentRow{
		Sid:   "x",
		Rules: []editor_state.ContentRuleRow{{Name: "Guarded", IsGuarded: new(true)}},
	}

	// Act
	normalized := editor_state_helpers.NormalizeZoneContentRow(row)

	// Assert
	assert.Equal(t, 1, normalized.Count)
}

func TestWhenCountIsPositive_KeepsCountUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state.ZoneContentRow{Sid: "x", Count: 3}

	// Act
	normalized := editor_state_helpers.NormalizeZoneContentRow(row)

	// Assert
	assert.Equal(t, 3, normalized.Count)
}
