package pathResolutionService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
)

func TestWhenPathResolutionServiceIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := file_system.NewPathResolutionService()

	// Assert
	assert.NotNil(t, service)
}
