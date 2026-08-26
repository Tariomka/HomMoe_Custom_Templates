package panels

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_topologies"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
)

func (this *LayoutPanel) getManualZoneEditWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Manual zone editing", []layout.Widget{
		widgets.NewBrightButtonLargeWidget(theme, "Manual zone editor...", &this.editConnectionsBtn, false),
		widgets.NewDimmedLabelWidget(theme, "Visually add, move and edit zones and connections on the generated map."),
	})
}

func (this *LayoutPanel) getZonesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Zones", []layout.Widget{
		widgets.NewSliderRowWidget(theme, "Player Owned castles per zone", constants.DefaultLabelWidthLong,
			&this.sldPlayerOwnedCastles, utils.RoundedRangeFormatter(0, 4)),
		widgets.NewSliderRowWidget(theme, "Player Unclaimed castles per zone", constants.DefaultLabelWidthLong,
			&this.sldPlayerCastles, utils.RoundedRangeFormatter(0, 4)),
		widgets.NewBrightButtonLargeWidget(theme, "Edit player zone content...", &this.btnPlayerContent, false),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkAdvancedZones,
			"Advanced zone control (split lowest / low / medium / high tiers)"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkAdvancedZones.Value {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(
						widgets.NewSliderRowWidget(theme, "Total neutral zones", constants.DefaultLabelWidthLong,
							&this.sldNeutralCount, utils.RoundedRangeFormatter(0, 16))),
					layout.Rigid(
						widgets.NewSliderRowWidget(theme, "Neutral castles per zone", constants.DefaultLabelWidthLong,
							&this.sldNeutralCastles, utils.RoundedRangeFormatter(0, 4))))
			}

			return this.getAdvancedZonesWidget(theme)(gtx)
		},
	})
}

func (this *LayoutPanel) getAdvancedZonesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Advanced options", []layout.Widget{
		this.getNeutralTierSectionWidget(theme, "Lowest tier",
			&this.sldNeutralLowestNoCastle, &this.sldNeutralLowestCastle,
			&this.sldNeutralLowestCastlesPerZone, &this.btnLowestContent),
		this.getNeutralTierSectionWidget(theme, "Low tier",
			&this.sldNeutralLowNoCastle, &this.sldNeutralLowCastle,
			&this.sldNeutralLowCastlesPerZone, &this.btnLowContent),
		this.getNeutralTierSectionWidget(theme, "Medium tier",
			&this.sldNeutralMedNoCastle, &this.sldNeutralMedCastle,
			&this.sldNeutralMedCastlesPerZone, &this.btnMedContent),
		this.getNeutralTierSectionWidget(theme, "High tier",
			&this.sldNeutralHighNoCastle, &this.sldNeutralHighCastle,
			&this.sldNeutralHighCastlesPerZone, &this.btnHighContent),
		this.getHubTierSectionWidget(theme),
	})
}

// getNeutralTierSectionWidget renders one advanced neutral-tier sub-section: the
// zone-count sliders, the per-tier castles-per-zone slider, and a button that
// opens the tier's zone-content editor dialog.
func (this *LayoutPanel) getNeutralTierSectionWidget(theme *material.Theme, title string,
	noCastle, withCastle, castlesPerZone *widget.Float, contentBtn *widget.Clickable) layout.Widget {
	return widgets.NewSectionWidget(theme, title, []layout.Widget{
		widgets.NewSliderRowWidget(theme, "No castle", constants.DefaultLabelWidthShort,
			noCastle, utils.RoundedRangeFormatter(0, 8)),
		widgets.NewSliderRowWidget(theme, "With castle", constants.DefaultLabelWidthShort,
			withCastle, utils.RoundedRangeFormatter(0, 8)),
		widgets.NewSliderRowWidget(theme, "Neutral castles per zone", constants.DefaultLabelWidth,
			castlesPerZone, utils.RoundedRangeFormatter(1, 4)),
		widgets.NewBrightButtonLargeWidget(theme, "Edit zone content...", contentBtn, false),
	})
}

