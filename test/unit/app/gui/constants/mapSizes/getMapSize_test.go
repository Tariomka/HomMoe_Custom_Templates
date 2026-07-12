package mapSizes_test

import (
	"testing"

	gui_constants "github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSizeIsKnown_ReturnsSameEntryAsInternalConstants(t *testing.T) {
	t.Parallel()
	// Arrange
	size := constants.AllMapSizes[gofakeit.Number(0, len(constants.AllMapSizes)-1)].Size

	// Act
	result := gui_constants.GetMapSize(size)

	// Assert
	assert.Equal(t, constants.GetMapSize(size), result)
}
