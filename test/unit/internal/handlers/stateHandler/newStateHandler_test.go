package stateHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenHandlerIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	handler := handlers.NewStateHandler(
		&test_helpers.FileServiceMock{},
		newPassingValidator(),
		mappers.NewEditorStateMapper(),
	)

	// Assert
	assert.NotNil(t, handler)
}
