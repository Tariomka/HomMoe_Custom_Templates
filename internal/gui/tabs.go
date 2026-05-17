package gui

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

// ——————————————————————————————————————————————
// Tab 1: Map Setup
// ——————————————————————————————————————————————

func (this *WindowOld) tabMapSetup(theme *material.Theme) []layout.Widget {
	return []layout.Widget{
		NewSectionWidget(theme, "Template", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Template name", 160, widgets.NewTextboxWidget(theme, &this.templateName, "Enter template name")),
			widgets.NewLabeledRowWidget(theme, "Game mode", 160, func(gtx layout.Context) layout.Dimensions {
				return this.gameMode.Layout(gtx, theme)
			}),
		}),
		NewSectionWidget(theme, "Map", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Players", 160, widgets.NewLabeledSlider(theme, &this.playerCnt, fmt.Sprintf("%d", roundedRange(this.playerCnt.Value, 2, 8)))),
			widgets.NewLabeledRowWidget(theme, "Map size", 160, func(gtx layout.Context) layout.Dimensions {
				size := sliderToMapSize(this.mapSizeSld.Value, this.chkExpSizes.Value)
				label := fmt.Sprintf("%d × %d  (%s)", size, size, mapSizeLabelInt(size))
				return widgets.NewLabeledSlider(theme, &this.mapSizeSld, label)(gtx)
			}),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkExpSizes, "Allow experimental large map sizes (>240)"),
		}),
		NewSectionWidget(theme, "Topology", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Topology", 160, func(gtx layout.Context) layout.Dimensions {
				return this.topology.Layout(gtx, theme)
			}),
			func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, topologyDescription(topologyValues[this.topology.Selected]))
				label.Color = colTextDim
				label.TextSize = unit.Sp(12)
				return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(168)}.Layout(gtx, label.Layout)
			},
		}),
	}
}

// ——————————————————————————————————————————————
// Tab 2: Generation Options (advanced map gen)
// ——————————————————————————————————————————————

