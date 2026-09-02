package template_variant_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
)

type Zone struct {
	Name string

	// Quality is the tier the generator planned for this zone, or the tier the
	// user stamped on it in the manual editor. A nil pointer means "not
	// recorded, infer it" - and it has to be a pointer, because the Quality
	// enum counts from iota - 1, so a plain field would read back as
	// QualityLowest and silently down-tier every zone nobody set.
	Quality *neutral_zone.Quality

	// GeneratorPosition is the normalized [0,1] hint the position-driven
	// topologies stamp so the preview reproduces the generated geometry.
	GeneratorPosition *[2]float64

	// GeneratorRing is the concentric-ring index stamped for Circles layouts.
	GeneratorRing *int

	// ManualPosition is the normalized [0,1] position assigned when the user
	// moves or adds a zone in the manual editor.
	ManualPosition *[2]float64

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

// ToZoneModel lifts a persisted zone with no recorded tier. Callers that know
// the tier stamp it afterwards; nothing here can invent one.
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

// ToZoneEntity drops the tier: the .rmg.json schema has nowhere to put it, and
// it is persisted with the editor state instead.
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
