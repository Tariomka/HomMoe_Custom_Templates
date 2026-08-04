// Package templateGenerator_test contains shared arrangement helpers for the
// templateGenerator.go unit tests.
package templateGenerator_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// newTemplateGenerator builds the generator with the same collaborators the
// application wires.
func newTemplateGenerator(configuration *config.GeneratorConfig) *template_generator.TemplateGenerator {
	castleFactory := zones.NewCastleFactory()
	roadFactory := zones.NewRoadFactory()
	return template_generator.NewTemplateGenerator(
		configuration,
		castleFactory,
		roadFactory,
		zones.NewZoneFactory(castleFactory, roadFactory))
}
