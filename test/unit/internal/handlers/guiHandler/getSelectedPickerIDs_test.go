package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheSelectionIsRequested_TheIdsFollowTheEntryOrder(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()
	entries := []dtos.PickerEntryDto{{ID: "a"}, {ID: "b"}}

	// Act
	ids := handler.GetSelectedPickerIDs(entries, map[string]bool{"b": true, "a": true})

	// Assert
	assert.Equal(t, []string{"a", "b"}, ids)
}
