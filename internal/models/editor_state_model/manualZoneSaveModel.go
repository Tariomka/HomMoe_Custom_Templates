package editor_state_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

// ManualZoneSaveModel adds behaviour to the behaviour-free
// editor_state.ManualZoneSave entity.
type ManualZoneSaveModel struct {
	editor_state.ManualZoneSave
}

// ToManualZoneSaves converts live editor zones into their serializable form,
// preserving each zone's ManualPosition outside the entities.Zone JSON.
func ToManualZoneSaves(zones []entities.Zone) []editor_state.ManualZoneSave {
	if len(zones) == 0 {
		return nil
	}

	saves := make([]editor_state.ManualZoneSave, 0, len(zones))
	for _, zone := range zones {
		saves = append(saves, editor_state.ManualZoneSave{Zone: zone, ManualPosition: zone.ManualPosition})
	}
	return saves
}

// FromManualZoneSaves rebuilds live editor zones from their serialized form,
// restoring each zone's ManualPosition.
func FromManualZoneSaves(saves []editor_state.ManualZoneSave) []entities.Zone {
	if len(saves) == 0 {
		return nil
	}

	zones := make([]entities.Zone, 0, len(saves))
	for _, save := range saves {
		zone := save.Zone
		zone.ManualPosition = save.ManualPosition
		zones = append(zones, zone)
	}
	return zones
}

// Clone returns a copy that shares no backing array or pointer with the
// receiver. entities.Zone lives in the protected template tree and therefore
// carries no Clone of its own, so every one of its reference-typed fields is
// copied here; a field added there must be added to cloneZone as well.
func (this ManualZoneSaveModel) Clone() ManualZoneSaveModel {
	return ManualZoneSaveModel{
		Zone:           cloneZone(this.Zone),
		ManualPosition: helpers.ClonePointer(this.ManualPosition),
	}
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
