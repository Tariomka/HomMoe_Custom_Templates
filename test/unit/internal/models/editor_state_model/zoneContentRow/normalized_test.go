package zoneContentRow_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenARowHasNoCount_NormalizingRaisesItToOne(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state_model.ZoneContentRow{Sid: "wood_mine"}

	// Act
	normalized := row.Normalized()

	// Assert
	assert.Equal(t, 1, normalized.Count)
}

func TestWhenARowAlreadyHasACount_NormalizingLeavesItAlone(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state_model.ZoneContentRow{Sid: "wood_mine", Count: 4}

	// Act
	normalized := row.Normalized()

	// Assert
	assert.Equal(t, 4, normalized.Count)
}

func TestWhenARowIsNormalized_TheReceiverIsNotMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	row := editor_state_model.ZoneContentRow{Sid: "wood_mine"}

	// Act
	_ = row.Normalized()

	// Assert
	assert.Equal(t, 0, row.Count)
}
