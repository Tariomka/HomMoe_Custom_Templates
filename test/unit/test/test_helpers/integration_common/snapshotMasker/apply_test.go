package snapshotMasker_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

// checkeredImage builds a deterministic non-uniform RGBA image so masked and
// unmasked pixels are distinguishable.
func checkeredImage(width, height int) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	for row := range height {
		for column := range width {
			shade := uint8(50 + 20*((column+row)%2))
			result.SetRGBA(column, row, color.RGBA{R: shade, G: shade, B: shade, A: 255})
		}
	}
	return result
}

// maskedExpectation clones the source and paints the given rects black, the
// way a correct Apply should.
func maskedExpectation(source *image.RGBA, rects ...image.Rectangle) *image.RGBA {
	expected := image.NewRGBA(source.Bounds())
	copy(expected.Pix, source.Pix)
	for _, rect := range rects {
		clamped := rect.Intersect(source.Bounds())
		for row := clamped.Min.Y; row < clamped.Max.Y; row++ {
			for column := clamped.Min.X; column < clamped.Max.X; column++ {
				expected.SetRGBA(column, row, color.RGBA{A: 255})
			}
		}
	}
	return expected
}

func TestWhenRectIsInsideBounds_MasksOnlyThatRegion(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshot := checkeredImage(8, 6)
	maskRect := image.Rect(2, 1, 5, 4)
	expected := maskedExpectation(screenshot, maskRect)
	masker := integration_common.SnapshotMasker{}
	masker.AddRect(maskRect)

	// Act
	masker.Apply(screenshot)

	// Assert
	assert.Equal(t, expected, screenshot)
}

func TestWhenRectExceedsBounds_ClampsToScreenshot(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshot := checkeredImage(8, 6)
	maskRect := image.Rect(5, 3, 20, 20)
	expected := maskedExpectation(screenshot, maskRect)
	masker := integration_common.SnapshotMasker{}
	masker.AddRect(maskRect)

	// Act
	masker.Apply(screenshot)

	// Assert
	assert.Equal(t, expected, screenshot)
}

func TestWhenNoRects_LeavesScreenshotUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshot := checkeredImage(8, 6)
	expected := maskedExpectation(screenshot)
	masker := integration_common.SnapshotMasker{}

	// Act
	masker.Apply(screenshot)

	// Assert
	assert.Equal(t, expected, screenshot)
}

func TestWhenMultipleRects_MasksAllRegions(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshot := checkeredImage(8, 6)
	firstRect := image.Rect(0, 0, 2, 2)
	secondRect := image.Rect(5, 4, 8, 6)
	expected := maskedExpectation(screenshot, firstRect, secondRect)
	masker := integration_common.SnapshotMasker{}
	masker.AddRect(firstRect)
	masker.AddRect(secondRect)

	// Act
	masker.Apply(screenshot)

	// Assert
	assert.Equal(t, expected, screenshot)
}

func TestWhenRectIsFullyOutsideBounds_LeavesScreenshotUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshot := checkeredImage(8, 6)
	expected := maskedExpectation(screenshot)
	masker := integration_common.SnapshotMasker{}
	masker.AddRect(image.Rect(10, 10, 20, 20))

	// Act
	masker.Apply(screenshot)

	// Assert
	assert.Equal(t, expected, screenshot)
}
