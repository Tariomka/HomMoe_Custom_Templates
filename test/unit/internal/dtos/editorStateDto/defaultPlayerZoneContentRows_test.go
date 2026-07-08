package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenDefaultRowsAreBuilt_ReturnsFourteenRows(t *testing.T) {
	// Arrange - defaults require no setup.

	// Act
	rows := dtos.DefaultPlayerZoneContentRows()

	// Assert
	assert.Len(t, rows, 14)
}

func TestWhenDefaultRowsAreBuilt_MarksOnlyFirstSevenRowsAsMines(t *testing.T) {
	// Arrange
	expectedMineFlags := []bool{
		true, true, true, true, true, true, true,
		false, false, false, false, false, false, false,
	}

	// Act
	rows := dtos.DefaultPlayerZoneContentRows()

	// Assert
	actualMineFlags := make([]bool, 0, len(rows))
	for _, row := range rows {
		actualMineFlags = append(actualMineFlags, row.IsMine)
	}
	assert.Equal(t, expectedMineFlags, actualMineFlags)
}

func TestWhenDefaultRowsAreBuilt_GuardsEveryRow(t *testing.T) {
	// Arrange - defaults require no setup.

	// Act
	rows := dtos.DefaultPlayerZoneContentRows()

	// Assert
	guardedRowCount := 0
	for _, row := range rows {
		for _, rule := range row.Rules {
			if rule.Name == "Guarded" && rule.IsGuarded != nil && *rule.IsGuarded {
				guardedRowCount++
				break
			}
		}
	}
	assert.Equal(t, len(rows), guardedRowCount)
}
