package editorState_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenManualConnectionsWereStored_ConnectionsRoundTripWithUserAddedFlag(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()
	connections := []entities.Connection{
		{Name: "A-B", From: "Zone A", To: "Zone B", GuardValue: gofakeit.Number(0, 30000), IsUserAdded: true},
		{Name: "B-C", From: "Zone B", To: "Zone C"},
	}
	state.SetManualEdits(nil, connections)

	// Act
	restored := state.GetManualConnections()

	// Assert
	assert.Equal(t, connections, restored)
}

func TestWhenNoManualConnectionsWereStored_NilConnectionsAreReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	state := models.NewEditorState()

	// Act
	restored := state.GetManualConnections()

	// Assert
	assert.Nil(t, restored)
}
