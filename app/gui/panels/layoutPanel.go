package panels

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/components"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	service_constants "github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

type LayoutPanel struct {
	topology *components.DropdownSelector

	chkRoads               widget.Bool
	chkPortals             widget.Bool
	sldMaxPortals          widget.Float
	chkFootholds           widget.Bool
	sldRemoteFootholds     widget.Float
	chkPlayerIsolation     widget.Bool
	chkMatchPlayerFactions widget.Bool
	chkAbandonedOutposts   widget.Bool
	sldAbandonedOutposts   widget.Float
	sldMinNeutralBetween   widget.Float

	chkAdvancedZones       widget.Bool
	sldNeutralLowNoCastle  widget.Float
	sldNeutralLowCastle    widget.Float
	sldNeutralMedNoCastle  widget.Float
	sldNeutralMedCastle    widget.Float
	sldNeutralHighNoCastle widget.Float
	sldNeutralHighCastle   widget.Float

	sldNeutralLowCastlesPerZone  widget.Float
	sldNeutralMedCastlesPerZone  widget.Float
	sldNeutralHighCastlesPerZone widget.Float

	sldNeutralCount       widget.Float
	sldPlayerOwnedCastles widget.Float
	sldPlayerCastles      widget.Float
	sldNeutralCastles     widget.Float
	sldHubSize            widget.Float
	sldHubCastles         widget.Float
	sldPlayerZoneSize     widget.Float
	sldNeutralZoneSize    widget.Float
	sldGuardRandom        widget.Float
	sldResourceDensity    widget.Float
	sldStructureDensity   widget.Float
	sldNeutralStack       widget.Float
	sldBorderGuard        widget.Float

	editConnectionsBtn widget.Clickable
	btnPlayerContent   widget.Clickable
	btnLowContent      widget.Clickable
	btnMedContent      widget.Clickable
	btnHighContent     widget.Clickable
	btnHubContent      widget.Clickable

	scroll widget.List

	state *drivers.State
}

func NewLayoutPanel(state *drivers.State) *LayoutPanel {
	panel := &LayoutPanel{
		topology: components.NewDropdownSelector(func() []string {
			labels := make([]string, 0)
			for topology := range service_constants.GetTopologyDescriptorSeq() {
				labels = append(labels, topology.Label)
			}
			return labels
		}()),
		state: state,
	}
	panel.scroll.Axis = layout.Vertical
	panel.LoadFromState()
	return panel
}

func (this *LayoutPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{
		widgets.NewHorizontallySplitWidget(theme,
			func(theme *material.Theme) layout.Widget {
				return func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(this.getTopologySectionWidget(theme)),
						layout.Rigid(this.getConnectivityWidget(theme)),
						layout.Rigid(this.getZoneSizesWidget(theme)),
						layout.Rigid(this.getDifficultyAndDensityWidget(theme)),
					)
				}
			},
			func(theme *material.Theme) layout.Widget {
				return func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
						layout.Rigid(this.getManualZoneEditWidget(theme)),
						layout.Rigid(this.getZonesWidget(theme)),
					)
				}
			}),
	}

	return func(gtx layout.Context) layout.Dimensions {
		this.handleConnectionEditorClick(gtx)
		this.handleZoneContentDialogClicks(gtx)
		return material.List(theme, &this.scroll).Layout(
			gtx, len(widgetsList),
			func(gtx layout.Context, index int) layout.Dimensions { return widgetsList[index](gtx) })
	}
}

