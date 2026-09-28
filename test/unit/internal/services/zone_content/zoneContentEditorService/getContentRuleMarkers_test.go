package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/content_rule_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenSeveralRulesCarryMarkers_TheyAreJoinedWithASeparator(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	markers := service.GetContentRuleMarkers([]content_rule_model.ContentRuleDescription{
		{Valid: true, Marker: "G"},
		{Valid: true, Marker: "R"},
	})

	// Assert
	assert.Equal(t, "G · R", markers)
}

func TestWhenARuleIsInvalid_ItsMarkerIsSkipped(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	markers := service.GetContentRuleMarkers([]content_rule_model.ContentRuleDescription{
		{Valid: false, Marker: "G"},
		{Valid: true, Marker: "R"},
	})

	// Assert
	assert.Equal(t, "R", markers)
}

func TestWhenARuleHasNoMarker_ItIsSkipped(t *testing.T) {
	t.Parallel()
	// Arrange
	service := zone_content.NewZoneContentEditorService()

	// Act
	markers := service.GetContentRuleMarkers([]content_rule_model.ContentRuleDescription{
		{Valid: true},
		{Valid: true, Marker: "S"},
	})

	// Assert
	assert.Equal(t, "S", markers)
}
