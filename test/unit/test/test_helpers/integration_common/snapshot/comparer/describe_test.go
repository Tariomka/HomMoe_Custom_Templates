package comparer_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenDifferenceIsDescribed_NamesBothMeasurementsAndTheirThresholds(t *testing.T) {
	t.Parallel()
	// Arrange
	comparer := snapshot.NewComparer()
	difference := snapshot.Difference{MeanDistance: 0.0066, ChangedPixelFraction: 0.0342}

	// Act
	description := comparer.Describe(difference)

	// Assert
	assert.Equal(
		t,
		"mean 0.6600% (allowed < 0.2500%), changed pixels 3.4200% (allowed < 0.0500%, tolerance 64/255)",
		description)
}
