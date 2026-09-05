// Package templateGenerator_test contains shared arrangement helpers for the
// templateGenerator.go unit tests.
package templateGenerator_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
)

// generateTemplate flattens the generated template back to the .rmg.json shape
// for the tests that assert on it; the tests that assert on the planned tiers
// call Generate directly. Going through the real mapper here also makes every
// one of those assertions a proof that the round trip is lossless.
func generateTemplate(
	generator template_generator.ITemplateGenerator) (*entities.RmgTemplate, []string) {
	generated, warnings := generator.Generate()
	return new(mappers.NewTemplateMapper().ToEntity(*generated)), warnings
}
