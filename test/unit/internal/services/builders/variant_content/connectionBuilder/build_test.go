package connectionBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenMultipleOptionsAreChained_ReturnsConnectionWithAllAccumulatedValues(t *testing.T) {
	// Arrange
	expectedName := gofakeit.Word()
	expectedFrom := gofakeit.Word()
	expectedTo := gofakeit.Word()
	expectedGuardValue := gofakeit.Number(1, 60000)
	builder := variant_content.NewConnectionBuilder()

	// Act
	connection := builder.
		WithName(expectedName).
		WithFrom(expectedFrom).
		WithTo(expectedTo).
		WithConnectionTypeDirect().
		WithGuardValue(expectedGuardValue).
		Build()

	// Assert
	assert.Equal(t, entities.Connection{
		Name:           expectedName,
		From:           expectedFrom,
		To:             expectedTo,
		ConnectionType: "Direct",
		GuardValue:     expectedGuardValue,
	}, connection)
}
