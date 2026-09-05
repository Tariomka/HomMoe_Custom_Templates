package difference_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenDifferenceIsFormatted_RendersBothMeasurementsAsPercentages(t *testing.T) {
	t.Parallel()
	// Arrange
	difference := snapshot.Difference{MeanDistance: 0.0066, ChangedPixelFraction: 0.0342}

	// Act
	formatted := difference.String()

	// Assert
	assert.Equal(t, "mean 0.6600%, changed pixels 3.4200%", formatted)
}
