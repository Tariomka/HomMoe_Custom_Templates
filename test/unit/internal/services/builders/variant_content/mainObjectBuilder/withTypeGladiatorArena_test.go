package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenGladiatorArenaTypeIsChosen_SetsGladiatorArenaTypeOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithTypeGladiatorArena().Build()

	// Assert
	assert.Equal(t, entities.MainObject{Type: "GladiatorArena"}, mainObject)
}
