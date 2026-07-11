package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenZeroAngleZoneIsProvided_SetsZeroAngleZoneOnBuiltOrientation(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedZone := gofakeit.Word()
	builder := variant_content.NewOrientationBuilder()

	// Act
	orientation := builder.WithZeroAngleZone(expectedZone).Build()

	// Assert
	assert.Equal(t, entities.Orientation{ZeroAngleZone: expectedZone}, orientation)
}
