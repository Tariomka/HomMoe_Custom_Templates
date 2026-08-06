package fileSystemHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenFileSystemHandlerIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	handler, _ := newHandlerWithMocks()

	// Assert
	assert.NotNil(t, handler)
}
