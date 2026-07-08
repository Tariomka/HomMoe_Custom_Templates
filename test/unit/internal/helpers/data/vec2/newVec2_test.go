package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFloatComponentsGiven_ReturnsVectorWithThoseComponents(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Float64Range(-1000, 1000)
	yComponent := gofakeit.Float64Range(-1000, 1000)

	// Act
	vector := data.NewVec2(xComponent, yComponent)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: xComponent, Y: yComponent}, vector)
}

func TestWhenIntComponentsGiven_ReturnsVectorWithThoseComponents(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Number(-1000, 1000)
	yComponent := gofakeit.Number(-1000, 1000)

	// Act
	vector := data.NewVec2(xComponent, yComponent)

	// Assert
	assert.Equal(t, data.Vec2[int]{X: xComponent, Y: yComponent}, vector)
}
