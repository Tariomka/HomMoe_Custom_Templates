package pickerEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheFilterIsNormalized_ItIsTrimmedAndLowercased(t *testing.T) {
	t.Parallel()
	// Arrange
	service := pickers.NewPickerEntryService()

	// Act
	filter := service.NormalizePickerFilter("  SwOrd  ")

	// Assert
	assert.Equal(t, "sword", filter)
}
