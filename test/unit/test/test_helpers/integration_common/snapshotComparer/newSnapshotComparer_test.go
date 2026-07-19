package snapshotComparer_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers/integration_common"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_UsesDefaultThreshold(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	comparer := integration_common.NewSnapshotComparer()

	// Assert
	assert.InEpsilon(t, integration_common.DefaultSnapshotThreshold, comparer.Threshold, 1e-12)
}
