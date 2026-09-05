package editor_state_model

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/linq"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ToManualZoneSaveEntities converts live editor zones into their serializable
// form. The manual position and the recorded tier travel beside the zone
// because entities.Zone has nowhere to put either.
func ToManualZoneSaveEntities(zones []template_model.Zone) []editor_state.ManualZoneSave {
	if len(zones) == 0 {
		return nil
	}

	return linq.FromSlice(zones).
		Select(func(zone template_model.Zone) editor_state.ManualZoneSave {
			return editor_state.ManualZoneSave{
				Zone:           template_model.ToZoneEntity(zone),
				ManualPosition: toPositionArray(zone.ManualPosition),
				Quality:        toQualityOrdinal(zone.Quality)}
		}).ToSlice()
}

// ToManualZoneModels rebuilds live editor zones from their serialized form. A
// save written before the tier was persisted has none, and the zone falls back
// to inference.
func ToManualZoneModels(saves []editor_state.ManualZoneSave) []template_model.Zone {
	if len(saves) == 0 {
		return nil
	}

	return linq.FromSlice(saves).
		Select(func(save editor_state.ManualZoneSave) template_model.Zone {
			zone := template_model.ToZoneModel(save.Zone)
			zone.ManualPosition = fromPositionArray(save.ManualPosition)
			zone.Quality = fromQualityOrdinal(save.Quality)
			return zone
		}).ToSlice()
}

// toPositionArray and fromPositionArray bridge the persisted [x, y] array to
// the vector the model works in. Batch S changes the persisted shape to match
// and deletes both.
func toPositionArray(position *data.Vec2[float64]) *[2]float64 {
	if position == nil {
		return nil
	}

	return &[2]float64{position.X, position.Y}
}

func fromPositionArray(position *[2]float64) *data.Vec2[float64] {
	if position == nil {
		return nil
	}

	return new(data.NewVec2(position[0], position[1]))
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
