package directoryBrowserService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/stretchr/testify/assert"
)

func TestWhenDirectoryBrowserServiceIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := file_system.NewDirectoryBrowserService()

	// Assert
	assert.NotNil(t, service)
}
