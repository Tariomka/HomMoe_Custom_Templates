package zoneNameType_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/zone_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheNameIsExactlyTheSharedHub_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Hub"

	// Act
	result := zone_helpers.IsSharedHubZoneName(zoneName)

	// Assert
	assert.True(t, result)
}

func TestWhenTheNameCarriesTheHubPrefix_IsNotTheSharedHub(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "Hub-" + gofakeit.LetterN(2)

	// Act
	result := zone_helpers.IsSharedHubZoneName(zoneName)

	// Assert
	assert.False(t, result)
}

func TestWhenTheNameOnlyDiffersInCase_IsNotTheSharedHub(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneName := "hub"

	// Act
	result := zone_helpers.IsSharedHubZoneName(zoneName)

	// Assert
	assert.False(t, result)
}