func (this *LayoutPanel) LoadFromState() {
	settings := this.state.GetStateData()

	this.topology.SelectByName(service_constants.GetTopologyDescriptorFromType(settings.Topology).Label)

	this.chkRoads.Value = settings.GenerateRoads
	this.chkPortals.Value = settings.RandomPortals
	this.sldMaxPortals.Value = utils.Normalize(float32(settings.MaxPortalConnections), 1, 32)
	this.chkFootholds.Value = settings.SpawnRemoteFootholds
	this.sldRemoteFootholds.Value = utils.Normalize(float32(settings.RemoteFootholdCount), 0, 4)
	this.chkPlayerIsolation.Value = settings.NoDirectPlayerConn
	this.chkMatchPlayerFactions.Value = settings.MatchPlayerCastleFactions
	this.chkAbandonedOutposts.Value = settings.SpawnAbandonedOutposts
	this.sldAbandonedOutposts.Value = utils.Normalize(float32(settings.AbandonedOutpostCount), 0, 4)
	this.sldMinNeutralBetween.Value = utils.Normalize(float32(settings.MinNeutralZonesBetweenPlayers), 0, 8)

	this.chkAdvancedZones.Value = settings.AdvancedMode
	this.sldNeutralCount.Value = utils.Normalize(float32(settings.NeutralZoneCount), 0, 16)
	this.sldPlayerOwnedCastles.Value = utils.Normalize(float32(settings.PlayerOwnedCastles), 0, 4)
	this.sldPlayerCastles.Value = utils.Normalize(float32(settings.PlayerZoneCastles), 0, 4)
	this.sldNeutralCastles.Value = utils.Normalize(float32(settings.NeutralZoneCastles), 0, 4)
	this.sldNeutralLowNoCastle.Value = utils.Normalize(float32(settings.NeutralLowNoCastleCount), 0, 8)
	this.sldNeutralLowCastle.Value = utils.Normalize(float32(settings.NeutralLowCastleCount), 0, 8)
	this.sldNeutralMedNoCastle.Value = utils.Normalize(float32(settings.NeutralMediumNoCastleCount), 0, 8)
	this.sldNeutralMedCastle.Value = utils.Normalize(float32(settings.NeutralMediumCastleCount), 0, 8)
	this.sldNeutralHighNoCastle.Value = utils.Normalize(float32(settings.NeutralHighNoCastleCount), 0, 8)
	this.sldNeutralHighCastle.Value = utils.Normalize(float32(settings.NeutralHighCastleCount), 0, 8)
	this.sldNeutralLowCastlesPerZone.Value = utils.Normalize(float32(settings.NeutralLowCastlesPerZone), 0, 4)
	this.sldNeutralMedCastlesPerZone.Value = utils.Normalize(float32(settings.NeutralMediumCastlesPerZone), 0, 4)
	this.sldNeutralHighCastlesPerZone.Value = utils.Normalize(float32(settings.NeutralHighCastlesPerZone), 0, 4)
	this.sldHubSize.Value = float32((settings.HubZoneSize - 0.5) / 1.5)
	this.sldHubCastles.Value = utils.Normalize(float32(settings.HubZoneCastles), 0, 4)
	this.sldPlayerZoneSize.Value = float32((settings.PlayerZoneSize - 0.5) / 1.5)
	this.sldNeutralZoneSize.Value = float32((settings.NeutralZoneSize - 0.5) / 1.5)
	this.sldGuardRandom.Value = utils.Normalize(float32(settings.GuardRandomization), 0, 0.5)
	this.sldResourceDensity.Value = utils.Normalize(float32(settings.ResourceDensityPercent), 25, 200)
	this.sldStructureDensity.Value = utils.Normalize(float32(settings.StructureDensityPercent), 25, 200)
	this.sldNeutralStack.Value = utils.Normalize(float32(settings.NeutralStackStrengthPercent), 25, 200)
	this.sldBorderGuard.Value = utils.Normalize(float32(settings.BorderGuardStrengthPercent), 25, 200)
}

