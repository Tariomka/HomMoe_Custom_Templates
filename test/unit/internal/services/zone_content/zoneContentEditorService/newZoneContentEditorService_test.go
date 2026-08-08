package zoneContentEditorService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheServiceIsConstructed_ReturnsAUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := zone_content.NewZoneContentEditorService()

	// Assert
	assert.NotNil(t, service)
}
