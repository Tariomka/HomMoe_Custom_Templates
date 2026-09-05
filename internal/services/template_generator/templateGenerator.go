package template_generator

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/generation_tuning"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers/provider_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones/zone_interfaces"
)

type TemplateGenerator struct {
	configuration     *config.GeneratorConfig
	zoneLabelProvider zone_interfaces.IZoneLabelProvider
	tuningFactory     generation_tuning.IGenerationTuningFactory

	contentLimitProvider provider_interfaces.IContentLimitProvider
	contentProvider      provider_interfaces.IMandatoryContentProvider
	gameRulesProvider    provider_interfaces.IGameRulesProvider
	gladiatorProvider    provider_interfaces.IGladiatorArenaProvider
	topologyProvider     provider_interfaces.ITopologyProvider
	zoneLayoutProvider   provider_interfaces.IZoneLayoutProvider
}

func NewTemplateGenerator(
	configuration *config.GeneratorConfig,
	zoneLabelProvider zone_interfaces.IZoneLabelProvider,
	tuningFactory generation_tuning.IGenerationTuningFactory,
	contentLimitProvider provider_interfaces.IContentLimitProvider,
	contentProvider provider_interfaces.IMandatoryContentProvider,
	gameRulesProvider provider_interfaces.IGameRulesProvider,
	gladiatorProvider provider_interfaces.IGladiatorArenaProvider,
	topologyProvider provider_interfaces.ITopologyProvider,
	zoneLayoutProvider provider_interfaces.IZoneLayoutProvider) ITemplateGenerator {
	return &TemplateGenerator{
		configuration:        configuration,
		zoneLabelProvider:    zoneLabelProvider,
		tuningFactory:        tuningFactory,
		contentLimitProvider: contentLimitProvider,
		contentProvider:      contentProvider,
		gameRulesProvider:    gameRulesProvider,
		gladiatorProvider:    gladiatorProvider,
		topologyProvider:     topologyProvider,
		zoneLayoutProvider:   zoneLayoutProvider,
	}
}

func (this *TemplateGenerator) SetConfiguration(configuration *config.GeneratorConfig) {
	if configuration != nil {
		this.configuration = configuration
	}
}

// Generate builds the template for the current configuration and returns it,
// with the tier the zone factory planned recorded on every neutral and hub
// zone, plus the warnings raised while parsing the configuration's free-text
// fields.
func (this *TemplateGenerator) Generate() (*template_model.Template, []string) {
	this.configuration.EnsureNameExists()
	playerLabels := this.zoneLabelProvider.CreatePlayerLabels(this.configuration.PlayerCount)
	neutralZones := this.zoneLabelProvider.CreateNeutralZonePlans(*this.configuration)
	holdCityLabel := this.zoneLabelProvider.GetHoldCityLabel(*this.configuration, playerLabels, neutralZones)
	tuning := this.tuningFactory.Create(this.configuration, this.configuration.PlayerCount+len(neutralZones))
	valueOverrides, warnings := this.gameRulesProvider.CreateValueOverrides(*this.configuration)

	variant := this.topologyProvider.
		CreateTopologyVariant(*this.configuration, playerLabels, neutralZones, tuning, holdCityLabel)

	generated := template_model.Template{
		Name:                this.configuration.TemplateName,
		GameMode:            this.configuration.GameMode,
		Description:         this.createTemplateDescription(len(neutralZones)),
		DisplayWinCondition: this.configuration.GetVictoryCondition(),
		SizeX:               this.configuration.MapSize,
		SizeZ:               this.configuration.MapSize,
		ValueOverrides:      valueOverrides,
		GlobalBans:          this.gameRulesProvider.CreateGlobalBans(*this.configuration),
		GameRules:           this.gameRulesProvider.CreateGameRules(*this.configuration),
		Variants:            []template_model.Variant{variant},
		ZoneLayouts:         this.zoneLayoutProvider.CreateZoneLayouts(),
		MandatoryContent:    this.contentProvider.CreateContents(*this.configuration, playerLabels, neutralZones),
		ContentCountLimits:  this.contentLimitProvider.CreateContentCountLimits(*this.configuration),
		ContentPools:        []template_model.ContentPool{},
		ContentLists:        []template_model.ContentList{},
	}

	this.gladiatorProvider.PlaceArena(*this.configuration, &generated.Variants[0])

	return &generated, warnings
}

func (this *TemplateGenerator) createTemplateDescription(neutralCount int) string {
	parts := []string{
		common_topologies.GetTopologyDescriptorFromType(this.configuration.Topology).Label + " layout",
		formatPhraseWithCount(neutralCount, "neutral zone", "neutral zones"),
		formatPhraseWithCount(
			1+this.configuration.ZoneConfiguration.PlayerZoneCastles+
				this.configuration.ZoneConfiguration.PlayerOwnedCastles,
			"castle", "castles") + " per player zone",
	}
	if neutralCount > 0 {
		if this.configuration.ZoneConfiguration.Advanced.Enabled {
			parts = append(parts, "mixed neutral zone tiers")
		} else {
			parts = append(parts,
				formatPhraseWithCount(
					this.configuration.ZoneConfiguration.NeutralZoneCastles,
					"castle",
					"castles")+" per neutral zone")
		}
	}
	var options []string
	if this.configuration.NoDirectPlayerConnections {
		options = append(options, "isolated player starts")
	}
	if this.configuration.RandomPortals {
		options = append(options, "random portals")
	}
	if !this.configuration.SpawnRemoteFootholds {
		options = append(options, "no remote footholds")
	}
	if !this.configuration.GenerateRoads {
		options = append(options, "roads disabled")
	}
	if len(options) > 0 {
		parts = append(parts, "options: "+strings.Join(options, ", "))
	}
	return "Generated with Custom Template Editor: " + strings.Join(parts, ", ") + "."
}

func formatPhraseWithCount(count int, singular, plural string) string {
	if count == 0 {
		return "no " + plural
	}
	word := singular
	if count != 1 {
		word = plural
	}
	return fmt.Sprintf("%d %s", count, word)
}
