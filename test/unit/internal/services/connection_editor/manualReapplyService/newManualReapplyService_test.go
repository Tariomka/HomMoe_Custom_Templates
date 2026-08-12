package manualReapplyService_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreProvided_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := newManualReapplyService()

	// Assert
	assert.NotNil(t, service)
}
