package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheEditWasMadeByTheGenerator_IsAddedIsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	added := edit.IsAdded()

	// Assert
	assert.False(t, added)
}

func TestWhenTheEditWasMarkedAdded_IsAddedIsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	edit = edit.AsAdded(2)

	// Act
	added := edit.IsAdded()

	// Assert
	assert.True(t, added)
}
