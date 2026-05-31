package generator

import (
	"strings"

	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/topology"
)

type TemplateGenerator struct {
	configuration     *config.GeneratorConfig
	zoneLabelProvider *ZoneLabelProvider
}

func NewTemplateGenerator(configuration *config.GeneratorConfig) *TemplateGenerator {
	if configuration == nil {
		configuration = config.NewGeneratorConfig()
	}
	return &TemplateGenerator{
		configuration:     configuration,
		zoneLabelProvider: NewZoneLabelProvider(),
	}
}

func (this *TemplateGenerator) SetConfiguration(configuration *config.GeneratorConfig) {
	if configuration != nil {
		this.configuration = configuration
	}
}

func (this *TemplateGenerator) Generate() *template.RmgTemplateModel {
	this.configuration.EnsureNameExists()
	playerLabels := this.zoneLabelProvider.CreatePlayerLabels(this.configuration.PlayerCount)
	neutralZones := this.zoneLabelProvider.CreateNeutralZonePlans(*this.configuration)
	holdCityLabel := this.zoneLabelProvider.GetHoldCityLabel(*this.configuration, playerLabels, neutralZones)
	tuning := this.createGenerationTuning(this.configuration.PlayerCount + len(neutralZones))
	victoryCondition := this.configuration.GetVictoryCondition()

	return &template.RmgTemplateModel{
		Name:                this.configuration.TemplateName,
		GameMode:            this.configuration.GameMode,
		Description:         this.createTemplateDescription(len(neutralZones)),
		DisplayWinCondition: victoryCondition,
		SizeX:               this.configuration.MapSize,
		SizeZ:               this.configuration.MapSize,
		GameRules:           buildGameRules(this.configuration, victoryCondition),
		Variants: []template.Variant{
			topology.NewTopologyFactory().
				ShufflePlayerZones(true).
				CreateTopologyVariant(this.configuration, playerLabels, neutralZones, tuning, holdCityLabel)},
		ZoneLayouts:        buildZoneLayouts(),
		MandatoryContent:   buildAllMandatoryContent(playerLabels, neutralZones, this.configuration),
		ContentCountLimits: BuildAllContentCountLimits(this.configuration),
		ContentPools:       []template.ContentPool{},
		ContentLists:       []template.ContentList{},
	}
}

func (this *TemplateGenerator) createGenerationTuning(totalZoneCount int) models.GenerationTuning {
	return models.GenerationTuning{
		ContentScale:                   computeContentScale(this.configuration.MapSize, totalZoneCount),
		ResourceDensityMultiplier:      float64(this.configuration.ZoneConfiguration.ResourceDensityPercent) / 200.0,
		StructureDensityMultiplier:     float64(this.configuration.ZoneConfiguration.StructureDensityPercent) / 100.0,
		NeutralStackStrengthMultiplier: float64(this.configuration.ZoneConfiguration.NeutralStackStrengthPercent) / 100.0,
		BorderGuardStrengthMultiplier:  float64(this.configuration.ZoneConfiguration.BorderGuardStrengthPercent) / 100.0,
		GuardRandomization:             effectiveGuardRandomization(this.configuration),
	}
}

func (this *TemplateGenerator) createTemplateDescription(neutralCount int) string {
	parts := []string{
		constants.GetTopologyDescriptor(this.configuration.Topology).Label + " layout",
		countPhrase(neutralCount, "neutral zone", "neutral zones"),
		countPhrase(this.configuration.ZoneConfiguration.PlayerZoneCastles, "castle", "castles") + " per player zone",
	}
	if neutralCount > 0 {
		if this.configuration.ZoneConfiguration.Advanced.Enabled {
			parts = append(parts, "mixed neutral zone tiers")
		} else {
			parts = append(parts, countPhrase(this.configuration.ZoneConfiguration.NeutralZoneCastles, "castle", "castles")+" per neutral zone")
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
	return "Generated with Olden Era Template Generator: " + strings.Join(parts, ", ") + "."
}
