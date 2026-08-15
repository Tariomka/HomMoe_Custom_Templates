package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheNameCarriesTheHubPrefix_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Hub-" + gofakeit.LetterN(2)

	// Act
	result := zone_helpers.IsClusterHubZoneName(zoneName)

	// Assert
	assert.True(t, result)
}

func TestWhenTheNameIsTheSharedHub_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Hub"

	// Act
	result := zone_helpers.IsClusterHubZoneName(zoneName)

	// Assert
	assert.False(t, result)
}

func TestWhenTheNameIsAPlayerZone_IsNotAClusterHub(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Spawn-" + gofakeit.LetterN(2)

	// Act
	result := zone_helpers.IsClusterHubZoneName(zoneName)

	// Assert
	assert.False(t, result)
}
