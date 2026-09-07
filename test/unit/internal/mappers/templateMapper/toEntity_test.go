package templateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The round trip is the drift guard for all thirty converter pairs: a field
// added to the schema and forgotten in a converter comes back zeroed here.
func TestWhenAFullyPopulatedTemplateRoundTrips_TheEntityComesBackUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewTemplateMapper()
	expected := test_helpers.NewAllFieldsTemplate()
	// The model carries one MapSize, so only a square template round trips.
	expected.SizeZ = expected.SizeX

	// Act
	actual := mapper.ToEntity(mapper.ToModel(expected))

	// Assert
	assert.Equal(t, expected, actual)
}

// nil and empty are different on the wire - null versus [] - so the converters
// have to preserve the distinction rather than normalising it away.
func TestWhenATemplateHasEmptyCollections_TheyStayEmptyRatherThanBecomingNil(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewTemplateMapper()
	expected := test_helpers.NewAllFieldsTemplate()
	expected.SizeZ = expected.SizeX
	expected.ContentPools = []template_entity.ContentPool{}
	expected.ContentLists = []template_entity.ContentList{}

	// Act
	actual := mapper.ToEntity(mapper.ToModel(expected))

	// Assert
	assert.Equal(t, expected, actual)
}

// The schema has no field for a zone's tier, so flattening has to drop it
// rather than fail.
func TestWhenAZoneCarriesATier_TheEntityIsProducedWithoutIt(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewTemplateMapper()
	model := mapper.ToModel(test_helpers.NewAllFieldsTemplate())
	require.NotEmpty(t, model.Variants)
	require.NotEmpty(t, model.Variants[0].Zones)
	model.Variants[0].Zones[0].Quality = new(neutral_zone.QualityHighest)

	// Act
	actual := mapper.ToEntity(model)

	// Assert
	expected := test_helpers.NewAllFieldsTemplate()
	expected.SizeZ = expected.SizeX
	assert.Equal(t, expected, actual)
}

// The model has a single MapSize where the schema has two dimensions, so
// flattening writes it to both and a non-square template comes back square.
func TestWhenAModelIsFlattened_BothMapDimensionsGetTheMapSize(t *testing.T) {
	t.Parallel()
	// Arrange
	mapper := mappers.NewTemplateMapper()
	entity := test_helpers.NewAllFieldsTemplate()
	entity.SizeX = 160
	entity.SizeZ = 96

	// Act
	actual := mapper.ToEntity(mapper.ToModel(entity))

	// Assert
	assert.Equal(t, 160, actual.SizeZ)
}
