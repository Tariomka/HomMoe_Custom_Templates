package guiHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenTheSpellCountLabelRequested_ReadsAsAPluralCaption(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	label := handler.GetSpellCountLabel(3)

	// Assert
	assert.Equal(t, "3 spells picked", label)
}
