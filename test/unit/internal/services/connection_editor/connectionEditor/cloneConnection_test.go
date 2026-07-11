package connectionEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func makeArbitraryConnection() entities.Connection {
	return entities.Connection{
		Name:                 gofakeit.Word(),
		From:                 "Spawn-A",
		To:                   "Neutral-B",
		ConnectionType:       "Portal",
		GuardValue:           gofakeit.Number(1, 60000),
		GuardWeeklyIncrement: gofakeit.Float64Range(0.05, 0.25),
		GuardZone:            "Spawn-A",
		GuardMatchGroup:      gofakeit.Word(),
		Length:               gofakeit.Float64Range(0.1, 3.0),
	}
}

func TestWhenIsUserAddedIsRequestedTrue_ReturnsIdenticalCopyWithFlagSet(t *testing.T) {
	t.Parallel()
	// Arrange
	original := makeArbitraryConnection()
	expected := original
	expected.IsUserAdded = true

	// Act
	clone := connection_editor.CloneConnection(original, true)

	// Assert
	assert.Equal(t, expected, clone)
}

func TestWhenIsUserAddedIsRequestedFalse_ReturnsIdenticalCopyWithFlagCleared(t *testing.T) {
	t.Parallel()
	// Arrange
	original := makeArbitraryConnection()
	original.IsUserAdded = true
	expected := original
	expected.IsUserAdded = false

	// Act
	clone := connection_editor.CloneConnection(original, false)

	// Assert
	assert.Equal(t, expected, clone)
}

func TestWhenConnectionIsCloned_LeavesOriginalFlagUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	original := makeArbitraryConnection()

	// Act
	connection_editor.CloneConnection(original, true)

	// Assert
	assert.False(t, original.IsUserAdded)
}
