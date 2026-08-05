package fileService_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenFileServiceIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service, _ := newServiceWithMocks()

	// Assert
	assert.NotNil(t, service)
}
