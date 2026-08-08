package pickerEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheServiceIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := pickers.NewPickerEntryService()

	// Assert
	assert.NotNil(t, service)
}
