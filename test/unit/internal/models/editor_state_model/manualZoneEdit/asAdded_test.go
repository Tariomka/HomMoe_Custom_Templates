package manualZoneEdit_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenAZoneIsMarkedAdded_OriginIsAdded(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	added := edit.AsAdded(2)

	// Assert
	assert.Equal(t, editor_state.ZoneOriginAdded, added.Origin)
}

func TestWhenAZoneIsMarkedAdded_CastlesAreRecorded(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	added := edit.AsAdded(2)

	// Assert
	assert.Equal(t, new(2), added.Castles)
}

func TestWhenAZoneIsMarkedAdded_TheReceiverIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit

	// Act
	_ = edit.AsAdded(2)

	// Assert
	assert.Nil(t, edit.Castles)
}

func TestWhenAZoneIsMarkedAdded_CastlesDoNotAliasTheArgument(t *testing.T) {
	t.Parallel()
	// Arrange
	var edit editor_state_model.ManualZoneEdit
	first := edit.AsAdded(2)
	second := edit.AsAdded(2)

	// Act
	*first.Castles = 9

	// Assert
	assert.Equal(t, 2, *second.Castles)
}
