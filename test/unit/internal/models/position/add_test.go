package position_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenPositionIsAdded_AppendsItToTheList(t *testing.T) {
	t.Parallel()
	// Arrange
	existing := data.NewVec2(gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1))
	added := data.NewVec2(gofakeit.Float64Range(0, 1), gofakeit.Float64Range(0, 1))
	positions := models.Positions{existing}

	// Act
	positions.Add(added)

	// Assert
	assert.Equal(t, models.Positions{existing, added}, positions)
}
