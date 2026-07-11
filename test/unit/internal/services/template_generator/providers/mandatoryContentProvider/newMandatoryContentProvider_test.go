package mandatoryContentProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_ReturnsNonNilProvider(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	provider := providers.NewMandatoryContentProvider()

	// Assert
	assert.NotNil(t, provider)
}
