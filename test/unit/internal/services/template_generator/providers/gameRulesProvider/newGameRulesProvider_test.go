package gameRulesProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

func TestWhenConstructed_ReturnsNonNilProvider(t *testing.T) {
	// Arrange & Act
	provider := providers.NewGameRulesProvider()

	// Assert
	assert.NotNil(t, provider)
}
