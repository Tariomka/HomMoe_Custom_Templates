package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenCrossroadsPositionIsProvided_SetsCrossroadsPositionPointerOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedPosition := gofakeit.Number(0, 10)
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithCrossroadsPosition(expectedPosition).Build()

	// Assert
	assert.Equal(t, template_model.Zone{CrossroadsPosition: &expectedPosition}, zone)
}
