package blockingIssueError_test

import (
	"errors"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/stretchr/testify/assert"
)

func TestWhenTopologyIsUnknown_UnwrapsToUnsupportedTopology(t *testing.T) {
	t.Parallel()
	// Arrange
	rejection := rejectionFor(t, "NotARealTopology")

	// Act
	actual := errors.Unwrap(rejection)

	// Assert
	assert.Same(t, common_errors.ErrUnsupportedTopology, actual)
}
