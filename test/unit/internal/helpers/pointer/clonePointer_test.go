package pointer_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSourceIsNil_ReturnsNil(t *testing.T) {
	t.Parallel()
	// Arrange
	var source *int

	// Act
	clone := helpers.ClonePointer(source)

	// Assert
	assert.Nil(t, clone)
}

func TestWhenSourceIsSet_ClonePointsToAnEqualValue(t *testing.T) {
	t.Parallel()
	// Arrange
	source := gofakeit.Int()

	// Act
	clone := helpers.ClonePointer(&source)

	// Assert
	require.NotNil(t, clone)
	assert.Equal(t, source, *clone)
}

func TestWhenSourceIsSet_CloneIsADistinctPointer(t *testing.T) {
	t.Parallel()
	// Arrange
	source := gofakeit.Int()

	// Act
	clone := helpers.ClonePointer(&source)

	// Assert
	assert.NotSame(t, &source, clone)
}

func TestWhenCloneIsMutated_SourceValueIsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	original := gofakeit.IntRange(1, 1000)
	source := original
	clone := helpers.ClonePointer(&source)

	// Act
	*clone = original + 1

	// Assert
	assert.Equal(t, original, source)
}
