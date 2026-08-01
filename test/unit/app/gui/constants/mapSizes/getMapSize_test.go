package mapSizes_test

import (
	"testing"

	gui_constants "github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/common"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenSizeIsKnown_ReturnsSameEntryAsInternalConstants(t *testing.T) {
	t.Parallel()
	// Arrange
	mapSizes := common.GetMapSizes(true)
	size := mapSizes[gofakeit.Number(0, len(mapSizes)-1)].Size

	// Act
	result := gui_constants.GetMapSize(size)

	// Assert
	assert.Equal(t, common.GetMapSize(size), result)
}
