package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ToManualZoneSaveEntities converts live editor zones into their serializable
// form. The positions and the recorded tier travel beside the zone because
// entities.Zone has nowhere to put either.
func ToManualZoneSaveEntities(zones []template_model.Zone) []editor_state.ManualZoneSave {
	if len(zones) == 0 {
		return nil
	}

	return linq.FromSlice(zones).
		Select(func(zone template_model.Zone) editor_state.ManualZoneSave {
			return editor_state.ManualZoneSave{
				Zone:              template_model.ToZoneEntity(zone),
				GeneratorPosition: helpers.ClonePointer(zone.GeneratorPosition),
				GeneratorRing:     helpers.ClonePointer(zone.GeneratorRing),
				ManualPosition:    helpers.ClonePointer(zone.ManualPosition),
				Quality:           toQualityOrdinal(zone.Quality)}
		}).ToSlice()
}

// ToManualZoneModels rebuilds live editor zones from their serialized form. A
// save written before the tier was persisted has none, and the zone falls back
// to inference; a save written before schema v2 has no generator stamps, and
// the preview falls back to a computed layout.
func ToManualZoneModels(saves []editor_state.ManualZoneSave) []template_model.Zone {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save editor_state.ManualZoneSave) template_model.Zone {
			zone := template_model.ToZoneModel(save.Zone)
			zone.GeneratorPosition = helpers.ClonePointer(save.GeneratorPosition)
			zone.GeneratorRing = helpers.ClonePointer(save.GeneratorRing)
			zone.ManualPosition = helpers.ClonePointer(save.ManualPosition)
			zone.Quality = fromQualityOrdinal(save.Quality)
			return zone
		}).ToSlice()
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
