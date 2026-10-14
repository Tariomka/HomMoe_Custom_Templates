package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoQualityIsRecorded_GetQualityReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	quality := edit.GetQuality()

	// Assert
	assert.Nil(t, quality)
}

func TestWhenTheLowestTierIsRecorded_GetQualityReturnsTheLowestTier(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit.Quality = new(int8(neutral_zone.QualityLowest))

	// Act
	quality := edit.GetQuality()

	// Assert
	assert.Equal(t, new(neutral_zone.QualityLowest), quality)
}

func TestWhenATierIsRecorded_GetQualityReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit.Quality = new(int8(neutral_zone.QualityHighest))

	// Act
	quality := edit.GetQuality()

	// Assert
	assert.Equal(t, new(neutral_zone.QualityHighest), quality)
}
