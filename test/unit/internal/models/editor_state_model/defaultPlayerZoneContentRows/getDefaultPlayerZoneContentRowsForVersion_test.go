package defaultPlayerZoneContentRows_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheFirstVersionIsRequested_ReturnsTodaysDefaults(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	rows, known := editor_state_model.GetDefaultPlayerZoneContentRowsForVersion(1)

	// Assert
	assert.Equal(t,
		[]any{editor_state_model.GetDefaultPlayerZoneContentRows(), true},
		[]any{rows, known})
}

func TestWhenAnUnknownVersionIsRequested_ReportsItAsUnknown(t *testing.T) {
	t.Parallel()
	cases := map[string]int{"VersionZero_IsUnknown": 0, "NewerVersion_IsUnknown": 2}
	for name, version := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange & Act
			_, known := editor_state_model.GetDefaultPlayerZoneContentRowsForVersion(version)

			// Assert
			assert.False(t, known)
		})
	}
}
