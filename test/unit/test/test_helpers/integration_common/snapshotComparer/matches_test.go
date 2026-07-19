package snapshotComparer_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

func TestWhenDifferenceIsBelowThreshold_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := integration_common.NewSnapshotComparer()

	// Act
	matches := comparer.Matches(0.0099)

	// Assert
	assert.True(t, matches)
}

func TestWhenDifferenceEqualsThreshold_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := integration_common.NewSnapshotComparer()

	// Act
	matches := comparer.Matches(integration_common.DefaultSnapshotThreshold)

	// Assert
	assert.False(t, matches)
}

func TestWhenDifferenceIsAboveThreshold_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := integration_common.NewSnapshotComparer()

	// Act
	matches := comparer.Matches(0.02)

	// Assert
	assert.False(t, matches)
}
