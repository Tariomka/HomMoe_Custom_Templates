package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMinimumBaseAngleIsProvided_SetsBaseAngleMinOnBuiltOrientation(t *testing.T) {
	// Arrange
	expectedAngle := gofakeit.Number(1, 359)
	builder := variant_content.NewOrientationBuilder()

	// Act
	orientation := builder.WithBaseAngleMin(expectedAngle).Build()

	// Assert
	assert.Equal(t, entities.Orientation{BaseAngleMin: expectedAngle}, orientation)
}
