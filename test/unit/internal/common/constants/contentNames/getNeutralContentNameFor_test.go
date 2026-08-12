package contentNames_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenALabelIsGiven_PrefixesItWithTheNeutralContentPrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	label := gofakeit.LetterN(2)

	// Act
	result := constants.GetNeutralContentNameFor(label)

	// Assert
	assert.Equal(t, "mandatory_content_neutral_"+label, result)
}
