package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenSimTurnSquadIsChosen_MarksBuiltConnectionAsSimTurnSquad(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.WithSimTurnSquad().Build()

	// Assert
	assert.Equal(t, entities.Connection{SimTurnSquad: true}, connection)
}
