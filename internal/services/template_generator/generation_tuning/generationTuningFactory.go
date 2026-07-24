package generation_tuning

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/utils"
)

type GenerationTuningFactory struct{}

func NewGenerationTuningFactory() *GenerationTuningFactory {
	return &GenerationTuningFactory{}
}

func (this *GenerationTuningFactory) Create(
	configuration *config.GeneratorConfig,
	totalZoneCount int,
) models.GenerationTuning {
	remoteFootholdCount := 0
	if configuration.SpawnRemoteFootholds {
		remoteFootholdCount = configuration.RemoteFootholdCount
	}
	abandonedOutpostCount := 0
	if configuration.ZoneConfiguration.SpawnAbandonedOutposts {
		abandonedOutpostCount = configuration.ZoneConfiguration.AbandonedOutpostCount
	}
	return models.GenerationTuning{
		ContentScale:                   utils.ComputeContentScale(configuration.MapSize, totalZoneCount),
		ResourceDensityMultiplier:      float64(configuration.ZoneConfiguration.ResourceDensityPercent) / 200.0,
		StructureDensityMultiplier:     float64(configuration.ZoneConfiguration.StructureDensityPercent) / 100.0,
		NeutralStackStrengthMultiplier: float64(configuration.ZoneConfiguration.NeutralStackStrengthPercent) / 100.0,
		BorderGuardStrengthMultiplier:  float64(configuration.ZoneConfiguration.BorderGuardStrengthPercent) / 100.0,
		GuardRandomization:             configuration.ZoneConfiguration.GetEffectiveGuardRandomization(),
		RemoteFootholdCount:            remoteFootholdCount,
		AbandonedOutpostCount:          abandonedOutpostCount,
		PlayerOwnedCastles:             configuration.ZoneConfiguration.PlayerOwnedCastles,
	}
}
