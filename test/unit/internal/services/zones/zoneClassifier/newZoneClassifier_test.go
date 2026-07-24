package zoneClassifier_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenClassifierIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	classifier := zones.NewZoneClassifier()

	// Assert
	assert.NotNil(t, classifier)
}
