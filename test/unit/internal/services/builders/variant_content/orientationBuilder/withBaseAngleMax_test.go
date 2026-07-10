package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMaximumBaseAngleIsProvided_SetsBaseAngleMaxOnBuiltOrientation(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedAngle := gofakeit.Number(1, 359)
	builder := variant_content.NewOrientationBuilder()

	// Act
	orientation := builder.WithBaseAngleMax(expectedAngle).Build()

	// Assert
	assert.Equal(t, entities.Orientation{BaseAngleMax: expectedAngle}, orientation)
}
