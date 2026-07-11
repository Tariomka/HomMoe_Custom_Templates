package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenNoGuardWhenOwnedIsChosen_MarksBuiltObjectToRemoveGuardIfOwned(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithNoGuardWhenOwned().Build()

	// Assert
	assert.Equal(t, entities.MainObject{RemoveGuardIfHasOwner: true}, mainObject)
}
