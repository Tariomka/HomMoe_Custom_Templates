package math_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenBooleanIsTrue_ReturnsOne(t *testing.T) {
	t.Parallel()
	// Arrange
	boolean := true

	// Act
	converted := helpers.BoolToInt(boolean)

	// Assert
	assert.Equal(t, 1, converted)
}

func TestWhenBooleanIsFalse_ReturnsZero(t *testing.T) {
	t.Parallel()
	// Arrange
	boolean := false

	// Act
	converted := helpers.BoolToInt(boolean)

	// Assert
	assert.Equal(t, 0, converted)
}
