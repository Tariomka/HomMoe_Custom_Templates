package masker_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenRectIsInsideBounds_MasksOnlyThatRegion(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshotImage := checkeredImage(8, 6)
	mask := image.Rect(2, 1, 5, 4)
	expected := maskedExpectation(screenshotImage, mask)
	masker := snapshot.Masker{}
	masker.AddRect(mask)

	// Act
	masker.Apply(screenshotImage)

	// Assert
	assert.Equal(t, expected, screenshotImage)
}

func TestWhenRectExceedsBounds_ClampsToScreenshot(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshotImage := checkeredImage(8, 6)
	mark := image.Rect(5, 3, 20, 20)
	expected := maskedExpectation(screenshotImage, mark)
	masker := snapshot.Masker{}
	masker.AddRect(mark)

	// Act
	masker.Apply(screenshotImage)

	// Assert
	assert.Equal(t, expected, screenshotImage)
}

func TestWhenNoRects_LeavesScreenshotUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshotImage := checkeredImage(8, 6)
	expected := maskedExpectation(screenshotImage)
	masker := snapshot.Masker{}

	// Act
	masker.Apply(screenshotImage)

	// Assert
	assert.Equal(t, expected, screenshotImage)
}

func TestWhenMultipleRects_MasksAllRegions(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshotImage := checkeredImage(8, 6)
	firstMark := image.Rect(0, 0, 2, 2)
	secondMark := image.Rect(5, 4, 8, 6)
	expected := maskedExpectation(screenshotImage, firstMark, secondMark)
	masker := snapshot.Masker{}
	masker.AddRect(firstMark)
	masker.AddRect(secondMark)

	// Act
	masker.Apply(screenshotImage)

	// Assert
	assert.Equal(t, expected, screenshotImage)
}

func TestWhenRectIsFullyOutsideBounds_LeavesScreenshotUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	screenshotImage := checkeredImage(8, 6)
	expected := maskedExpectation(screenshotImage)
	masker := snapshot.Masker{}
	masker.AddRect(image.Rect(10, 10, 20, 20))

	// Act
	masker.Apply(screenshotImage)

	// Assert
	assert.Equal(t, expected, screenshotImage)
}

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

// maskedExpectation clones the source and paints the given marks black, the
// way a correct Apply should.
func maskedExpectation(source *image.RGBA, marks ...image.Rectangle) *image.RGBA {
	expected := image.NewRGBA(source.Bounds())
	copy(expected.Pix, source.Pix)
	for _, rect := range marks {
		clamped := rect.Intersect(source.Bounds())
		for row := clamped.Min.Y; row < clamped.Max.Y; row++ {
			for column := clamped.Min.X; column < clamped.Max.X; column++ {
				expected.SetRGBA(column, row, color.RGBA{A: 255})
			}
		}
	}
	return expected
}