// getHubTierSectionWidget renders the advanced Hub sub-section. It only appears
// for the topologies that create a Hub zone and (being nested inside the
// advanced options) only while advanced zone control is enabled.
func (this *LayoutPanel) getHubTierSectionWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if !common_topologies.GetTopologyCapabilities(this.state.GetTopology()).UsesHub {
			return layout.Dimensions{}
		}

		return widgets.NewSectionWidget(theme, "Hub", []layout.Widget{
			widgets.NewSliderRowWidget(theme, "Hub zone castles", constants.DefaultLabelWidth,
				&this.sldHubCastles, utils.RoundedRangeFormatter(0, 4)),
			widgets.NewBrightButtonLargeWidget(theme, "Edit zone content...", &this.btnHubContent, false),
		})(gtx)
	}
}

// handleConnectionEditorClick opens the manual zone editor over the most
// recently generated template, or reports that one must be generated first.
func (this *LayoutPanel) handleConnectionEditorClick(gtx layout.Context) {
	if !this.editConnectionsBtn.Clicked(gtx) {
		return
	}

	lastTemplate := this.state.GetLastTemplate()
	if lastTemplate == nil || len(lastTemplate.Variants) == 0 {
		this.state.SetStatus("Generate a template first to edit its zones.", true)
		return
	}

	activeVariant := lastTemplate.Variants[0]
	options := this.zoneEditorHandler.GetZoneEditorOptions(this.state.GetStateDto(), len(activeVariant.Zones))
	this.state.GetDialogHost().Open(dialogs.NewZoneEditorDialog(
		activeVariant.Zones,
		activeVariant.Connections,
		options.Topology,
		options.Tuning,
		options.GenerateRoads,
		this.zoneEditorHandler,
		this.state.ApplyEditedZones,
		this.state.PreviewBaseZones))
}

// handleZoneContentDialogClicks opens the single-tier zone-content editor for
// whichever per-zone "Edit zone content..." button was clicked this frame.
func (this *LayoutPanel) handleZoneContentDialogClicks(gtx layout.Context) {
	settings := this.state.GetStateData()
	switch {
	case this.btnPlayerContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Player", true, settings.PlayerZoneContentRows,
			func(state *editor_state_model.EditorState, rows []editor_state_model.ZoneContentRow) {
				state.PlayerZoneContentRows = rows
			})
	case this.btnLowestContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Lowest Neutral", false, settings.LowestNeutralContentRows,
			func(state *editor_state_model.EditorState, rows []editor_state_model.ZoneContentRow) {
				state.LowestNeutralContentRows = rows
			})
	case this.btnLowContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Low Neutral", false, settings.LowNeutralContentRows,
			func(state *editor_state_model.EditorState, rows []editor_state_model.ZoneContentRow) {
				state.LowNeutralContentRows = rows
			})
	case this.btnMedContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Medium Neutral", false, settings.MediumNeutralContentRows,
			func(state *editor_state_model.EditorState, rows []editor_state_model.ZoneContentRow) {
				state.MediumNeutralContentRows = rows
			})
	case this.btnHighContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: High Neutral", false, settings.HighNeutralContentRows,
			func(state *editor_state_model.EditorState, rows []editor_state_model.ZoneContentRow) {
				state.HighNeutralContentRows = rows
			})
	case this.btnHubContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Hub", false, settings.HubZoneContentRows,
			func(state *editor_state_model.EditorState, rows []editor_state_model.ZoneContentRow) {
				state.HubZoneContentRows = rows
			})
	}
}

// openZoneContentDialog opens a ZoneContentDialog for a single tier and writes
// the edited rows back into the editor state through the given setter.
func (this *LayoutPanel) openZoneContentDialog(
	title string,
	isPlayerTier bool,
	rows []editor_state_model.ZoneContentRow,
	set func(*editor_state_model.EditorState, []editor_state_model.ZoneContentRow)) {
	this.state.GetDialogHost().Open(dialogs.NewZoneContentDialog(
		title, isPlayerTier, rows, this.contentRuleHandler, this.state.GetDialogHost().Open,
		func(updated []editor_state_model.ZoneContentRow) {
			this.state.UpdateState(func(s *editor_state_model.EditorState) { set(s, updated) })
		}))
}
