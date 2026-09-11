package math_test

import (
	"testing"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenAVectorIsConverted_BothAxesLandOnTheGioPoint(t *testing.T) {
	t.Parallel()
	// Arrange
	vector := data.NewVec2(gofakeit.Float64Range(-500, 500), gofakeit.Float64Range(-500, 500))

	// Act
	point := utils.ToF32Point(vector)

	// Assert
	assert.Equal(t, f32.Pt(float32(vector.X), float32(vector.Y)), point)
}
