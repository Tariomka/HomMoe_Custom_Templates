package manualReapplyService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreProvided_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := connection_editor.NewManualReapplyServiceWithDependencies(
		connection_editor.NewZoneEditorService(),
		zone_services.NewZoneClassifier(),
		generation_tuning.NewGenerationTuningFactory(),
	)

	// Assert
	assert.NotNil(t, service)
}
