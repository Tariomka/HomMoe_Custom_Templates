package components

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

type GenerationPanel struct {
	chkRoads               widget.Bool
	chkPortals             widget.Bool
	sldMaxPortals          widget.Float
	chkFootholds           widget.Bool
	chkBalancedZones       widget.Bool
	chkPlayerIsolation     widget.Bool
	chkMatchPlayerFactions widget.Bool
	sldMinNeutralBetween   widget.Float

	chkAdvancedZones       widget.Bool
	sldNeutralLowNoCastle  widget.Float
	sldNeutralLowCastle    widget.Float
	sldNeutralMedNoCastle  widget.Float
	sldNeutralMedCastle    widget.Float
	sldNeutralHighNoCastle widget.Float
	sldNeutralHighCastle   widget.Float
	sldNeutralCount        widget.Float
	sldPlayerCastles       widget.Float
	sldNeutralCastles      widget.Float
	sldHubSize             widget.Float
	sldHubCastles          widget.Float
	sldPlayerZoneSize      widget.Float
	sldNeutralZoneSize     widget.Float
	sldGuardRandom         widget.Float
	sldResourceDensity     widget.Float
	sldStructureDensity    widget.Float
	sldNeutralStack        widget.Float
	sldBorderGuard         widget.Float

	scroll widget.List

	state *State
}

func NewGenerationPanel(state *State) *GenerationPanel {
	panel := &GenerationPanel{state: state}
	panel.scroll.Axis = layout.Vertical
	panel.LoadFromState()
	return panel
}

