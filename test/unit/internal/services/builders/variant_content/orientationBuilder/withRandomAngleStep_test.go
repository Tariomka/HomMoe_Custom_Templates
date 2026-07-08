package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomAngleStepIsProvided_SetsRandomAngleStepOnBuiltOrientation(t *testing.T) {
	// Arrange
	expectedStep := gofakeit.Number(1, 90)
	builder := variant_content.NewOrientationBuilder()

	// Act
	orientation := builder.WithRandomAngleStep(expectedStep).Build()

	// Assert
	assert.Equal(t, entities.Orientation{RandomAngleStep: expectedStep}, orientation)
}
