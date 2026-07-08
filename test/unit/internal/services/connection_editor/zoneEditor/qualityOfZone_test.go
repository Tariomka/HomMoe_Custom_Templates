package zoneEditor_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardedPoolIsInspected_ReturnsMatchingQuality(t *testing.T) {
	testCases := []struct {
		name     string
		pool     []string
		expected models.NeutralZoneQuality
	}{
		{"WhenPoolIsTier4_ReturnsHigh", []string{"pool_t4_x"}, models.QualityHigh},
		{"WhenPoolIsTier5_ReturnsHigh", []string{"pool_t5_x"}, models.QualityHigh},
		{"WhenPoolIsTier1_ReturnsLow", []string{"pool_t1_x"}, models.QualityLow},
		{"WhenPoolIsTier2_ReturnsLow", []string{"pool_t2_x"}, models.QualityLow},
		{"WhenPoolIsTier3_ReturnsMedium", []string{"pool_t3_x"}, models.QualityMedium},
		{"WhenPoolIsEmpty_ReturnsMedium", nil, models.QualityMedium},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			zone := entities.Zone{Name: "Neutral-Z", GuardedContentPool: testCase.pool}

			// Act
			quality := connection_editor.QualityOfZone(zone)

			// Assert
			assert.Equal(t, testCase.expected, quality)
		})
	}
}
