package assetProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsDrawn_CanvasIsMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	blankPixels := append([]uint8(nil), newCanvas().Pix...)

	// Act
	canvas := renderNeutral(t, preview.Zone{Tier: 1})

	// Assert
	assert.NotEqual(t, blankPixels, canvas.Pix)
}

func TestWhenTierIsPlastic_DrawsDifferentSpriteThanLow(t *testing.T) {
	t.Parallel()
	// Arrange
	lowTierCanvas := renderNeutral(t, preview.Zone{Tier: preview.TierBronze})

	// Act
	canvas := renderNeutral(t, preview.Zone{Tier: preview.TierPlastic})

	// Assert
	assert.NotEqual(t, lowTierCanvas.Pix, canvas.Pix)
}

func TestWhenTierIsMedium_DrawsDifferentSpriteThanLow(t *testing.T) {
	t.Parallel()
	// Arrange
	lowTierCanvas := renderNeutral(t, preview.Zone{Tier: 1})

	// Act
	canvas := renderNeutral(t, preview.Zone{Tier: 2})

	// Assert
	assert.NotEqual(t, lowTierCanvas.Pix, canvas.Pix)
}

func TestWhenTierIsHigh_DrawsDifferentSpriteThanMedium(t *testing.T) {
	t.Parallel()
	// Arrange
	mediumTierCanvas := renderNeutral(t, preview.Zone{Tier: 2})

	// Act
	canvas := renderNeutral(t, preview.Zone{Tier: 3})

	// Assert
	assert.NotEqual(t, mediumTierCanvas.Pix, canvas.Pix)
}

func TestWhenTierIsPlatinum_DrawsDifferentSpriteThanHigh(t *testing.T) {
	t.Parallel()
	// Arrange
	highTierCanvas := renderNeutral(t, preview.Zone{Tier: preview.TierGold})

	// Act
	canvas := renderNeutral(t, preview.Zone{Tier: preview.TierPlatinum})

	// Assert
	assert.NotEqual(t, highTierCanvas.Pix, canvas.Pix)
}

func TestWhenZoneHasCastle_DrawsDifferentSpriteThanWithoutCastle(t *testing.T) {
	t.Parallel()
	// Arrange
	castleLessCanvas := renderNeutral(t, preview.Zone{Tier: 2})

	// Act
	canvas := renderNeutral(t, preview.Zone{Tier: 2, HasCastle: true})

	// Assert
	assert.NotEqual(t, castleLessCanvas.Pix, canvas.Pix)
}
