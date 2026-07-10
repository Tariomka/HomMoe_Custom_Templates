package orientationBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenModeIsProvided_SetsModeOnBuiltOrientation(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedMode := gofakeit.Word()
	builder := variant_content.NewOrientationBuilder()

	// Act
	orientation := builder.WithMode(expectedMode).Build()

	// Assert
	assert.Equal(t, entities.Orientation{Mode: expectedMode}, orientation)
}
