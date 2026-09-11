package zoneEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/stretchr/testify/assert"
)

func TestWhenZoneListIsEmpty_ReturnsLabelA(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	label := test_helpers.NewZoneEditorService().NextFreeZoneLabel(nil)

	// Assert
	assert.Equal(t, "A", label)
}

func TestWhenFirstLettersAreUsed_ReturnsNextFreeLetter(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-B"},
	}

	// Act
	label := test_helpers.NewZoneEditorService().NextFreeZoneLabel(zones)

	// Assert
	assert.Equal(t, "C", label)
}

func TestWhenSameLetterIsUsedAcrossPrefixes_CountsItOnce(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []template_model.Zone{
		{Name: "Spawn-A"},
		{Name: "Neutral-A"},
	}

	// Act
	label := test_helpers.NewZoneEditorService().NextFreeZoneLabel(zones)

	// Assert
	assert.Equal(t, "B", label)
}

// The label table is finite, and the editor's "Add zone" button asks for a
// label before it knows whether one is left, so the exhausted case has to come
// back empty rather than panic or hand out a duplicate.
func TestWhenEveryLabelIsTaken_ReturnsNoLabel(t *testing.T) {
	t.Parallel()
	// Arrange
	labels := constants.GetZoneLabels()
	zones := make([]template_model.Zone, 0, len(labels))
	for _, label := range labels {
		zones = append(zones, template_model.Zone{Name: "Neutral-" + label})
	}

	// Act
	label := test_helpers.NewZoneEditorService().NextFreeZoneLabel(zones)

	// Assert
	assert.Empty(t, label)
}
