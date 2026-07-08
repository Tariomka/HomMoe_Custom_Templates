package previewLayout_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameHasSpawnPrefix_ReturnsTrailingLetters(t *testing.T) {
	// Arrange
	zoneName := "Spawn-AB"

	// Act
	letter := services.ExtractZoneLetter(zoneName)

	// Assert
	assert.Equal(t, "AB", letter)
}

func TestWhenNameHasNeutralPrefix_ReturnsTrailingLetter(t *testing.T) {
	// Arrange
	zoneName := "Neutral-C"

	// Act
	letter := services.ExtractZoneLetter(zoneName)

	// Assert
	assert.Equal(t, "C", letter)
}

func TestWhenNameHasNoKnownPrefix_ReturnsNameUnchanged(t *testing.T) {
	// Arrange
	zoneName := "Hub"

	// Act
	letter := services.ExtractZoneLetter(zoneName)

	// Assert
	assert.Equal(t, "Hub", letter)
}
