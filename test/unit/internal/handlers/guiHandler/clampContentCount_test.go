package guiHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenTheContentCountIsBelowOne_ItIsClampedToOne(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	count := handler.ClampContentCount(0, 5)

	// Assert
	assert.Equal(t, 1, count)
}
