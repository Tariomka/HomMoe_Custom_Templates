package hubLabel_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/constants"
	"github.com/stretchr/testify/assert"
)

func TestWhenLabelIsTheSharedHub_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	label := constants.HubZoneName

	// Act
	isHub := constants.IsHubLabel(label)

	// Assert
	assert.True(t, isHub)
}

func TestWhenLabelIsATournamentPerPlayerHub_ReturnsTrue(t *testing.T) {
	t.Parallel()
	// Arrange
	label := constants.HubZonePrefix + "A"

	// Act
	isHub := constants.IsHubLabel(label)

	// Assert
	assert.True(t, isHub)
}

func TestWhenLabelIsAnOrdinaryZoneLabel_ReturnsFalse(t *testing.T) {
	t.Parallel()
	// Arrange
	var hubLabels []string

	// Act & Assert
	for _, label := range constants.GetZoneLabels() {
		if constants.IsHubLabel(label) {
			hubLabels = append(hubLabels, label)
		}
	}
	assert.Empty(t, hubLabels)
}