func (this *WindowOld) tabGenerationOptions(theme *material.Theme) []layout.Widget {
	widgetLayout := []layout.Widget{
		NewSectionWidget(theme, "Connectivity", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkRoads, "Generate roads between zones"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPortals, "Random portals (instead of fixed connections)"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkPortals.Value {
					return layout.Dimensions{}
				}
				return widgets.NewLabeledRowWidget(theme, "Max portal connections", 200, widgets.NewLabeledSlider(theme, &this.sldMaxPortals, fmt.Sprintf("%d", roundedRange(this.sldMaxPortals.Value, 1, 32))))(gtx)
			},
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkFootholds, "Spawn remote footholds"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkBalancedZones, "Experimental balanced zone placement"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPlayerIsolation, "Disallow direct player-to-player connections"),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkMatchPlayerFactions, "Match player castle factions"),
			widgets.NewLabeledRowWidget(theme, "Min neutrals between players", 200, widgets.NewLabeledSlider(theme, &this.sldMinNeutralBetween, fmt.Sprintf("%d", roundedRange(this.sldMinNeutralBetween.Value, 0, 8)))),
		}),
		NewSectionWidget(theme, "Zone sizes", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Player zone size", 200, widgets.NewLabeledSlider(theme, &this.sldPlayerZoneSize, fmt.Sprintf("× %.2f", 0.5+float64(this.sldPlayerZoneSize.Value)*1.5))),
			widgets.NewLabeledRowWidget(theme, "Neutral zone size", 200, widgets.NewLabeledSlider(theme, &this.sldNeutralZoneSize, fmt.Sprintf("× %.2f", 0.5+float64(this.sldNeutralZoneSize.Value)*1.5))),
			func(gtx layout.Context) layout.Dimensions {
				if topologyValues[this.topology.Selected] != generator.TopologyHubAndSpoke {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Hub zone size", 200, widgets.NewLabeledSlider(theme, &this.sldHubSize, fmt.Sprintf("× %.2f", 0.5+float64(this.sldHubSize.Value)*1.5)))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Hub zone castles", 200, widgets.NewLabeledSlider(theme, &this.sldHubCastles, fmt.Sprintf("%d", roundedRange(this.sldHubCastles.Value, 0, 4))))),
				)
			},
		}),
		NewSectionWidget(theme, "Difficulty & Density", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Resource density %", 200, widgets.NewLabeledSlider(theme, &this.sldResourceDensity, fmt.Sprintf("%d%%", roundedRange(this.sldResourceDensity.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Structure density %", 200, widgets.NewLabeledSlider(theme, &this.sldStructureDensity, fmt.Sprintf("%d%%", roundedRange(this.sldStructureDensity.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Neutral stack strength %", 200, widgets.NewLabeledSlider(theme, &this.sldNeutralStack, fmt.Sprintf("%d%%", roundedRange(this.sldNeutralStack.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Border guard strength %", 200, widgets.NewLabeledSlider(theme, &this.sldBorderGuard, fmt.Sprintf("%d%%", roundedRange(this.sldBorderGuard.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Guard randomization", 200, widgets.NewLabeledSlider(theme, &this.sldGuardRandom, fmt.Sprintf("± %.2f", utils.Denormalize(this.sldGuardRandom.Value, 0, 0.5)))),
		}),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkAdvancedZones, "Advanced zone control (split low / medium / high tiers)"),
	}

	if this.chkAdvancedZones.Value {
		widgetLayout = append(widgetLayout, NewSectionWidget(theme, "Zones (advanced)", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Total neutral zones", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralCount, fmt.Sprintf("%d", roundedRange(this.sldNeutralCount.Value, 0, 16)))),
			widgets.NewLabeledRowWidget(theme, "Player castles per zone", 220, widgets.NewLabeledSlider(theme, &this.sldPlayerCastles, fmt.Sprintf("%d", roundedRange(this.sldPlayerCastles.Value, 0, 4)))),
			widgets.NewLabeledRowWidget(theme, "Neutral castles per zone", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralCastles, fmt.Sprintf("%d", roundedRange(this.sldNeutralCastles.Value, 0, 4)))),
			widgets.NewDimmedLabelWidget(theme, "Low tier"),
			widgets.NewLabeledRowWidget(theme, "  no castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralLowNoCastle, fmt.Sprintf("%d", roundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)))),
			widgets.NewLabeledRowWidget(theme, "  with castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralLowCastle, fmt.Sprintf("%d", roundedRange(this.sldNeutralLowCastle.Value, 0, 8)))),
			widgets.NewDimmedLabelWidget(theme, "Medium tier"),
			widgets.NewLabeledRowWidget(theme, "  no castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralMedNoCastle, fmt.Sprintf("%d", roundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)))),
			widgets.NewLabeledRowWidget(theme, "  with castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralMedCastle, fmt.Sprintf("%d", roundedRange(this.sldNeutralMedCastle.Value, 0, 8)))),
			widgets.NewDimmedLabelWidget(theme, "High tier"),
			widgets.NewLabeledRowWidget(theme, "  no castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralHighNoCastle, fmt.Sprintf("%d", roundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)))),
			widgets.NewLabeledRowWidget(theme, "  with castle", 220, widgets.NewLabeledSlider(theme, &this.sldNeutralHighCastle, fmt.Sprintf("%d", roundedRange(this.sldNeutralHighCastle.Value, 0, 8)))),
		}))
	}
	return widgetLayout
}

// ——————————————————————————————————————————————
// Tab 3: Game Rules
// ——————————————————————————————————————————————

func (this *WindowOld) tabGameRules(theme *material.Theme) []layout.Widget {
	widgetLayout := []layout.Widget{
		NewSectionWidget(theme, "Victory Condition", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Victory", 160, func(gtx layout.Context) layout.Dimensions {
				return this.victory.Layout(gtx, theme)
			}),
		}),
		NewSectionWidget(theme, "Loss Conditions", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkLostStartCity, "Lose if start city is captured"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkLostStartCity.Value {
					return layout.Dimensions{}
				}
				return widgets.NewLabeledRowWidget(theme, "Grace period (days)", 200, widgets.NewLabeledSlider(theme, &this.sldLostCityDay, fmt.Sprintf("%d", roundedRange(this.sldLostCityDay.Value, 1, 30))))(gtx)
			},
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkLostStartHero, "Lose if start hero is killed"),
		}),
		NewSectionWidget(theme, "City Hold", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkCityHold, "Win if you hold a target city"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkCityHold.Value && this.victory.Selected != 2 {
					return layout.Dimensions{}
				}
				return widgets.NewLabeledRowWidget(theme, "Days to hold", 200, widgets.NewLabeledSlider(theme, &this.sldCityHoldDays, fmt.Sprintf("%d", roundedRange(this.sldCityHoldDays.Value, 1, 30))))(gtx)
			},
		}),
		NewSectionWidget(theme, "Gladiator Arena", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkGladiatorArena, "Enable gladiator arena"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkGladiatorArena.Value {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Days delay start", 200, widgets.NewLabeledSlider(theme, &this.sldGladiatorDelay, fmt.Sprintf("%d", roundedRange(this.sldGladiatorDelay.Value, 1, 90))))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Count days", 200, widgets.NewLabeledSlider(theme, &this.sldGladiatorCountDay, fmt.Sprintf("%d", roundedRange(this.sldGladiatorCountDay.Value, 1, 14))))),
				)
			},
		}),
		NewSectionWidget(theme, "Tournament", []layout.Widget{
			widgets.NewLabeledCheckboxRowWidget(theme, &this.chkTournament, "Enable tournament"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkTournament.Value && this.victory.Selected != 3 {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "First tournament day", 200, widgets.NewLabeledSlider(theme, &this.sldTournamentDay, fmt.Sprintf("%d", roundedRange(this.sldTournamentDay.Value, 1, 60))))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Interval (days)", 200, widgets.NewLabeledSlider(theme, &this.sldTournamentInterval, fmt.Sprintf("%d", roundedRange(this.sldTournamentInterval.Value, 1, 30))))),
					layout.Rigid(widgets.NewLabeledRowWidget(theme, "Points to win", 200, widgets.NewLabeledSlider(theme, &this.sldTournamentPoints, fmt.Sprintf("%d", roundedRange(this.sldTournamentPoints.Value, 1, 10))))),
					layout.Rigid(widgets.NewLabeledCheckboxRowWidget(theme, &this.chkTournamentSaveArmy, "Save army between rounds")),
				)
			},
		}),
		NewSectionWidget(theme, "Heroes", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Hero count min", 200, widgets.NewLabeledSlider(theme, &this.sldHeroMin, fmt.Sprintf("%d", roundedRange(this.sldHeroMin.Value, 1, 16)))),
			widgets.NewLabeledRowWidget(theme, "Hero count max", 200, widgets.NewLabeledSlider(theme, &this.sldHeroMax, fmt.Sprintf("%d", roundedRange(this.sldHeroMax.Value, 1, 16)))),
			widgets.NewLabeledRowWidget(theme, "Increment", 200, widgets.NewLabeledSlider(theme, &this.sldHeroIncr, fmt.Sprintf("%d", roundedRange(this.sldHeroIncr.Value, 1, 5)))),
		}),
		NewSectionWidget(theme, "Experience modifiers", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Faction laws exp %", 200, widgets.NewLabeledSlider(theme, &this.sldFactionLawsExp, fmt.Sprintf("%d%%", roundedRange(this.sldFactionLawsExp.Value, 25, 200)))),
			widgets.NewLabeledRowWidget(theme, "Astrology exp %", 200, widgets.NewLabeledSlider(theme, &this.sldAstrologyExp, fmt.Sprintf("%d%%", roundedRange(this.sldAstrologyExp.Value, 25, 200)))),
		}),
	}
	return widgetLayout
}

