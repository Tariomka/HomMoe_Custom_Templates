package vec2_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFloatVectorTransformedToInt_TruncatesComponents(t *testing.T) {
	// Arrange
	vector := data.NewVec2(7.9, -3.2)

	// Act
	transformed := data.Transform[float64, int](vector)

	// Assert
	assert.Equal(t, data.Vec2[int]{X: 7, Y: -3}, transformed)
}

func TestWhenIntVectorTransformedToFloat_KeepsComponentValues(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Number(-1000, 1000)
	yComponent := gofakeit.Number(-1000, 1000)
	vector := data.NewVec2(xComponent, yComponent)

	// Act
	transformed := data.Transform[int, float64](vector)

	// Assert
	assert.Equal(t, data.Vec2[float64]{X: float64(xComponent), Y: float64(yComponent)}, transformed)
}
