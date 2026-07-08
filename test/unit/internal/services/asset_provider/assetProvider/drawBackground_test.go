package assetProvider_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenBackgroundIsDrawn_CanvasIsMutated(t *testing.T) {
	// Arrange
	provider := mustNewProvider(t)
	canvas := newCanvas()
	blankPixels := append([]uint8(nil), canvas.Pix...)

	// Act
	provider.DrawBackground(canvas)

	// Assert
	assert.NotEqual(t, blankPixels, canvas.Pix)
}

func TestWhenBackgroundIsDrawn_EveryPixelIsOpaque(t *testing.T) {
	// Arrange
	provider := mustNewProvider(t)
	canvas := newCanvas()

	// Act
	provider.DrawBackground(canvas)

	// Assert
	minAlpha := uint8(255)
	for offset := 3; offset < len(canvas.Pix); offset += 4 {
		minAlpha = min(minAlpha, canvas.Pix[offset])
	}
	assert.Equal(t, uint8(255), minAlpha)
}