// ——————————————————————————————————————————————
// Tab 4: Zone Content
// ——————————————————————————————————————————————

func (this *WindowOld) tabZoneContent(theme *material.Theme) []layout.Widget {
	return []layout.Widget{
		widgets.NewWarningBannerWidget(theme, "EXPERIMENTAL — Player zone mandatory content. Effects only apply on generation."),
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, layout.Rigid(widgets.NewButtonWidget(theme, "↺  Reset to defaults", &this.btnZoneReset, false)))
		},
		this.zcMines.Layout(theme),
		this.zcTreasures.Layout(theme),
		this.zcHires.Layout(theme),
		this.zcBanks.Layout(theme),
	}
}

// NewSectionWidget wraps a group of widgets in a bordered panel under a header.
func NewSectionWidget(theme *material.Theme, title string, rows []layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(widgets.NewSectionHeaderWidget(theme, title)),
				layout.Rigid(widgets.NewPanelWidget(unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
					children := make([]layout.FlexChild, 0, len(rows)*2)
					for i, rowWidget := range rows {
						if i > 0 {
							children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
						}
						children = append(children, layout.Rigid(rowWidget))
					}
					return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
				})),
			)
		})
	}
}

func topologyDescription(topology models.MapTopology) string {
	switch topology {
	case generator.TopologyDefault:
		return "Ring: each player borders two neighbors in a closed loop."
	case generator.TopologyHubAndSpoke:
		return "Hub: central neutral hub connects all player zones."
	case generator.TopologyChain:
		return "Chain: linear series, harder for outer players to interact."
	case generator.TopologySharedWeb:
		return "Shared web: heavy interconnection through central neutral mesh."
	default:
		return "Random: layout decided by the generator."
	}
}
