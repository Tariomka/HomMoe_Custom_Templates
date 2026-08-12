package variantMappingCatalog_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenCatalogIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	catalog := content_rules.NewVariantMappingCatalog()

	// Assert
	assert.NotNil(t, catalog)
}
