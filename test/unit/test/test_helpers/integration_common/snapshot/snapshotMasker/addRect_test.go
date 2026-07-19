package snapshotMasker_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenRectAdded_ApplyMasksItsPixels(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshotImage := checkeredImage(6, 7)
	masker := snapshot.SnapshotMasker{}

	// Act
	masker.AddRect(image.Rect(1, 1, 3, 3))
	masker.Apply(screenshotImage)

	// Assert
	assert.Equal(t, color.RGBA{A: 255}, screenshotImage.RGBAAt(2, 2))
}
