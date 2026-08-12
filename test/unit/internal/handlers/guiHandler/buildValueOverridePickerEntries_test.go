package guiHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenValueOverrideEntriesAreBuilt_TheSidBecomesTheLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	entries := handler.BuildValueOverridePickerEntries([]string{"Sid.Gold"})

	// Assert
	assert.Equal(t, "Sid.Gold", entries[0].Label)
}
