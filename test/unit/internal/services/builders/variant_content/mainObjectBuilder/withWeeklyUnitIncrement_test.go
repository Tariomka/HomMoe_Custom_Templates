package mainObjectBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenWeeklyUnitIncrementIsChosen_EnablesWeeklyUnitIncrementOnBuiltObject(t *testing.T) {
	t.Parallel()
	// Arrange
	builder := variant_content.NewObjectBuilder()

	// Act
	mainObject := builder.WithWeeklyUnitIncrement().Build()

	// Assert
	assert.Equal(t, entities.MainObject{EnableWeeklyUnitIncrement: true}, mainObject)
}
