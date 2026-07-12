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
	canvas := renderPlayer(t, preview.Zone{Owner: 1, IsPlayer: true})

	// Assert
	assert.NotEqual(t, blankPixels, canvas.Pix)
}

func TestWhenOwnerIsBelowRange_FallsBackToFirstPlayerSprite(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPlayerCanvas := renderPlayer(t, preview.Zone{Owner: 1, IsPlayer: true})

	// Act
	canvas := renderPlayer(t, preview.Zone{Owner: 0, IsPlayer: true})

	// Assert
	assert.Equal(t, firstPlayerCanvas.Pix, canvas.Pix)
}

func TestWhenOwnerIsAboveRange_FallsBackToLastPlayerSprite(t *testing.T) {
	t.Parallel()
	// Arrange
	lastPlayerCanvas := renderPlayer(t, preview.Zone{Owner: 8, IsPlayer: true})

	// Act
	canvas := renderPlayer(t, preview.Zone{Owner: 99, IsPlayer: true})

	// Assert
	assert.Equal(t, lastPlayerCanvas.Pix, canvas.Pix)
}

func TestWhenOwnersDiffer_DrawsDifferentSprites(t *testing.T) {
	t.Parallel()
	// Arrange
	firstPlayerCanvas := renderPlayer(t, preview.Zone{Owner: 1, IsPlayer: true})

	// Act
	canvas := renderPlayer(t, preview.Zone{Owner: 2, IsPlayer: true})

	// Assert
	assert.NotEqual(t, firstPlayerCanvas.Pix, canvas.Pix)
}
