// Package templateGenerator_test contains shared arrangement helpers for the
// templateGenerator.go unit tests.
package templateGenerator_test

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// newTemplateGenerator builds the generator with the same collaborators the
// application wires.
func newTemplateGenerator(configuration *config.GeneratorConfig) *template_generator.TemplateGenerator {
	castleFactory := zones.NewCastleFactory()
	roadFactory := zones.NewRoadFactory()
	zoneFactory := zones.NewZoneFactory(castleFactory, roadFactory)
	zoneClassifier := zones.NewZoneClassifier()
	zoneEditor := connection_editor.NewZoneEditorService(castleFactory, roadFactory, zoneFactory)

	return template_generator.NewTemplateGenerator(
		configuration,
		zones.NewZoneLabelProvider(),
		generation_tuning.NewGenerationTuningFactory(),
		providers.NewContentLimitProvider(),
		providers.NewMandatoryContentProvider(zoneClassifier, zoneEditor),
		providers.NewGameRulesProvider(),
		providers.NewTopologyProvider(zoneFactory, roadFactory),
		providers.NewZoneLayoutProvider())
}
