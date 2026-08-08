package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheCountIsClamped_ItStaysInsideTheAllowedRange(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		count    int
		maxCount int
		expected int
	}{
		{"WhenTheCountIsBelowOne_ClampsToOne", 0, 5, 1},
		{"WhenTheCountIsAboveTheMaximum_ClampsToTheMaximum", 9, 5, 5},
		{"WhenTheCountIsInRange_ItIsKept", 3, 5, 3},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			service := zone_content.NewZoneContentEditorService()

			// Act
			count := service.ClampContentCount(testCase.count, testCase.maxCount)

			// Assert
			assert.Equal(t, testCase.expected, count)
		})
	}
}
