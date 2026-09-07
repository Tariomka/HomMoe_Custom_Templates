package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardRandomizationIsProvided_SetsGuardRandomizationOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedRandomization := gofakeit.Float64Range(0.01, 1)
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithGuardRandomization(expectedRandomization).Build()

	// Assert
	assert.Equal(t, template_model.MainObject{GuardRandomization: expectedRandomization}, mainObject)
}
