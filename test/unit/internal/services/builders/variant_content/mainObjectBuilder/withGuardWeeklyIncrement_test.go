package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardWeeklyIncrementIsProvided_SetsGuardWeeklyIncrementOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedIncrement := gofakeit.Float64Range(0.01, 1)
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithGuardWeeklyIncrement(expectedIncrement).Build()

	// Assert
	assert.Equal(t, entities.MainObject{GuardWeeklyIncrement: expectedIncrement}, mainObject)
}
