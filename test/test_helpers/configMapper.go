package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
)

// NewConfigMapper builds the mapper with the same collaborators
// internal/composition wires for production.
func NewConfigMapper() *mappers.GeneratorConfigMapper {
	return mappers.NewConfigMapper(mappers.NewMandatoryContentItemMapper(content_rules.NewContentRuleService()))
}
