package assetProvider_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenArenaMarkerIsDrawn_CanvasIsMutated(t *testing.T) {
	t.Parallel()
	// Arrange
	blankPixels := append([]uint8(nil), newCanvas().Pix...)

	// Act
	canvas := renderArenaMarker(t)

	// Assert
	assert.NotEqual(t, blankPixels, canvas.Pix)
}
