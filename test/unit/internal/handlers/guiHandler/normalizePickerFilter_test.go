package guiHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenAPickerFilterIsNormalized_ItIsTrimmedAndLowercased(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	filter := handler.NormalizePickerFilter("  SwOrd ")

	// Assert
	assert.Equal(t, "sword", filter)
}
