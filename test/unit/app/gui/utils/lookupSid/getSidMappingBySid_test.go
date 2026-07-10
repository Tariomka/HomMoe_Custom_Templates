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

func TestWhenSidMatches_MappingIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := constants.ContentIDs.DragonUtopia

	// Act
	mapping, _ := utils.GetSidMappingBySid(expected.Sid)

	// Assert
	assert.Equal(t, expected, mapping)
}

func TestWhenSidMatches_LookupReportsFound(t *testing.T) {
	t.Parallel()
	// Arrange
	sid := constants.ContentIDs.DragonUtopia.Sid

	// Act
	_, found := utils.GetSidMappingBySid(sid)

	// Assert
	assert.True(t, found)
}

func TestWhenSidBelongsToIncludeList_MappingIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := constants.IncludeListIDs.ResourceBanksTier2

	// Act
	mapping, _ := utils.GetSidMappingBySid(expected.Sid)

	// Assert
	assert.Equal(t, expected, mapping)
}

func TestWhenSidCaseDiffers_LookupReportsNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	upperCasedSid := strings.ToUpper(constants.ContentIDs.DragonUtopia.Sid)

	// Act
	_, found := utils.GetSidMappingBySid(upperCasedSid)

	// Assert
	assert.False(t, found)
}

func TestWhenSidIsUnknown_LookupReportsNotFound(t *testing.T) {
	t.Parallel()
	// Arrange
	unknownSid := gofakeit.UUID()

	// Act
	_, found := utils.GetSidMappingBySid(unknownSid)

	// Assert
	assert.False(t, found)
}

func TestWhenSidIsUnknown_ZeroMappingIsReturned(t *testing.T) {
	t.Parallel()
	// Arrange
	unknownSid := gofakeit.UUID()

	// Act
	mapping, _ := utils.GetSidMappingBySid(unknownSid)

	// Assert
	assert.Equal(t, models.SidMapping{}, mapping)
}
