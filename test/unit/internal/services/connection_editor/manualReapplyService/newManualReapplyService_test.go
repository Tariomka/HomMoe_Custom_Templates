package manualReapplyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := connection_editor.NewManualReapplyService()

	// Assert
	assert.NotNil(t, service)
}