func (this *LayoutPanel) SaveToState() {
	// TODO: check `.Update(gtx)` and on true update the value
	this.state.UpdateState(func(settings *dtos.EditorStateDto) {
		settings.Topology = this.getCurrentTopology().Type

		settings.GenerateRoads = this.chkRoads.Value
		settings.RandomPortals = this.chkPortals.Value
		settings.MaxPortalConnections = utils.RoundedRange(this.sldMaxPortals.Value, 1, 32)
		settings.SpawnRemoteFootholds = this.chkFootholds.Value
		settings.RemoteFootholdCount = utils.RoundedRange(this.sldRemoteFootholds.Value, 0, 4)
		settings.NoDirectPlayerConn = this.chkPlayerIsolation.Value
		settings.MatchPlayerCastleFactions = this.chkMatchPlayerFactions.Value
		settings.SpawnAbandonedOutposts = this.chkAbandonedOutposts.Value
		settings.AbandonedOutpostCount = utils.RoundedRange(this.sldAbandonedOutposts.Value, 0, 4)
		settings.MinNeutralZonesBetweenPlayers = utils.RoundedRange(this.sldMinNeutralBetween.Value, 0, 8)

		settings.AdvancedMode = this.chkAdvancedZones.Value
		settings.NeutralZoneCount = utils.RoundedRange(this.sldNeutralCount.Value, 0, 16)
		settings.PlayerOwnedCastles = utils.RoundedRange(this.sldPlayerOwnedCastles.Value, 0, 4)
		settings.PlayerZoneCastles = utils.RoundedRange(this.sldPlayerCastles.Value, 0, 4)
		settings.NeutralZoneCastles = utils.RoundedRange(this.sldNeutralCastles.Value, 0, 4)
		settings.NeutralLowNoCastleCount = utils.RoundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)
		settings.NeutralLowCastleCount = utils.RoundedRange(this.sldNeutralLowCastle.Value, 0, 8)
		settings.NeutralMediumNoCastleCount = utils.RoundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)
		settings.NeutralMediumCastleCount = utils.RoundedRange(this.sldNeutralMedCastle.Value, 0, 8)
		settings.NeutralHighNoCastleCount = utils.RoundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)
		settings.NeutralHighCastleCount = utils.RoundedRange(this.sldNeutralHighCastle.Value, 0, 8)
		settings.NeutralLowCastlesPerZone = utils.RoundedRange(this.sldNeutralLowCastlesPerZone.Value, 0, 4)
		settings.NeutralMediumCastlesPerZone = utils.RoundedRange(this.sldNeutralMedCastlesPerZone.Value, 0, 4)
		settings.NeutralHighCastlesPerZone = utils.RoundedRange(this.sldNeutralHighCastlesPerZone.Value, 0, 4)
		settings.HubZoneSize = float64(0.5 + this.sldHubSize.Value*1.5)
		settings.HubZoneCastles = utils.RoundedRange(this.sldHubCastles.Value, 0, 4)
		settings.PlayerZoneSize = float64(0.5 + this.sldPlayerZoneSize.Value*1.5)
		settings.NeutralZoneSize = float64(0.5 + this.sldNeutralZoneSize.Value*1.5)
		settings.GuardRandomization = float64(utils.Denormalize(this.sldGuardRandom.Value, 0, 0.5))
		settings.ResourceDensityPercent = utils.RoundedRange(this.sldResourceDensity.Value, 25, 200)
		settings.StructureDensityPercent = utils.RoundedRange(this.sldStructureDensity.Value, 25, 200)
		settings.NeutralStackStrengthPercent = utils.RoundedRange(this.sldNeutralStack.Value, 25, 200)
		settings.BorderGuardStrengthPercent = utils.RoundedRange(this.sldBorderGuard.Value, 25, 200)
	})
}

func (this *LayoutPanel) getTopologySectionWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Topology", []layout.Widget{
		widgets.NewLabeledRowWidget(
			theme, "Topology", constants.DefaultLabelWidth,
			func(gtx layout.Context) layout.Dimensions { return this.topology.Layout(gtx, theme) }),
		func(gtx layout.Context) layout.Dimensions {
			label := material.Caption(theme, this.getCurrentTopology().Description)
			label.Color = themes.ColorTextDim
			return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(constants.DefaultLabelWidth + 8)}.Layout(gtx, label.Layout)
		},
	})
}

func (this *LayoutPanel) getConnectivityWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Connectivity", []layout.Widget{
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkRoads, "Generate roads between zones"),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPortals, "Random portals (instead of fixed connections)"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkPortals.Value {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(
				theme, "Max portal connections", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(
					theme, &this.sldMaxPortals,
					utils.RoundedRangeString(this.sldMaxPortals.Value, 1, 32)))(gtx)
		},
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPlayerIsolation, "Disallow direct player-to-player connections"),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkMatchPlayerFactions, "Match player castle factions"),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkFootholds, "Spawn remote footholds"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkFootholds.Value {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(
				theme, "Remote footholds", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(
					theme, &this.sldRemoteFootholds,
					utils.RoundedRangeString(this.sldRemoteFootholds.Value, 0, 4)))(gtx)
		},
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkAbandonedOutposts, "Spawn abandoned outposts alongside neutral castles"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkAbandonedOutposts.Value {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(
				theme, "Abandoned outposts", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(
					theme, &this.sldAbandonedOutposts,
					utils.RoundedRangeString(this.sldAbandonedOutposts.Value, 0, 4)))(gtx)
		},
		// TODO: Investigate this. Is it used? How does it work? Seems like it does not do anything
		// or at minimum does not work as expected. Also if it works, range needs to be dynamic
		// based on current neutral and player zone counts
		// widgets.NewLabeledRowWidget(theme, "Min neutrals between players", 200, widgets.NewLabeledSliderWidget(theme, &this.sldMinNeutralBetween, fmt.Sprintf("%d", utils.RoundedRange(this.sldMinNeutralBetween.Value, 0, 8)))),
	})
}

