package lookupSid_test

import (
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenNameMatchesExactly_MappingIsReturned(t *testing.T) {
	// Arrange
	expected := constants.ContentIDs.AlchemyLab

	// Act
	mapping, _ := utils.GetSidMappingByName(expected.Name)

	// Assert
	assert.Equal(t, expected, mapping)
}

func TestWhenNameMatchesExactly_LookupReportsFound(t *testing.T) {
	// Arrange
	name := constants.ContentIDs.AlchemyLab.Name

	// Act
	_, found := utils.GetSidMappingByName(name)

	// Assert
	assert.True(t, found)
}

func TestWhenNameDiffersOnlyByCase_MappingIsReturned(t *testing.T) {
	// Arrange
	expected := constants.ContentIDs.WiseOwl
	upperCasedName := strings.ToUpper(expected.Name)

	// Act
	mapping, _ := utils.GetSidMappingByName(upperCasedName)

	// Assert
	assert.Equal(t, expected, mapping)
}

func TestWhenNameBelongsToIncludeList_MappingIsReturned(t *testing.T) {
	// Arrange
	expected := constants.IncludeListIDs.RandomHiresLowTier

	// Act
	mapping, _ := utils.GetSidMappingByName(expected.Name)

	// Assert
	assert.Equal(t, expected, mapping)
}

func TestWhenNameIsUnknown_LookupReportsNotFound(t *testing.T) {
	// Arrange
	unknownName := gofakeit.UUID()

	// Act
	_, found := utils.GetSidMappingByName(unknownName)

	// Assert
	assert.False(t, found)
}

func TestWhenNameIsUnknown_ZeroMappingIsReturned(t *testing.T) {
	// Arrange
	unknownName := gofakeit.UUID()

	// Act
	mapping, _ := utils.GetSidMappingByName(unknownName)

	// Assert
	assert.Equal(t, models.SidMapping{}, mapping)
}
