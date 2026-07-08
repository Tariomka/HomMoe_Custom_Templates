package assetProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenNeutralZoneIsDrawn_CanvasIsMutated(t *testing.T) {
	// Arrange
	blankPixels := append([]uint8(nil), newCanvas().Pix...)

	// Act
	canvas := renderNeutral(t, preview.PreviewZone{Tier: 1})

	// Assert
	assert.NotEqual(t, blankPixels, canvas.Pix)
}

func TestWhenTierIsUnknown_FallsBackToLowQualitySprite(t *testing.T) {
	// Arrange
	lowTierCanvas := renderNeutral(t, preview.PreviewZone{Tier: 1})

	// Act
	canvas := renderNeutral(t, preview.PreviewZone{Tier: 0})

	// Assert
	assert.Equal(t, lowTierCanvas.Pix, canvas.Pix)
}

func TestWhenTierIsMedium_DrawsDifferentSpriteThanLow(t *testing.T) {
	// Arrange
	lowTierCanvas := renderNeutral(t, preview.PreviewZone{Tier: 1})

	// Act
	canvas := renderNeutral(t, preview.PreviewZone{Tier: 2})

	// Assert
	assert.NotEqual(t, lowTierCanvas.Pix, canvas.Pix)
}

func TestWhenTierIsHigh_DrawsDifferentSpriteThanMedium(t *testing.T) {
	// Arrange
	mediumTierCanvas := renderNeutral(t, preview.PreviewZone{Tier: 2})

	// Act
	canvas := renderNeutral(t, preview.PreviewZone{Tier: 3})

	// Assert
	assert.NotEqual(t, mediumTierCanvas.Pix, canvas.Pix)
}

func TestWhenZoneHasCastle_DrawsDifferentSpriteThanWithoutCastle(t *testing.T) {
	// Arrange
	castleLessCanvas := renderNeutral(t, preview.PreviewZone{Tier: 2})

	// Act
	canvas := renderNeutral(t, preview.PreviewZone{Tier: 2, HasCastle: true})

	// Assert
	assert.NotEqual(t, castleLessCanvas.Pix, canvas.Pix)
}
