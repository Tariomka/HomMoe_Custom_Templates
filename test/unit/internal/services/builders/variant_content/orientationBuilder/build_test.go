package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsOrientationWithAllAccumulatedValues(t *testing.T) {
	// Arrange
	expectedMode := gofakeit.Word()
	expectedZone := gofakeit.Word()
	expectedMinimumAngle := gofakeit.Number(1, 179)
	expectedMaximumAngle := gofakeit.Number(180, 359)
	builder := variant_content.NewOrientationBuilder()

	// Act
	orientation := builder.
		WithMode(expectedMode).
		WithZeroAngleZone(expectedZone).
		WithBaseAngleMin(expectedMinimumAngle).
		WithBaseAngleMax(expectedMaximumAngle).
		Build()

	// Assert
	assert.Equal(t, entities.Orientation{
		Mode:          expectedMode,
		ZeroAngleZone: expectedZone,
		BaseAngleMin:  expectedMinimumAngle,
		BaseAngleMax:  expectedMaximumAngle,
	}, orientation)
}
