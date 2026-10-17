package editor_state_model

import (
	"slices"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/neutral_zone"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

const (
	roadsFieldName       = "roads"
	mainObjectsFieldName = "mainObjects"
)

type ManualZoneEdit struct {
	editor_state.ManualZoneEdit
}

func NewManualZoneEdit(derived, edited template_model.Zone, alwaysStoreRoads bool) ManualZoneEdit {
	editedEntity := template_model.ToZoneEntity(edited.Clone())
	overrides, unset := DiffFieldOverrides[template_entity.Zone, editor_state.ZoneFieldOverrides](
		template_model.ToZoneEntity(derived), editedEntity)
	if alwaysStoreRoads {
		overrides.Roads = editedEntity.Roads
		unset = slices.DeleteFunc(unset, func(name string) bool { return name == roadsFieldName })
		if editedEntity.Roads == nil {
			unset = append(unset, roadsFieldName)
		}
	}

	var edit ManualZoneEdit
	edit.Name = edited.Name
	edit.Quality = toQualityOrdinal(edited.Quality)
	edit.GeneratorPosition = helpers.ClonePointer(edited.GeneratorPosition)
	edit.GeneratorRing = helpers.ClonePointer(edited.GeneratorRing)
	edit.ManualPosition = helpers.ClonePointer(edited.ManualPosition)
	edit.FieldOverrides = overrides
	edit.Unset = unset
	return edit
}

// Apply rebuilds the edited zone on top of derived.
func (this ManualZoneEdit) Apply(derived template_model.Zone) template_model.Zone {
	zone := template_model.ToZoneModel(ApplyFieldOverrides(
		template_model.ToZoneEntity(derived), this.FieldOverrides, this.Unset)).Clone()
	zone.Name = this.Name
	zone.Quality = fromQualityOrdinal(this.Quality)
	zone.GeneratorPosition = helpers.ClonePointer(this.GeneratorPosition)
	zone.GeneratorRing = helpers.ClonePointer(this.GeneratorRing)
	zone.ManualPosition = helpers.ClonePointer(this.ManualPosition)
	return zone
}

// IsAdded reports whether the zone was made in the editor rather than by the generator.
func (this ManualZoneEdit) IsAdded() bool { return this.Origin == editor_state.ZoneOriginAdded }

// GetCastles returns the recorded castle count of an added or re-profiled zone.
func (this ManualZoneEdit) GetCastles() (int, bool) {
	if this.Castles == nil {
		return 0, false
	}

	return *this.Castles, true
}

// GetQuality returns the recorded tier, or nil when none was recorded.
func (this ManualZoneEdit) GetQuality() *neutral_zone.Quality {
	return fromQualityOrdinal(this.Quality)
}

func (this ManualZoneEdit) HasMainObjectsOverride() bool {
	return this.FieldOverrides.MainObjects != nil || slices.Contains(this.Unset, mainObjectsFieldName)
}

// AsAdded returns the edit recorded as editor-made with the castle count it
// derives from.
func (this ManualZoneEdit) AsAdded(castles int) ManualZoneEdit {
	this.Origin = editor_state.ZoneOriginAdded
	this.Castles = &castles
	return this
}

// AsReprofiled returns the edit recorded as a generated zone whose profile
// was changed, with the castle count it derives from.
func (this ManualZoneEdit) AsReprofiled(castles int) ManualZoneEdit {
	this.Origin = ""
	this.Castles = &castles
	return this
}

func (this ManualZoneEdit) Clone() ManualZoneEdit {
	clone := this
	clone.Quality = helpers.ClonePointer(this.Quality)
	clone.GeneratorPosition = helpers.ClonePointer(this.GeneratorPosition)
	clone.GeneratorRing = helpers.ClonePointer(this.GeneratorRing)
	clone.ManualPosition = helpers.ClonePointer(this.ManualPosition)
	clone.Castles = helpers.ClonePointer(this.Castles)
	clone.Unset = slices.Clone(this.Unset)
	clone.FieldOverrides = cloneFieldOverrides(this.FieldOverrides)
	return clone
}

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
