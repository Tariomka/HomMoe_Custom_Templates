package comparer_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenBothMeasurementsAreBelowThresholds_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := snapshot.NewComparer()
	difference := snapshot.Difference{
		MeanDistance:         snapshot.DefaultMeanThreshold / 2,
		ChangedPixelFraction: snapshot.DefaultChangedPixelThreshold / 2,
	}

	// Act
	matches := comparer.Matches(difference)

	// Assert
	assert.True(t, matches)
}

func TestWhenMeanDistanceEqualsThreshold_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := snapshot.NewComparer()
	difference := snapshot.Difference{MeanDistance: snapshot.DefaultMeanThreshold}

	// Act
	matches := comparer.Matches(difference)

	// Assert
	assert.False(t, matches)
}

func TestWhenOnlyMeanDistanceExceedsThreshold_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := snapshot.NewComparer()
	difference := snapshot.Difference{MeanDistance: snapshot.DefaultMeanThreshold * 2}

	// Act
	matches := comparer.Matches(difference)

	// Assert
	assert.False(t, matches)
}

func TestWhenChangedPixelFractionEqualsThreshold_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := snapshot.NewComparer()
	difference := snapshot.Difference{ChangedPixelFraction: snapshot.DefaultChangedPixelThreshold}

	// Act
	matches := comparer.Matches(difference)

	// Assert
	assert.False(t, matches)
}

// TestWhenOnlyChangedPixelFractionExceedsThreshold_ReturnsFalse is the case the
// old mean-only gate could not catch: a change concentrated on few pixels.
func TestWhenOnlyChangedPixelFractionExceedsThreshold_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := snapshot.NewComparer()
	difference := snapshot.Difference{ChangedPixelFraction: snapshot.DefaultChangedPixelThreshold * 2}

	// Act
	matches := comparer.Matches(difference)

	// Assert
	assert.False(t, matches)
}
