package generatorConfigMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapperIsCreated_ReturnsNonNilInstance(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	mapper := mappers.NewConfigMapper(mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService()))

	// Assert
	assert.NotNil(t, mapper)
}
