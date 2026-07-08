package assetProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlayerZoneIsDrawn_CanvasIsMutated(t *testing.T) {
	// Arrange
	blankPixels := append([]uint8(nil), newCanvas().Pix...)

	// Act
	canvas := renderPlayer(t, preview.PreviewZone{Owner: 1, IsPlayer: true})

	// Assert
	assert.NotEqual(t, blankPixels, canvas.Pix)
}

func TestWhenOwnerIsBelowRange_FallsBackToFirstPlayerSprite(t *testing.T) {
	// Arrange
	firstPlayerCanvas := renderPlayer(t, preview.PreviewZone{Owner: 1, IsPlayer: true})

	// Act
	canvas := renderPlayer(t, preview.PreviewZone{Owner: 0, IsPlayer: true})

	// Assert
	assert.Equal(t, firstPlayerCanvas.Pix, canvas.Pix)
}

func TestWhenOwnerIsAboveRange_FallsBackToLastPlayerSprite(t *testing.T) {
	// Arrange
	lastPlayerCanvas := renderPlayer(t, preview.PreviewZone{Owner: 8, IsPlayer: true})

	// Act
	canvas := renderPlayer(t, preview.PreviewZone{Owner: 99, IsPlayer: true})

	// Assert
	assert.Equal(t, lastPlayerCanvas.Pix, canvas.Pix)
}

func TestWhenOwnersDiffer_DrawsDifferentSprites(t *testing.T) {
	// Arrange
	firstPlayerCanvas := renderPlayer(t, preview.PreviewZone{Owner: 1, IsPlayer: true})

	// Act
	canvas := renderPlayer(t, preview.PreviewZone{Owner: 2, IsPlayer: true})

	// Assert
	assert.NotEqual(t, firstPlayerCanvas.Pix, canvas.Pix)
}
