package connectionNames_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTwoLabelsAreGiven_JoinsThemWithTheFallbackPrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	labelFrom := gofakeit.LetterN(2)
	labelTo := gofakeit.LetterN(2)

	// Act
	result := constants.GetFallbackConnectionNameFor(labelFrom, labelTo)

	// Assert
	assert.Equal(t, "Fallback-"+labelFrom+"-"+labelTo, result)
}
