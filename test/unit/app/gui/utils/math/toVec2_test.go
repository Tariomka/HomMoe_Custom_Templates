package math_test

import (
	"testing"

	"gioui.org/f32"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/stretchr/testify/assert"
)

func TestWhenAGioPointIsConverted_BothAxesLandOnTheVector(t *testing.T) {
	t.Parallel()
	// Arrange - halves stay exact in both float widths, so the round trip is not
	// measuring float32 truncation.
	point := f32.Pt(12.5, -37.25)

	// Act
	vector := utils.ToVec2(point)

	// Assert
	assert.Equal(t, data.NewVec2(12.5, -37.25), vector)
}
