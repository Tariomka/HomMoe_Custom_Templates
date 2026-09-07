package drivers

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

// ApplyEditedZones writes zones and connections edited in the manual zone
// editor back into the live template and stores them in the editor state as
// the authoritative manual snapshot, reapplied on later regenerations and
// saved with the rest of the .gen.json state.
//
// An apply that follows an untouched revert to base stores no snapshot: the
// base zones would otherwise be pinned and reapplied over every later
// regeneration, which is the very thing the revert undoes.
func (this *State) ApplyEditedZones(request dtos.ZoneEditorZonesDto) {
	pendingBase := this.pendingBaseZones
	this.pendingBaseZones = dtos.ZoneEditorZonesDto{}
	if !this.hasTemplateVariants() {
		return
	}

	this.handleUpdateTemplate(request.Zones, request.Connections)
	if request.RevertToBase && matchesZoneSet(request, pendingBase) {
		this.innerState.ClearManualEdits()
		return
	}

	this.innerState.SetManualEdits(request.Zones, request.Connections)
}

// PreviewBaseZones generates a manual-edit-free layout and returns it for an
// open zone editor to display. It commits NOTHING - neither the live template
// nor the stored manual edits change until the user applies - so cancelling
// the editor leaves the edited template exactly as it was.
//
// It reports false when generation produced nothing, with the reason in the
// status line. Regeneration is random, so this is a NEW base layout rather
// than the one the manual edits were originally made on; that layout is not
// retained anywhere.
func (this *State) PreviewBaseZones() (dtos.ZoneEditorZonesDto, bool) {
	dto, err := this.handler.GenerateTemplate(this.GetStateDto())
	if err != nil {
		this.SetStatus(fmt.Sprintf("Generation failed: %v.", err), true)
		return dtos.ZoneEditorZonesDto{}, false
	}
	if dto.Template == nil || len(dto.Template.Variants) == 0 {
		return dtos.ZoneEditorZonesDto{}, false
	}

	variant := dto.Template.Variants[0]
	this.pendingBaseZones = dtos.ZoneEditorZonesDto{
		Zones:       variant.Zones,
		Connections: variant.Connections,
	}

	return this.pendingBaseZones, true
}

func matchesZoneSet(left, right dtos.ZoneEditorZonesDto) bool {
	return reflect.DeepEqual(left.Zones, right.Zones) &&
		reflect.DeepEqual(left.Connections, right.Connections)
}

func (this *State) handleUpdateTemplate(zones []template_model.Zone, connections []template_model.Connection) {
	dto, err := this.handler.UpdateTemplate(dtos.TemplateUpdateDto{
		Template:    this.lastTemplate,
		Zones:       zones,
		Connections: connections,
		EditorState: new(this.GetStateDto()),
	})

	if err != nil && errors.Is(err, common_errors.ErrProvidedTemplateInvalid) {
		this.SetStatus(
			fmt.Sprintf("Unable to update template, possibly because template was not generated. ‼ Error: %v", err),
			true)
		return
	}

	this.setLastTemplate(dto.Template)
	if err != nil {
		this.SetStatus(
			fmt.Sprintf(
				"Applied %d zones and %d connections. ‼ Error: %v; fix before export.",
				len(zones), len(connections), err),
			true)
		return
	}

	this.SetStatus(
		fmt.Sprintf("Applied %d zones and %d connections from the editor.", len(zones), len(connections)),
		false)
}

// reapplyManualEdits restores the manual zone/connection snapshot over the
// freshly generated template. When castle-count options changed since the
// last generation - the only generator options that override manual edits -
// the new counts are first pushed into the snapshot and the updated snapshot
// is stored back so later regenerations and saves carry it.
func (this *State) reapplyManualEdits(castleChanges editor_state_model.CastleSettingChanges) {
	zones := this.innerState.GetManualZones()
	connections := this.innerState.GetManualConnections()
	if castleChanges.Any() {
		zones = this.handler.ReapplyCastleSettings(dtos.CastleSettingsReapplyRequestDto{
			Zones:       zones,
			Changes:     castleChanges,
			EditorState: this.GetStateDto(),
		})
		this.innerState.SetManualEdits(zones, connections)
	}
	this.handleUpdateTemplate(zones, connections)
}
