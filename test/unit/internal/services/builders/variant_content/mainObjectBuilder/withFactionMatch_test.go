package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactionMatchIsApplied_SetsMatchFactionReferenceWithZeroArgument(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithFactionMatch().Build()

	// Assert
	assert.Equal(t, entities.MainObject{
		Faction: &entities.TypedRef{Type: registry.GetFactionTypeValues().Match, Args: []string{"0"}},
	}, mainObject)
}
