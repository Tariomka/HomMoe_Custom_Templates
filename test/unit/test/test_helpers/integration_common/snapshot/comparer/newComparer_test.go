package comparer_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common/snapshot"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_UsesDefaultThresholds(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	comparer := snapshot.NewComparer()

	// Assert
	assert.Equal(t, snapshot.Comparer{
		MeanThreshold:         snapshot.DefaultMeanThreshold,
		PixelTolerance:        snapshot.DefaultPixelTolerance,
		ChangedPixelThreshold: snapshot.DefaultChangedPixelThreshold,
	}, comparer)
}
