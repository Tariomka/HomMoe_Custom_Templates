package bannableItems_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenComparedItemsHaveDifferentCategories_OrdersByCategory(t *testing.T) {
	t.Parallel()
	// Arrange
	combat := constants.BannableItemEntry{Category: "Combat", Name: "Zephyr"}
	magic := constants.BannableItemEntry{Category: "Magic", Name: "Amulet"}

	// Act
	comparison := constants.CompareBannableItems(combat, magic)

	// Assert
	assert.Negative(t, comparison)
}

func TestWhenComparedItemsShareACategory_OrdersByName(t *testing.T) {
	t.Parallel()
	// Arrange
	category := gofakeit.Word()
	first := constants.BannableItemEntry{Category: category, Name: "Amulet"}
	second := constants.BannableItemEntry{Category: category, Name: "Zephyr"}

	// Act
	comparison := constants.CompareBannableItems(first, second)

	// Assert
	assert.Negative(t, comparison)
}

func TestWhenComparedItemsAreIdentical_ReportsEquality(t *testing.T) {
	t.Parallel()
	// Arrange
	item := constants.BannableItemEntry{Category: gofakeit.Word(), Name: gofakeit.Word()}

	// Act
	comparison := constants.CompareBannableItems(item, item)

	// Assert
	assert.Zero(t, comparison)
}
