package bannableItems_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSidIsInTheCatalog_ReturnsItsEntry(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := constants.GetBannableItemsWithExclusions(nil)[0]

	// Act
	item, found := constants.FindBannableItem(expected.Sid)

	// Assert
	assert.True(t, found)
	assert.Equal(t, expected, item)
}

func TestWhenSidIsNotInTheCatalog_ReportsItAsMissing(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	_, found := constants.FindBannableItem(gofakeit.UUID())

	// Assert
	assert.False(t, found)
}

func TestWhenSidIsNotInTheCatalog_ReturnsTheZeroEntry(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	item, _ := constants.FindBannableItem(gofakeit.UUID())

	// Assert
	assert.Equal(t, constants.BannableItemEntry{}, item)
}
