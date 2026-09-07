package zoneNames_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenALabelIsGiven_PrefixesItWithTheHubZonePrefix(t *testing.T) {
	t.Parallel()
	// Arrange
	label := gofakeit.LetterN(2)

	// Act
	result := constants.GetHubZoneNameFor(label)

	// Assert
	assert.Equal(t, "Hub-"+label, result)
}
