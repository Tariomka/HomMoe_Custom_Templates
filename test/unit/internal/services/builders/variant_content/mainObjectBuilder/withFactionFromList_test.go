package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenFactionFromListIsApplied_SetsFromListFactionReferenceWithoutArguments(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithFactionFromList().Build()

	// Assert
	assert.Equal(t, template_model.MainObject{
		Faction: &template_model.TypedRef{Type: registry.GetFactionTypeValues().FromList},
	}, mainObject)
}
