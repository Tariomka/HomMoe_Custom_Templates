package variantMappingCatalog_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenVariantIdExists_ReturnsItsSingleEntryMapping(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewVariantMappingCatalog()

	// Act
	mapping, ok := catalog.GetVariantForContentByID(models.SidMapping{Sid: dragonUtopiaSid}, 2)

	// Assert
	require.True(t, ok)
	assert.Equal(t, []data.Tuple[int, string]{data.NewTuple(2, "Large Guard")}, mapping.Variants)
}

func TestWhenVariantIdIsUnknown_ReturnsNotOk(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewVariantMappingCatalog()

	// Act
	_, ok := catalog.GetVariantForContentByID(models.SidMapping{Sid: dragonUtopiaSid}, 99)

	// Assert
	assert.False(t, ok)
}

func TestWhenContentHasNoVariants_ReturnsNotOk(t *testing.T) {
	t.Parallel()
	// Arrange
	catalog := content_rules.NewVariantMappingCatalog()

	// Act
	_, ok := catalog.GetVariantForContentByID(models.SidMapping{Sid: watchtowerSid}, 0)

	// Assert
	assert.False(t, ok)
}
