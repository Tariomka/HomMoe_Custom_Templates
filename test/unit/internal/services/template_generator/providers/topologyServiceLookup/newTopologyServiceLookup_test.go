package topologyServiceLookup_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_ReturnsNonNilLookup(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	lookup := test_helpers.NewTopologyServiceLookup(test_helpers.NewZoneFactories())

	// Assert
	assert.NotNil(t, lookup)
}
