package template_generator

import (
	"fmt"
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
)

type TemplateGenerator struct {
	configuration     *config.GeneratorConfig
	zoneLabelProvider *zones.ZoneLabelProvider

	contentLimitProvider *providers.ContentLimitProvider
	contentProvider      *providers.MandatoryContentProvider
	gameRulesProvider    *providers.GameRulesProvider
	topologyProvider     *providers.TopologyProvider
	zoneLayoutProvider   *providers.ZoneLayoutProvider
}

func NewTemplateGenerator(configuration *config.GeneratorConfig) *TemplateGenerator {
	if configuration == nil {
		configuration = config.NewGeneratorConfig()
	}
	return &TemplateGenerator{
		configuration:        configuration,
		zoneLabelProvider:    zones.NewZoneLabelProvider(),
		contentLimitProvider: providers.NewContentLimitProvider(),
		contentProvider:      providers.NewMandatoryContentProvider(),
		gameRulesProvider:    providers.NewGameRulesProvider(),
		topologyProvider:     providers.NewTopologyProvider(),
		zoneLayoutProvider:   providers.NewZoneLayoutProvider(),
	}
}

func (this *TemplateGenerator) SetConfiguration(configuration *config.GeneratorConfig) {
	if configuration != nil {
		this.configuration = configuration
	}
}

func (this *TemplateGenerator) Generate() *template.RmgTemplate {
	this.configuration.EnsureNameExists()
	playerLabels := this.zoneLabelProvider.CreatePlayerLabels(this.configuration.PlayerCount)
	neutralZones := this.zoneLabelProvider.CreateNeutralZonePlans(*this.configuration)
	holdCityLabel := this.zoneLabelProvider.GetHoldCityLabel(*this.configuration, playerLabels, neutralZones)
	tuning := this.createGenerationTuning(this.configuration.PlayerCount + len(neutralZones))

	return &template.RmgTemplate{
		Name:                this.configuration.TemplateName,
		GameMode:            this.configuration.GameMode,
		Description:         this.createTemplateDescription(len(neutralZones)),
		DisplayWinCondition: this.configuration.GetVictoryCondition(),
		SizeX:               this.configuration.MapSize,
		SizeZ:               this.configuration.MapSize,
		ValueOverrides:      this.gameRulesProvider.CreateValueOverrides(*this.configuration),
		GlobalBans:          this.gameRulesProvider.CreateGlobalBans(*this.configuration),
		GameRules:           this.gameRulesProvider.CreateGameRules(*this.configuration),
		Variants: []entities.Variant{
			this.topologyProvider.
				ShufflePlayerZones(this.configuration.ShufflePlayerZones).
				CreateTopologyVariant(*this.configuration, playerLabels, neutralZones, tuning, holdCityLabel),
		},
		ZoneLayouts:        this.zoneLayoutProvider.CreateZoneLayouts(),
		MandatoryContent:   this.contentProvider.CreateContents(*this.configuration, playerLabels, neutralZones),
		ContentCountLimits: this.contentLimitProvider.CreateContentCountLimits(*this.configuration),
		ContentPools:       []entities.ContentPool{},
		ContentLists:       []entities.ContentList{},
	}
}

func (this *TemplateGenerator) createGenerationTuning(totalZoneCount int) models.GenerationTuning {
	return models.NewGenerationTuning(this.configuration, totalZoneCount)
}

func (this *TemplateGenerator) createTemplateDescription(neutralCount int) string {
	parts := []string{
		constants.GetTopologyDescriptorFromType(this.configuration.Topology).Label + " layout",
		formatPhraseWithCount(neutralCount, "neutral zone", "neutral zones"),
		formatPhraseWithCount(
			1+this.configuration.ZoneConfiguration.PlayerZoneCastles+this.configuration.ZoneConfiguration.PlayerOwnedCastles,
			"castle",
			"castles") + " per player zone",
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
