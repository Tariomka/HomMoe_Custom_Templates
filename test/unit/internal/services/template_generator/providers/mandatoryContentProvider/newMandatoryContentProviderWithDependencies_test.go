package mandatoryContentProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreProvided_ReturnsUsableProvider(t *testing.T) {
	t.Parallel()
	// Arrange
	zoneClassifier := zones.NewZoneClassifier()
	zoneEditor := connection_editor.NewZoneEditorService()

	// Act
	provider := providers.NewMandatoryContentProviderWithDependencies(zoneClassifier, zoneEditor)

	// Assert
	assert.NotNil(t, provider)
}
