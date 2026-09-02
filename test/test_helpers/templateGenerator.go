package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

// NewTemplateGenerator builds the generator with the same collaborators
// internal/composition wires for production.
func NewTemplateGenerator(configuration *config.GeneratorConfig) template_generator.ITemplateGenerator {
	castleFactory := zones.NewCastleFactory()
	roadFactory := zones.NewRoadFactory()
	zoneFactory := zones.NewZoneFactory(castleFactory, roadFactory)
	zoneClassifier := zones.NewZoneTierService()
	zoneEditor := connection_editor.NewZoneEditorService(castleFactory, roadFactory, zoneFactory)
	zoneLabelProvider := zones.NewZoneLabelProvider()
	connectionService := base.NewTopologyConnectionService(zoneLabelProvider)

	return template_generator.NewTemplateGenerator(
		configuration,
		zoneLabelProvider,
		generation_tuning.NewGenerationTuningFactory(),
		mappers.NewTemplateMapper(),
		providers.NewContentLimitProvider(),
		providers.NewMandatoryContentProvider(zoneClassifier, zoneEditor),
		providers.NewGameRulesProvider(),
		providers.NewGladiatorArenaProvider(zoneClassifier),
		providers.NewTopologyProvider(
			NewTopologyServiceLookup(zoneFactory, roadFactory, zoneLabelProvider, connectionService)),
		providers.NewZoneLayoutProvider())
}
