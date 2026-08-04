package topologyProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_ReturnsNonNilProvider(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	provider := test_helpers.NewTopologyProvider()

	// Assert
	assert.NotNil(t, provider)
}
