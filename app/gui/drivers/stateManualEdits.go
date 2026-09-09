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

func (this *State) ApplyEditedZones(request dtos.ZoneEditorZonesDto) {
	pendingBase := this.pendingBaseZones
	this.pendingBaseZones = dtos.ZoneEditorZonesDto{}
	if !this.hasTemplateVariants() {
		return
	}

	revertsToUntouchedBase := request.RevertToBase && matchesZoneSet(request, pendingBase)
	if !this.handleUpdateTemplate(request.Zones, request.Connections) {
		return
	}

	if revertsToUntouchedBase && this.innerState.ClearManualEdits() {
		this.flagAsUnsaved()
		return
	}

	if this.innerState.SetManualEdits(request.Zones, request.Connections) {
		this.flagAsUnsaved()
	}
}

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

func (this *State) handleUpdateTemplate(zones []template_model.Zone, connections []template_model.Connection) bool {
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
		return false
	}

	this.setLastTemplate(dto.Template)
	if err != nil {
		this.SetStatus(
			fmt.Sprintf(
				"Applied %d zones and %d connections. ‼ Error: %v; fix before export.",
				len(zones), len(connections), err),
			true)
		return true
	}

	this.SetStatus(
		fmt.Sprintf("Applied %d zones and %d connections from the editor.", len(zones), len(connections)),
		false)
	return true
}

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
