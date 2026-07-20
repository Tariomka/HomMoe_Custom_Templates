package assetProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlayerZoneIsDrawn_CanvasIsMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	blankPixels := append([]uint8(nil), newCanvas().Pix...)

	// Act
	canvas := renderPlayer(t, preview.Zone{Owner: 1, Type: preview.ZoneTypePlayer})

	// Assert
	assert.NotEqual(t, blankPixels, canvas.Pix)
}

func TestWhenOwnerIsBelowRange_FallsBackToFirstPlayerSprite(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPlayerCanvas := renderPlayer(t, preview.Zone{Owner: 1, Type: preview.ZoneTypePlayer})

	// Act
	canvas := renderPlayer(t, preview.Zone{Owner: 0, Type: preview.ZoneTypePlayer})

	// Assert
	assert.Equal(t, firstPlayerCanvas.Pix, canvas.Pix)
}

func TestWhenOwnerIsAboveRange_FallsBackToLastPlayerSprite(t *testing.T) {
	t.Parallel()
	// Arrange
	lastPlayerCanvas := renderPlayer(t, preview.Zone{Owner: 8, Type: preview.ZoneTypePlayer})

	// Act
	canvas := renderPlayer(t, preview.Zone{Owner: 99, Type: preview.ZoneTypePlayer})

	// Assert
	assert.Equal(t, lastPlayerCanvas.Pix, canvas.Pix)
}

func TestWhenOwnersDiffer_DrawsDifferentSprites(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPlayerCanvas := renderPlayer(t, preview.Zone{Owner: 1, Type: preview.ZoneTypePlayer})

	// Act
	canvas := renderPlayer(t, preview.Zone{Owner: 2, Type: preview.ZoneTypePlayer})

	// Assert
	assert.NotEqual(t, firstPlayerCanvas.Pix, canvas.Pix)
}
