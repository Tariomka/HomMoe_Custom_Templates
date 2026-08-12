// Package composition holds the compile-time dependency graph of the
// application. Every collaborator is declared exactly once here, so the
// generated injector is the single place where the object graph is built.
//
//nolint:gochecknoglobals // Dependency injection
package composition

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/bonuses"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/connection_editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	editor_services "github.com/Tariomka/hommoe_custom_templates/internal/services/editor"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_system"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/pickers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/preview_service"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/topology/base"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zone_content"
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
	base.NewTopologyConnectionService,
)

// GenerationSet builds the template generator and the providers it delegates to.
var GenerationSet = wire.NewSet(
	config.NewGeneratorConfig,
	generation_tuning.NewGenerationTuningFactory,
	providers.NewContentLimitProvider,
	providers.NewGameRulesProvider,
	providers.NewGladiatorArenaProvider,
	providers.NewMandatoryContentProvider,
	providers.NewTopologyProvider,
	providers.NewZoneLayoutProvider,
	provideTopologyServices,
	template_generator.NewTemplateGenerator,
)

// EditorSet builds the services backing the manual zone editor and the previews.
var EditorSet = wire.NewSet(
	bonuses.NewBonusEntryService,
	pickers.NewPickerEntryService,
	connection_editor.NewConnectionEditorService,
	connection_editor.NewManualReapplyService,
	connection_editor.NewZoneEditorService,
	connection_editor.NewZoneEditorGeometryService,
	content_rules.NewContentRuleService,
	preview_service.NewPreviewLayoutService,
	providePreviewGenerator,
	zone_content.NewZoneContentEditorService,
)

// InfrastructureSet builds the persistence and mapping collaborators.
var InfrastructureSet = wire.NewSet(
	repositories.NewEditorStateRepository,
	repositories.NewPreviewRepository,
	repositories.NewTemplateRepository,
	file_service.NewFileService,
	mappers.NewConfigMapper,
	mappers.NewMandatoryContentItemMapper,
	validators.NewEditorStateValidator,
)

// HandlerSet builds the handlers the GUI facade delegates to.
var HandlerSet = wire.NewSet(
	handlers.NewBonusHandler,
	handlers.NewPickerHandler,
	handlers.NewContentRuleHandler,
	handlers.NewPreviewHandler,
	handlers.NewStateHandler,
	handlers.NewTemplateHandler,
	handlers.NewZoneEditorHandler,
	handlers.NewZoneContentHandler,
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

// FileSystemSet is the complete graph behind InitializeFileSystemHandler. It is
// deliberately disjoint from GuiHandlerSet: browsing the local disk shares no
// collaborator with template editing.
var FileSystemSet = wire.NewSet(
	file_system.NewDirectoryBrowserService,
	file_system.NewPathResolutionService,
	handlers.NewFileSystemHandler,
)

// RegenerationSet is the complete graph behind InitializeRegenerationHandler.
// Like FileSystemSet it is disjoint from GuiHandlerSet: deciding *when* to
// regenerate shares no collaborator with actually generating.
var RegenerationSet = wire.NewSet(
	editor_services.NewRegenerationDecisionService,
	handlers.NewRegenerationHandler,
)
