package neutralZoneQuality_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutralZone"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardedPoolIsInspected_ReturnsMatchingQuality(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		pool     []string
		expected neutralZone.Quality
	}{
		{"WhenPoolIsTier4_ReturnsHigh", []string{"pool_t4_x"}, neutralZone.QualityHigh},
		{"WhenPoolIsTier5_ReturnsHigh", []string{"pool_t5_x"}, neutralZone.QualityHigh},
		{"WhenPoolIsTier1_ReturnsLow", []string{"pool_t1_x"}, neutralZone.QualityLow},
		{"WhenPoolIsTier2_ReturnsLow", []string{"pool_t2_x"}, neutralZone.QualityLow},
		{"WhenPoolIsTier3_ReturnsMedium", []string{"pool_t3_x"}, neutralZone.QualityMedium},
		{"WhenPoolIsEmpty_ReturnsMedium", nil, neutralZone.QualityMedium},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			zone := entities.Zone{Name: "Neutral-Z", GuardedContentPool: testCase.pool}

			// Act
			quality := neutralZone.GetQualityFrom(zone)

			// Assert
			assert.Equal(t, testCase.expected, quality)
		})
	}
}
