package generator

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template"
)

type TemplateGenerator struct {
	configuration *config.GeneratorConfig
}

func NewTemplateGenerator(configuration *config.GeneratorConfig) *TemplateGenerator {
	if configuration == nil {
		configuration = config.NewGeneratorConfig()
	}
	return &TemplateGenerator{configuration: configuration}
}

func (this *TemplateGenerator) SetConfiguration(configuration *config.GeneratorConfig) {
	if configuration != nil {
		this.configuration = configuration
	}
}

func (this *TemplateGenerator) Generate() *template.RmgTemplateModel {
	this.ensureTemplateName()
	playerLetters := this.getPlayerLetters()

	neutralZones := buildNeutralZonePlan(this.configuration)
	var holdCityNeutralLetter string
	if this.configuration.IsHubCityToHold() {
		adj := buildTopologyAdjacency(this.configuration, playerLetters, neutralZones)
		holdCityNeutralLetter = pickHoldCityNeutralLetter(neutralZones, playerLetters, adj)
	}

	totalZones := this.configuration.PlayerCount + len(neutralZones)
	tuning := generationTuning{
		ContentScale:                   computeContentScale(this.configuration.MapSize, totalZones),
		ResourceDensityMultiplier:      float64(this.configuration.ZoneConfiguration.ResourceDensityPercent) / 200.0,
		StructureDensityMultiplier:     float64(this.configuration.ZoneConfiguration.StructureDensityPercent) / 100.0,
		NeutralStackStrengthMultiplier: float64(this.configuration.ZoneConfiguration.NeutralStackStrengthPercent) / 100.0,
		BorderGuardStrengthMultiplier:  float64(this.configuration.ZoneConfiguration.BorderGuardStrengthPercent) / 100.0,
		GuardRandomization:             effectiveGuardRandomization(this.configuration),
	}

	victoryCondition := this.configuration.GetVictoryCondition()

	template := &template.RmgTemplateModel{
		Name:                this.configuration.TemplateName,
		GameMode:            this.configuration.GameMode,
		Description:         buildTemplateDescription(this.configuration, len(neutralZones)),
		DisplayWinCondition: victoryCondition,
		SizeX:               this.configuration.MapSize,
		SizeZ:               this.configuration.MapSize,
		GameRules:           buildGameRules(this.configuration, victoryCondition),
		Variants:            []template.Variant{buildVariant(this.configuration, playerLetters, neutralZones, tuning, holdCityNeutralLetter, this.configuration.IsHubCityToHold())},
		ZoneLayouts:         buildZoneLayouts(),
		MandatoryContent:    buildAllMandatoryContent(playerLetters, neutralZones, this.configuration),
		ContentCountLimits:  BuildAllContentCountLimits(this.configuration),
		ContentPools:        []template.ContentPool{},
		ContentLists:        []template.ContentList{},
	}
	return template
}

func (this *TemplateGenerator) ensureTemplateName() {
	if this.configuration.TemplateName == "" {
		this.configuration.TemplateName = constants.DefaultTemplateName
	}
}

func (this *TemplateGenerator) getPlayerLetters() []string {
	letters := make([]string, this.configuration.PlayerCount)
	copy(letters, ZoneLetters[:this.configuration.PlayerCount])
	return letters
}
