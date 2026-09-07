package connectionNames_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenThreeLabelsAreGiven_JoinsThemWithTheTournamentHubRingPrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	labelFrom := gofakeit.LetterN(2)
	labelCentral := gofakeit.LetterN(2)
	labelTo := gofakeit.LetterN(2)

	// Act
	result := constants.GetTournamentHubRingConnectionNameFor(labelFrom, labelCentral, labelTo)

	// Assert
	assert.Equal(t, "THubRing-"+labelFrom+"-"+labelCentral+"-"+labelTo, result)
}