func (this *LayoutPanel) getZoneSizesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Zone sizes", []layout.Widget{
		widgets.NewLabeledRowWidget(
			theme, "Player zone size", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldPlayerZoneSize,
				utils.MultiplierString(this.sldPlayerZoneSize.Value, 0.5, 1.5))),
		widgets.NewLabeledRowWidget(
			theme, "Neutral zone size", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldNeutralZoneSize,
				utils.MultiplierString(this.sldNeutralZoneSize.Value, 0.5, 1.5))),
		func(gtx layout.Context) layout.Dimensions {
			if this.state.GetStateData().Topology != config.TopologyHubAndSpoke {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(
				theme, "Hub zone size", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(
					theme, &this.sldHubSize,
					utils.MultiplierString(this.sldHubSize.Value, 0.5, 1.5)))(gtx)
		},
	})
}

func (this *LayoutPanel) getDifficultyAndDensityWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Difficulty & Density", []layout.Widget{
		widgets.NewLabeledRowWidget(
			theme, "Resource density", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldResourceDensity,
				utils.RoundedRangePercentString(this.sldResourceDensity.Value, 25, 200))),
		widgets.NewLabeledRowWidget(
			theme, "Structure density", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldStructureDensity,
				utils.RoundedRangePercentString(this.sldStructureDensity.Value, 25, 200))),
		widgets.NewLabeledRowWidget(
			theme, "Neutral stack strength", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldNeutralStack,
				utils.RoundedRangePercentString(this.sldNeutralStack.Value, 25, 200))),
		widgets.NewLabeledRowWidget(
			theme, "Border guard strength", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldBorderGuard,
				utils.RoundedRangePercentString(this.sldBorderGuard.Value, 25, 200))),
		widgets.NewLabeledRowWidget(
			theme, "Guard randomization", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldGuardRandom,
				utils.DenormalizeString(this.sldGuardRandom.Value, 0, 0.5))),
	})
}

func (this *LayoutPanel) getManualZoneEditWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Manual zone editing", []layout.Widget{
		widgets.NewBrightButtonWidget(theme, "Manual zone editor...", &this.editConnectionsBtn, false),
		widgets.NewDimmedLabelWidget(theme, "Visually add, move and edit zones and connections on the generated map."),
	})
}

func (this *LayoutPanel) getZonesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Zones", []layout.Widget{
		widgets.NewLabeledRowWidget(
			theme, "Player Owned castles per zone", constants.DefaultLabelWidthLong,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldPlayerOwnedCastles,
				utils.RoundedRangeString(this.sldPlayerOwnedCastles.Value, 0, 4))),
		widgets.NewLabeledRowWidget(
			theme, "Player Unclaimed castles per zone", constants.DefaultLabelWidthLong,
			widgets.NewLabeledSliderWidget(
				theme, &this.sldPlayerCastles,
				utils.RoundedRangeString(this.sldPlayerCastles.Value, 0, 4))),
		widgets.NewBrightButtonWidget(theme, "Edit player zone content...", &this.btnPlayerContent, false),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkAdvancedZones, "Advanced zone control (split low / medium / high tiers)"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkAdvancedZones.Value {
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(widgets.NewLabeledRowWidget(
						theme, "Total neutral zones", constants.DefaultLabelWidthLong,
						widgets.NewLabeledSliderWidget(
							theme, &this.sldNeutralCount,
							utils.RoundedRangeString(this.sldNeutralCount.Value, 0, 16)))),
					layout.Rigid(widgets.NewLabeledRowWidget(
						theme, "Neutral castles per zone", constants.DefaultLabelWidthLong,
						widgets.NewLabeledSliderWidget(
							theme, &this.sldNeutralCastles,
							utils.RoundedRangeString(this.sldNeutralCastles.Value, 0, 4)))),
				)
			}

			return this.getAdvancedZonesWidget(theme)(gtx)
		},
	})
}

