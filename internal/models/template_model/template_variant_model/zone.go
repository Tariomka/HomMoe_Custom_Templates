package template_variant_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type Zone struct {
	Name    string
	Quality *neutral_zone.Quality

	GeneratorPosition *[2]float64
	ManualPosition    *[2]float64

	GeneratorRing *int

	Size   float64
	Layout string

	GuardCutoffValue          int
	GuardRandomization        float64
	GuardMultiplier           float64
	GuardWeeklyIncrement      float64
	GuardReactionDistribution []int
	DiplomacyModifier         float64

	EncounterHolesSettings *EncounterHolesSettings

	RandomHireEnableWeeklyUnitIncrement []bool
	RandomHireInitialUnitIncrement      []int

	GuardedContentPool   []string
	UnguardedContentPool []string
	ResourcesContentPool []string

	MandatoryContent   StringList
	ContentCountLimits StringList

	GuardedContentValue          int
	GuardedContentValuePerArea   int
	UnguardedContentValue        int
	UnguardedContentValuePerArea int
	ResourcesValue               int
	ResourcesValuePerArea        int

	MainObjects []MainObject

	ZoneBiome        TypedRef
	ContentBiome     TypedRef
	MetaObjectsBiome TypedRef

	CrossroadsPosition *int
	Roads              []Road
}

func (this Zone) Clone() Zone {
	clone := this
	clone.Quality = helpers.ClonePointer(this.Quality)
	clone.GeneratorPosition = helpers.ClonePointer(this.GeneratorPosition)
	clone.GeneratorRing = helpers.ClonePointer(this.GeneratorRing)
	clone.ManualPosition = helpers.ClonePointer(this.ManualPosition)
	clone.GuardReactionDistribution = slices.Clone(this.GuardReactionDistribution)
	clone.EncounterHolesSettings = helpers.ClonePointer(this.EncounterHolesSettings)
	clone.RandomHireEnableWeeklyUnitIncrement = slices.Clone(this.RandomHireEnableWeeklyUnitIncrement)
	clone.RandomHireInitialUnitIncrement = slices.Clone(this.RandomHireInitialUnitIncrement)
	clone.GuardedContentPool = slices.Clone(this.GuardedContentPool)
	clone.UnguardedContentPool = slices.Clone(this.UnguardedContentPool)
	clone.ResourcesContentPool = slices.Clone(this.ResourcesContentPool)
	clone.MandatoryContent = slices.Clone(this.MandatoryContent)
	clone.ContentCountLimits = slices.Clone(this.ContentCountLimits)
	clone.MainObjects = helpers.MapSlice(this.MainObjects, MainObject.Clone)
	clone.ZoneBiome = this.ZoneBiome.Clone()
	clone.ContentBiome = this.ContentBiome.Clone()
	clone.MetaObjectsBiome = this.MetaObjectsBiome.Clone()
	clone.CrossroadsPosition = helpers.ClonePointer(this.CrossroadsPosition)
	clone.Roads = helpers.MapSlice(this.Roads, Road.Clone)
	return clone
}

