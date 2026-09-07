package pickerEntry_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheFilterIsNormalized_ItIsTrimmedAndLowercased(t *testing.T) {
	t.Parallel()
	// Arrange
	// Act
	filter := models.NormalizePickerFilter("  SwOrd  ")

	// Assert
	assert.Equal(t, "sword", filter)
}