func (this *LayoutPanel) getAdvancedZonesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Advanced options", []layout.Widget{
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
		widgets.NewLabeledRowWidget(
			theme, "No castle", constants.DefaultLabelWidthShort,
			widgets.NewLabeledSliderWidget(
				theme, noCastle,
				utils.RoundedRangeString(noCastle.Value, 0, 8))),
		widgets.NewLabeledRowWidget(
			theme, "With castle", constants.DefaultLabelWidthShort,
			widgets.NewLabeledSliderWidget(
				theme, withCastle,
				utils.RoundedRangeString(withCastle.Value, 0, 8))),
		widgets.NewLabeledRowWidget(
			theme, "Neutral castles per zone", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(
				theme, castlesPerZone,
				utils.RoundedRangeString(castlesPerZone.Value, 0, 4))),
		widgets.NewBrightButtonWidget(theme, "Edit zone content...", contentBtn, false),
	})
}

// getHubTierSectionWidget renders the advanced Hub sub-section. It only appears
// for the Hub & Spoke topology and (being nested inside the advanced options)
// only while advanced zone control is enabled.
func (this *LayoutPanel) getHubTierSectionWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		if this.state.GetStateData().Topology != config.TopologyHubAndSpoke {
			return layout.Dimensions{}
		}

		return widgets.NewSectionWidget(theme, "Hub", []layout.Widget{
			widgets.NewLabeledRowWidget(
				theme, "Hub zone castles", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(
					theme, &this.sldHubCastles,
					utils.RoundedRangeString(this.sldHubCastles.Value, 0, 4))),
			widgets.NewBrightButtonWidget(theme, "Edit zone content...", &this.btnHubContent, false),
		})(gtx)
	}
}

func (this *LayoutPanel) getCurrentTopology() service_constants.TopologyDescriptor {
	return service_constants.GetTopologyDescriptorFromIndex(this.topology.GetSelectedIndex())
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
	settings := this.state.GetStateData()
	generatorConfig := this.state.GetGeneratorConfig()
	tuning := models.NewGenerationTuning(generatorConfig, len(activeVariant.Zones))
	this.state.Dialogs().Open(dialogs.NewZoneEditorDialog(
		activeVariant.Zones,
		activeVariant.Connections,
		settings.Topology,
		tuning,
		settings.GenerateRoads,
		func(zones []entities.Zone, conns []entities.Connection) { this.state.ApplyEditedZones(zones, conns) },
	))
}

// handleZoneContentDialogClicks opens the single-tier zone-content editor for
// whichever per-zone "Edit zone content..." button was clicked this frame.
func (this *LayoutPanel) handleZoneContentDialogClicks(gtx layout.Context) {
	settings := this.state.GetStateData()
	switch {
	case this.btnPlayerContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Player", true, settings.PlayerZoneContentRows,
			func(s *dtos.EditorStateDto, rows []models.ZoneContentRowSave) { s.PlayerZoneContentRows = rows })
	case this.btnLowContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Low Neutral", false, settings.LowNeutralContentRows,
			func(s *dtos.EditorStateDto, rows []models.ZoneContentRowSave) { s.LowNeutralContentRows = rows })
	case this.btnMedContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Medium Neutral", false, settings.MediumNeutralContentRows,
			func(s *dtos.EditorStateDto, rows []models.ZoneContentRowSave) { s.MediumNeutralContentRows = rows })
	case this.btnHighContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: High Neutral", false, settings.HighNeutralContentRows,
			func(s *dtos.EditorStateDto, rows []models.ZoneContentRowSave) { s.HighNeutralContentRows = rows })
	case this.btnHubContent.Clicked(gtx):
		this.openZoneContentDialog("Zone Content: Hub", false, settings.HubZoneContentRows,
			func(s *dtos.EditorStateDto, rows []models.ZoneContentRowSave) { s.HubZoneContentRows = rows })
	}
}

// openZoneContentDialog opens a ZoneContentDialog for a single tier and writes
// the edited rows back into the editor state through the given setter.
func (this *LayoutPanel) openZoneContentDialog(
	title string,
	isPlayerTier bool,
	rows []models.ZoneContentRowSave,
	set func(*dtos.EditorStateDto, []models.ZoneContentRowSave),
) {
	this.state.Dialogs().Open(dialogs.NewZoneContentDialog(
		title, isPlayerTier, rows, this.state.Dialogs().Open,
		func(updated []models.ZoneContentRowSave) {
			this.state.UpdateState(func(s *dtos.EditorStateDto) { set(s, updated) })
		}))
}
