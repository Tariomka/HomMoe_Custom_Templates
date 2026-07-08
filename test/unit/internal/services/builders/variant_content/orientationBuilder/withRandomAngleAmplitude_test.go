package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRandomAngleAmplitudeIsProvided_SetsRandomAngleAmplitudeOnBuiltOrientation(t *testing.T) {
	// Arrange
	expectedAmplitude := gofakeit.Number(1, 180)
	builder := variant_content.NewOrientationBuilder()

	// Act
	orientation := builder.WithRandomAngleAmplitude(expectedAmplitude).Build()

	// Assert
	assert.Equal(t, entities.Orientation{RandomAngleAmplitude: expectedAmplitude}, orientation)
}
