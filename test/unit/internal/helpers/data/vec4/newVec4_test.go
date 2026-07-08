package vec4_test

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
	zComponent := gofakeit.Float64Range(-1000, 1000)
	wComponent := gofakeit.Float64Range(-1000, 1000)

	// Act
	vector := data.NewVec4(xComponent, yComponent, zComponent, wComponent)

	// Assert
	assert.Equal(t, data.Vec4[float64]{X: xComponent, Y: yComponent, Z: zComponent, W: wComponent}, vector)
}

func TestWhenIntComponentsGiven_ReturnsVectorWithThoseComponents(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Number(-1000, 1000)
	yComponent := gofakeit.Number(-1000, 1000)
	zComponent := gofakeit.Number(-1000, 1000)
	wComponent := gofakeit.Number(-1000, 1000)

	// Act
	vector := data.NewVec4(xComponent, yComponent, zComponent, wComponent)

	// Assert
	assert.Equal(t, data.Vec4[int]{X: xComponent, Y: yComponent, Z: zComponent, W: wComponent}, vector)
}
