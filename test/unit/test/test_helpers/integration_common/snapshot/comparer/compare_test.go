package comparer_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// solidImage builds a width x height RGBA image filled with the given color.
func solidImage(width, height int, fill color.NRGBA) *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, width, height))
	for row := range height {
		for column := range width {
			result.SetRGBA(column, row, color.RGBA{R: fill.R, G: fill.G, B: fill.B, A: 255})
		}
	}
	return result
}

func TestWhenImagesAreIdentical_ReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	fill := color.NRGBA{
		R: uint8(gofakeit.Number(0, 255)),
		G: uint8(gofakeit.Number(0, 255)),
		B: uint8(gofakeit.Number(0, 255)),
	}
	golden := solidImage(4, 3, fill)
	actual := solidImage(4, 3, fill)
	comparer := snapshot.NewComparer()

	// Act
	difference, err := comparer.Compare(golden, actual)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, snapshot.Difference{}, difference)
}

func TestWhenImagesAreFullyInverted_ReturnsOne(t *testing.T) {
	t.Parallel()
	// Arrange
	golden := solidImage(4, 3, color.NRGBA{})
	actual := solidImage(4, 3, color.NRGBA{R: 255, G: 255, B: 255})
	comparer := snapshot.NewComparer()

	// Act
	difference, err := comparer.Compare(golden, actual)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, snapshot.Difference{MeanDistance: 1, ChangedPixelFraction: 1}, difference)
}

func TestWhenSinglePixelDiffers_ReturnsScaledMeanDistance(t *testing.T) {
	t.Parallel()
	// Arrange
	golden := solidImage(2, 2, color.NRGBA{R: 100, G: 100, B: 100})
	actual := solidImage(2, 2, color.NRGBA{R: 100, G: 100, B: 100})
	actual.SetRGBA(1, 0, color.RGBA{R: 110, G: 120, B: 130, A: 255})
	comparer := snapshot.NewComparer()

	// Act
	difference, err := comparer.Compare(golden, actual)

	// Assert
	require.NoError(t, err)
	assert.InDelta(t, 60.0/(2*2*3*255), difference.MeanDistance, 1e-12)
}

func TestWhenPixelDeltaStaysWithinTolerance_CountsNoChangedPixel(t *testing.T) {
	t.Parallel()
	// Arrange
	golden := solidImage(2, 2, color.NRGBA{R: 100, G: 100, B: 100})
	actual := solidImage(2, 2, color.NRGBA{R: 100, G: 100, B: 100})
	actual.SetRGBA(1, 0, color.RGBA{R: 100 + snapshot.DefaultPixelTolerance, G: 100, B: 100, A: 255})
	comparer := snapshot.NewComparer()

	// Act
	difference, err := comparer.Compare(golden, actual)

	// Assert
	require.NoError(t, err)
	assert.Zero(t, difference.ChangedPixelFraction)
}

func TestWhenPixelDeltaExceedsTolerance_CountsTheChangedPixel(t *testing.T) {
	t.Parallel()
	// Arrange
	golden := solidImage(2, 2, color.NRGBA{R: 100, G: 100, B: 100})
	actual := solidImage(2, 2, color.NRGBA{R: 100, G: 100, B: 100})
	actual.SetRGBA(1, 0, color.RGBA{R: 100, G: 100, B: 100 + snapshot.DefaultPixelTolerance + 1, A: 255})
	comparer := snapshot.NewComparer()

	// Act
	difference, err := comparer.Compare(golden, actual)

	// Assert
	require.NoError(t, err)
	assert.InDelta(t, 0.25, difference.ChangedPixelFraction, 1e-12)
}

func TestWhenDimensionsDiffer_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	golden := solidImage(4, 3, color.NRGBA{})
	actual := solidImage(3, 4, color.NRGBA{})
	comparer := snapshot.NewComparer()

	// Act
	_, err := comparer.Compare(golden, actual)

	// Assert
	assert.ErrorContains(t, err, "snapshot dimensions differ")
}

func TestWhenImagesAreEmpty_ReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	golden := image.NewRGBA(image.Rect(0, 0, 0, 0))
	actual := image.NewRGBA(image.Rect(0, 0, 0, 0))
	comparer := snapshot.NewComparer()

	// Act
	difference, err := comparer.Compare(golden, actual)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, snapshot.Difference{}, difference)
}