func ToZoneModel(entity template.Zone) Zone {
	return Zone{
		Name:                      entity.Name,
		GeneratorPosition:         entity.GeneratorPosition,
		GeneratorRing:             entity.GeneratorRing,
		ManualPosition:            entity.ManualPosition,
		Size:                      entity.Size,
		Layout:                    entity.Layout,
		GuardCutoffValue:          entity.GuardCutoffValue,
		GuardRandomization:        entity.GuardRandomization,
		GuardMultiplier:           entity.GuardMultiplier,
		GuardWeeklyIncrement:      entity.GuardWeeklyIncrement,
		GuardReactionDistribution: entity.GuardReactionDistribution,
		DiplomacyModifier:         entity.DiplomacyModifier,
		EncounterHolesSettings: helpers.MapPointer(
			entity.EncounterHolesSettings,
			ToEncounterHolesSettingsModel,
		),
		RandomHireEnableWeeklyUnitIncrement: entity.RandomHireEnableWeeklyUnitIncrement,
		RandomHireInitialUnitIncrement:      entity.RandomHireInitialUnitIncrement,
		GuardedContentPool:                  entity.GuardedContentPool,
		UnguardedContentPool:                entity.UnguardedContentPool,
		ResourcesContentPool:                entity.ResourcesContentPool,
		MandatoryContent:                    ToStringListModel(entity.MandatoryContent),
		ContentCountLimits:                  ToStringListModel(entity.ContentCountLimits),
		GuardedContentValue:                 entity.GuardedContentValue,
		GuardedContentValuePerArea:          entity.GuardedContentValuePerArea,
		UnguardedContentValue:               entity.UnguardedContentValue,
		UnguardedContentValuePerArea:        entity.UnguardedContentValuePerArea,
		ResourcesValue:                      entity.ResourcesValue,
		ResourcesValuePerArea:               entity.ResourcesValuePerArea,
		MainObjects:                         ToMainObjectModels(entity.MainObjects),
		ZoneBiome:                           ToTypedRefModel(entity.ZoneBiome),
		ContentBiome:                        ToTypedRefModel(entity.ContentBiome),
		MetaObjectsBiome:                    ToTypedRefModel(entity.MetaObjectsBiome),
		CrossroadsPosition:                  entity.CrossroadsPosition,
		Roads:                               ToRoadModels(entity.Roads),
	}
}

func ToZoneEntity(model Zone) template.Zone {
	return template.Zone{
		Name:                      model.Name,
		GeneratorPosition:         model.GeneratorPosition,
		GeneratorRing:             model.GeneratorRing,
		ManualPosition:            model.ManualPosition,
		Size:                      model.Size,
		Layout:                    model.Layout,
		GuardCutoffValue:          model.GuardCutoffValue,
		GuardRandomization:        model.GuardRandomization,
		GuardMultiplier:           model.GuardMultiplier,
		GuardWeeklyIncrement:      model.GuardWeeklyIncrement,
		GuardReactionDistribution: model.GuardReactionDistribution,
		DiplomacyModifier:         model.DiplomacyModifier,
		EncounterHolesSettings: helpers.MapPointer(
			model.EncounterHolesSettings,
			ToEncounterHolesSettingsEntity,
		),
		RandomHireEnableWeeklyUnitIncrement: model.RandomHireEnableWeeklyUnitIncrement,
		RandomHireInitialUnitIncrement:      model.RandomHireInitialUnitIncrement,
		GuardedContentPool:                  model.GuardedContentPool,
		UnguardedContentPool:                model.UnguardedContentPool,
		ResourcesContentPool:                model.ResourcesContentPool,
		MandatoryContent:                    ToStringListEntity(model.MandatoryContent),
		ContentCountLimits:                  ToStringListEntity(model.ContentCountLimits),
		GuardedContentValue:                 model.GuardedContentValue,
		GuardedContentValuePerArea:          model.GuardedContentValuePerArea,
		UnguardedContentValue:               model.UnguardedContentValue,
		UnguardedContentValuePerArea:        model.UnguardedContentValuePerArea,
		ResourcesValue:                      model.ResourcesValue,
		ResourcesValuePerArea:               model.ResourcesValuePerArea,
		MainObjects:                         ToMainObjectEntities(model.MainObjects),
		ZoneBiome:                           ToTypedRefEntity(model.ZoneBiome),
		ContentBiome:                        ToTypedRefEntity(model.ContentBiome),
		MetaObjectsBiome:                    ToTypedRefEntity(model.MetaObjectsBiome),
		CrossroadsPosition:                  model.CrossroadsPosition,
		Roads:                               ToRoadEntities(model.Roads),
	}
}

func ToZoneModels(entities []template.Zone) []Zone {
	return helpers.MapSlice(entities, ToZoneModel)
}

func ToZoneEntities(models []Zone) []template.Zone {
	return helpers.MapSlice(models, ToZoneEntity)
}
