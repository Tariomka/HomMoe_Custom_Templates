// Package composition holds the compile-time dependency graph of the
// application. Every collaborator is declared exactly once here, so the
// generated injector is the single place where the object graph is built.
//
// Provider sets are package-level vars rather than structs, so the
// one-struct-per-file rule of AGENTS.md 4.1 does not apply to this file.
package composition

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	zone_services "github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/goforj/wire"
)

// ZoneSet builds the zone factories every generation path shares.
var ZoneSet = wire.NewSet(
	zone_services.NewCastleFactory,
	zone_services.NewRoadFactory,
	zone_services.NewZoneFactory,
	zone_services.NewZoneClassifier,
	zone_services.NewZoneLabelProvider,
)

// GenerationSet builds the template generator and the providers it delegates to.
var GenerationSet = wire.NewSet(
	config.NewGeneratorConfig,
	generation_tuning.NewGenerationTuningFactory,
	providers.NewContentLimitProvider,
	providers.NewGameRulesProvider,
	providers.NewMandatoryContentProvider,
	providers.NewTopologyProvider,
	providers.NewZoneLayoutProvider,
	template_generator.NewTemplateGenerator,
)

// EditorSet builds the services backing the manual zone editor and the previews.
var EditorSet = wire.NewSet(
	connection_editor.NewConnectionEditorService,
	connection_editor.NewManualReapplyService,
	connection_editor.NewZoneEditorService,
	content_rules.NewContentRuleService,
	preview_service.NewPreviewLayoutService,
	providePreviewGenerator,
)

// InfrastructureSet builds the persistence and mapping collaborators.
var InfrastructureSet = wire.NewSet(
	file_service.NewFileService,
	mappers.NewConfigMapper,
	validators.NewEditorStateValidator,
)

// HandlerSet builds the five handlers the GUI facade delegates to.
var HandlerSet = wire.NewSet(
	handlers.NewContentRuleHandler,
	handlers.NewPreviewHandler,
	handlers.NewStateHandler,
	handlers.NewTemplateHandler,
	handlers.NewZoneEditorHandler,
)

// GuiHandlerSet is the complete graph behind InitializeGuiHandler.
var GuiHandlerSet = wire.NewSet(
	ZoneSet,
	GenerationSet,
	EditorSet,
	InfrastructureSet,
	HandlerSet,
	handlers.NewGuiHandler,
)