func (this *GenerationPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{
		widgets.NewSectionWidget(theme, "Connectivity", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkRoads, "Generate roads between zones"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPortals, "Random portals (instead of fixed connections)"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkPortals.Value {
					return layout.Dimensions{}
				}
				return widgets.NewLabeledRowWidget(theme, "Max portal connections", 200, widgets.NewLabeledSlider(theme, &this.sldMaxPortals, fmt.Sprintf("%d", utils.RoundedRange(this.sldMaxPortals.Value, 1, 32))))(gtx)
			},
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkFootholds, "Spawn remote footholds"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkBalancedZones, "Experimental balanced zone placement"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPlayerIsolation, "Disallow direct player-to-player connections"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkMatchPlayerFactions, "Match player castle factions"),
			widgets.NewLabeledRowWidget(theme, "Min neutrals between players", 200, widgets.NewLabeledSlider(theme, &this.sldMinNeutralBetween, fmt.Sprintf("%d", utils.RoundedRange(this.sldMinNeutralBetween.Value, 0, 8)))),
		}),
		widgets.NewSectionWidget(theme, "Zone sizes", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Player zone size", 200, widgets.NewLabeledSlider(theme, &this.sldPlayerZoneSize, fmt.Sprintf("× %.2f", 0.5+float64(this.sldPlayerZoneSize.Value)*1.5))),
			widgets.NewLabeledRowWidget(theme, "Neutral zone size", 200, widgets.NewLabeledSlider(theme, &this.sldNeutralZoneSize, fmt.Sprintf("× %.2f", 0.5+float64(this.sldNeutralZoneSize.Value)*1.5))),
			func(gtx layout.Context) layout.Dimensions {
				if this.state.GetSettingsFile().Topology != generator.TopologyHubAndSpoke {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Hub zone size", 200, widgets.NewLabeledSlider(theme, &this.sldHubSize, fmt.Sprintf("× %.2f", 0.5+float64(this.sldHubSize.Value)*1.5)))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Hub zone castles", 200, widgets.NewLabeledSlider(theme, &this.sldHubCastles, fmt.Sprintf("%d", utils.RoundedRange(this.sldHubCastles.Value, 0, 4))))),
				)
			},
		}),
		widgets.NewSectionWidget(theme, "Difficulty & Density", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Resource density %", 200, widgets.NewLabeledSlider(theme, &this.sldResourceDensity, fmt.Sprintf("%d%%", utils.RoundedRange(this.sldResourceDensity.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Structure density %", 200, widgets.NewLabeledSlider(theme, &this.sldStructureDensity, fmt.Sprintf("%d%%", utils.RoundedRange(this.sldStructureDensity.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Neutral stack strength %", 200, widgets.NewLabeledSlider(theme, &this.sldNeutralStack, fmt.Sprintf("%d%%", utils.RoundedRange(this.sldNeutralStack.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Border guard strength %", 200, widgets.NewLabeledSlider(theme, &this.sldBorderGuard, fmt.Sprintf("%d%%", utils.RoundedRange(this.sldBorderGuard.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Guard randomization", 200, widgets.NewLabeledSlider(theme, &this.sldGuardRandom, fmt.Sprintf("± %.2f", utils.Denormalize(this.sldGuardRandom.Value, 0, 0.5)))),
		}),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkAdvancedZones, "Advanced zone control (split low / medium / high tiers)"),
	}
	if this.chkAdvancedZones.Value {
		widgetsList = append(widgetsList, widgets.NewSectionWidget(theme, "Zones (advanced)", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Total neutral zones", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralCount, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralCount.Value, 0, 16)))),
			widgets.NewLabeledRowWidget(theme, "Player castles per zone", 220, widgets.NewLabeledSlider(theme, &this.sldPlayerCastles, fmt.Sprintf("%d", utils.RoundedRange(this.sldPlayerCastles.Value, 0, 4)))),
			widgets.NewLabeledRowWidget(theme, "Neutral castles per zone", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralCastles, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralCastles.Value, 0, 4)))),
			widgets.NewDimmedLabelWidget(theme, "Low tier"),
			widgets.NewLabeledRowWidget(theme, "  no castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralLowNoCastle, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)))),
			widgets.NewLabeledRowWidget(theme, "  with castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralLowCastle, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralLowCastle.Value, 0, 8)))),
			widgets.NewDimmedLabelWidget(theme, "Medium tier"),
			widgets.NewLabeledRowWidget(theme, "  no castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralMedNoCastle, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)))),
			widgets.NewLabeledRowWidget(theme, "  with castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralMedCastle, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralMedCastle.Value, 0, 8)))),
			widgets.NewDimmedLabelWidget(theme, "High tier"),
			widgets.NewLabeledRowWidget(theme, "  no castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralHighNoCastle, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)))),
			widgets.NewLabeledRowWidget(theme, "  with castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralHighCastle, fmt.Sprintf("%d", utils.RoundedRange(this.sldNeutralHighCastle.Value, 0, 8)))),
		}))
	}
	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

func (this *GenerationPanel) LoadFromState() {
	settings := this.state.GetSettingsFile()
	this.chkRoads.Value = settings.GenerateRoads
	this.chkPortals.Value = settings.RandomPortals
	this.sldMaxPortals.Value = utils.Normalize(float32(settings.MaxPortalConnections), 1, 32)
	this.chkFootholds.Value = settings.SpawnRemoteFootholds
	this.chkBalancedZones.Value = settings.ExperimentalBalancedZonePlacement
	this.chkPlayerIsolation.Value = settings.NoDirectPlayerConn
	this.chkMatchPlayerFactions.Value = settings.MatchPlayerCastleFactions
	this.sldMinNeutralBetween.Value = utils.Normalize(float32(settings.MinNeutralZonesBetweenPlayers), 0, 8)

	this.chkAdvancedZones.Value = settings.AdvancedMode
	this.sldNeutralCount.Value = utils.Normalize(float32(settings.NeutralZoneCount), 0, 16)
	this.sldPlayerCastles.Value = utils.Normalize(float32(settings.PlayerZoneCastles), 0, 4)
	this.sldNeutralCastles.Value = utils.Normalize(float32(settings.NeutralZoneCastles), 0, 4)
	this.sldNeutralLowNoCastle.Value = utils.Normalize(float32(settings.NeutralLowNoCastleCount), 0, 8)
	this.sldNeutralLowCastle.Value = utils.Normalize(float32(settings.NeutralLowCastleCount), 0, 8)
	this.sldNeutralMedNoCastle.Value = utils.Normalize(float32(settings.NeutralMediumNoCastleCount), 0, 8)
	this.sldNeutralMedCastle.Value = utils.Normalize(float32(settings.NeutralMediumCastleCount), 0, 8)
	this.sldNeutralHighNoCastle.Value = utils.Normalize(float32(settings.NeutralHighNoCastleCount), 0, 8)
	this.sldNeutralHighCastle.Value = utils.Normalize(float32(settings.NeutralHighCastleCount), 0, 8)
	this.sldHubSize.Value = float32((settings.HubZoneSize - 0.5) / 1.5)
	this.sldHubCastles.Value = utils.Normalize(float32(settings.HubZoneCastles), 0, 4)
	this.sldPlayerZoneSize.Value = float32((settings.PlayerZoneSize - 0.5) / 1.5)
	this.sldNeutralZoneSize.Value = float32((settings.NeutralZoneSize - 0.5) / 1.5)
	this.sldGuardRandom.Value = utils.Normalize(float32(settings.GuardRandomization), 0, 0.5)
	this.sldResourceDensity.Value = utils.Normalize(float32(settings.EffectiveResourceDensity()), 25, 200)
	this.sldStructureDensity.Value = utils.Normalize(float32(settings.EffectiveStructureDensity()), 25, 200)
	this.sldNeutralStack.Value = utils.Normalize(float32(settings.NeutralStackStrengthPercent), 25, 200)
	this.sldBorderGuard.Value = utils.Normalize(float32(settings.BorderGuardStrengthPercent), 25, 200)
}

// TODO: check `.Update(gtx)` and on true update the value
func (this *GenerationPanel) SaveToState() {
	this.state.UpdateState(func(settings *models.SettingsFile) {
		settings.GenerateRoads = this.chkRoads.Value
		settings.RandomPortals = this.chkPortals.Value
		settings.MaxPortalConnections = utils.RoundedRange(this.sldMaxPortals.Value, 1, 32)
		settings.SpawnRemoteFootholds = this.chkFootholds.Value
		settings.ExperimentalBalancedZonePlacement = this.chkBalancedZones.Value
		settings.NoDirectPlayerConn = this.chkPlayerIsolation.Value
		settings.MatchPlayerCastleFactions = this.chkMatchPlayerFactions.Value
		settings.MinNeutralZonesBetweenPlayers = utils.RoundedRange(this.sldMinNeutralBetween.Value, 0, 8)

		settings.AdvancedMode = this.chkAdvancedZones.Value
		settings.NeutralZoneCount = utils.RoundedRange(this.sldNeutralCount.Value, 0, 16)
		settings.PlayerZoneCastles = utils.RoundedRange(this.sldPlayerCastles.Value, 0, 4)
		settings.NeutralZoneCastles = utils.RoundedRange(this.sldNeutralCastles.Value, 0, 4)
		settings.NeutralLowNoCastleCount = utils.RoundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)
		settings.NeutralLowCastleCount = utils.RoundedRange(this.sldNeutralLowCastle.Value, 0, 8)
		settings.NeutralMediumNoCastleCount = utils.RoundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)
		settings.NeutralMediumCastleCount = utils.RoundedRange(this.sldNeutralMedCastle.Value, 0, 8)
		settings.NeutralHighNoCastleCount = utils.RoundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)
		settings.NeutralHighCastleCount = utils.RoundedRange(this.sldNeutralHighCastle.Value, 0, 8)
		settings.HubZoneSize = float64(0.5 + this.sldHubSize.Value*1.5)
		settings.HubZoneCastles = utils.RoundedRange(this.sldHubCastles.Value, 0, 4)
		settings.PlayerZoneSize = float64(0.5 + this.sldPlayerZoneSize.Value*1.5)
		settings.NeutralZoneSize = float64(0.5 + this.sldNeutralZoneSize.Value*1.5)
		settings.GuardRandomization = float64(utils.Denormalize(this.sldGuardRandom.Value, 0, 0.5))
		rd := utils.RoundedRange(this.sldResourceDensity.Value, 25, 200)
		sd := utils.RoundedRange(this.sldStructureDensity.Value, 25, 200)
		settings.ResourceDensityPercent = &rd
		settings.StructureDensityPercent = &sd
		settings.NeutralStackStrengthPercent = utils.RoundedRange(this.sldNeutralStack.Value, 25, 200)
		settings.BorderGuardStrengthPercent = utils.RoundedRange(this.sldBorderGuard.Value, 25, 200)
	})
}
