package vec2_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenFloatVectorConverted_TruncatesComponentsToPoint(t *testing.T) {
	// Arrange
	vector := data.NewVec2(4.8, -2.7)

	// Act
	point := vector.ToPoint()

	// Assert
	assert.Equal(t, image.Pt(4, -2), point)
}

func TestWhenIntVectorConverted_KeepsComponentValues(t *testing.T) {
	// Arrange
	xComponent := gofakeit.Number(-1000, 1000)
	yComponent := gofakeit.Number(-1000, 1000)
	vector := data.NewVec2(xComponent, yComponent)

	// Act
	point := vector.ToPoint()

	// Assert
	assert.Equal(t, image.Pt(xComponent, yComponent), point)
}
