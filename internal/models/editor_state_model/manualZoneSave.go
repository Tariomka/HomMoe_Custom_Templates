package editor_state_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ManualZoneSave adds behaviour to the behaviour-free
// editor_state.ManualZoneSave entity.
type ManualZoneSave struct {
	editor_state.ManualZoneSave
}

// ToManualZoneSaves converts live editor zones into their serializable form,
// preserving each zone's ManualPosition outside the entities.Zone JSON and its
// recorded tier, which the .rmg.json schema has nowhere to put.
func ToManualZoneSaves(zones []template_model.Zone) []ManualZoneSave {
	if len(zones) == 0 {
		return nil
	}

	return linq.FromSlice(zones).
		Select(func(zone template_model.Zone) ManualZoneSave {
			return ManualZoneSave{
				Zone:           template_model.ToZoneEntity(zone),
				ManualPosition: zone.ManualPosition,
				Quality:        toQualityOrdinal(zone.Quality)}
		}).ToSlice()
}

// FromManualZoneSaves rebuilds live editor zones from their serialized form,
// restoring each zone's ManualPosition and recorded tier. A save written before
// the tier was persisted has none, and the zone falls back to inference.
func FromManualZoneSaves(saves []ManualZoneSave) []template_model.Zone {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save ManualZoneSave) template_model.Zone {
			zone := template_model.ToZoneModel(save.Zone)
			zone.ManualPosition = save.ManualPosition
			zone.Quality = fromQualityOrdinal(save.Quality)
			return zone
		}).ToSlice()
}

// ToManualZoneSaveModels wraps persisted manual zones for use at the service
// layer. It converts the storage axis; ToManualZoneSaves converts the live
// entities.Zone axis.
func ToManualZoneSaveModels(saves []editor_state.ManualZoneSave) []ManualZoneSave {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save editor_state.ManualZoneSave) ManualZoneSave {
			return ManualZoneSave{ManualZoneSave: save}
		}).ToSlice()
}

// ToManualZoneSaveEntities unwraps manual zones back into their persisted form.
func ToManualZoneSaveEntities(saves []ManualZoneSave) []editor_state.ManualZoneSave {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save ManualZoneSave) editor_state.ManualZoneSave {
			return save.ManualZoneSave
		}).ToSlice()
}

// Clone returns a copy that shares no backing array or pointer with the
// receiver. entities.Zone lives in the protected template tree and therefore
// carries no Clone of its own, so every one of its reference-typed fields is
// copied here; a field added there must be added to cloneZone as well.
func (this ManualZoneSave) Clone() ManualZoneSave {
	return ManualZoneSave{
		Zone:           cloneZone(this.Zone),
		ManualPosition: helpers.ClonePointer(this.ManualPosition),
		Quality:        helpers.ClonePointer(this.Quality)}
}

func cloneZone(source entities.Zone) entities.Zone {
	clone := source

	clone.GeneratorPosition = helpers.ClonePointer(source.GeneratorPosition)
	clone.GeneratorRing = helpers.ClonePointer(source.GeneratorRing)
	clone.ManualPosition = helpers.ClonePointer(source.ManualPosition)
	clone.EncounterHolesSettings = helpers.ClonePointer(source.EncounterHolesSettings)
	clone.CrossroadsPosition = helpers.ClonePointer(source.CrossroadsPosition)

	clone.GuardReactionDistribution = slices.Clone(source.GuardReactionDistribution)
	clone.RandomHireEnableWeeklyUnitIncrement = slices.Clone(source.RandomHireEnableWeeklyUnitIncrement)
	clone.RandomHireInitialUnitIncrement = slices.Clone(source.RandomHireInitialUnitIncrement)
	clone.GuardedContentPool = slices.Clone(source.GuardedContentPool)
	clone.UnguardedContentPool = slices.Clone(source.UnguardedContentPool)
	clone.ResourcesContentPool = slices.Clone(source.ResourcesContentPool)
	clone.MandatoryContent = slices.Clone(source.MandatoryContent)
	clone.ContentCountLimits = slices.Clone(source.ContentCountLimits)

	clone.ZoneBiome = cloneTypedRef(source.ZoneBiome)
	clone.ContentBiome = cloneTypedRef(source.ContentBiome)
	clone.MetaObjectsBiome = cloneTypedRef(source.MetaObjectsBiome)

	clone.MainObjects = slices.Clone(source.MainObjects)
	for objectIndex := range clone.MainObjects {
		clone.MainObjects[objectIndex] = cloneMainObject(clone.MainObjects[objectIndex])
	}

	clone.Roads = slices.Clone(source.Roads)
	for roadIndex := range clone.Roads {
		clone.Roads[roadIndex] = cloneRoad(clone.Roads[roadIndex])
	}

	return clone
}

func cloneMainObject(source entities.MainObject) entities.MainObject {
	clone := source
	clone.Factions = slices.Clone(source.Factions)
	clone.PlacementArgs = slices.Clone(source.PlacementArgs)
	if source.Faction != nil {
		faction := cloneTypedRef(*source.Faction)
		clone.Faction = &faction
	}
	return clone
}

func cloneRoad(source entities.Road) entities.Road {
	clone := source
	clone.From = cloneTypedRef(source.From)
	clone.To = cloneTypedRef(source.To)
	clone.Road = helpers.ClonePointer(source.Road)
	return clone
}

func cloneTypedRef(source entities.TypedRef) entities.TypedRef {
	clone := source
	clone.Args = slices.Clone(source.Args)
	return clone
}

// toQualityOrdinal and fromQualityOrdinal cross the entity boundary the tier
// cannot: an entity may not name neutral_zone.Quality, so it stores the raw
// ordinal. Both directions keep nil meaning "not recorded" - the enum's zero
// value is QualityLowest, so collapsing nil onto it would down-tier silently.
func toQualityOrdinal(quality *neutral_zone.Quality) *int8 {
	if quality == nil {
		return nil
	}

	ordinal := int8(*quality)
	return &ordinal
}

func fromQualityOrdinal(ordinal *int8) *neutral_zone.Quality {
	if ordinal == nil {
		return nil
	}

	quality := neutral_zone.Quality(*ordinal)
	return &quality
}
