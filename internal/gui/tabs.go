package gui

import (
	"fmt"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

// ——————————————————————————————————————————————
// Tab 1: Map Setup
// ——————————————————————————————————————————————

func (this *State) tabMapSetup(theme *material.Theme) []layout.Widget {
	return []layout.Widget{
		section(theme, "Template", []layout.Widget{
			labeledRowW(theme, "Template name", 160, func(gtx layout.Context) layout.Dimensions {
				return drawEditor(gtx, theme, &this.templateName, "Enter template name")
			}),
			labeledRowW(theme, "Game mode", 160, func(gtx layout.Context) layout.Dimensions {
				return this.gameMode.Layout(gtx, theme)
			}),
		}),
		section(theme, "Map", []layout.Widget{
			labeledRowW(theme, "Players", 160, func(gtx layout.Context) layout.Dimensions {
				players := roundedRange(this.playerCnt.Value, 2, 8)
				return sliderLabeled(gtx, theme, &this.playerCnt, fmt.Sprintf("%d", players))
			}),
			labeledRowW(theme, "Map size", 160, func(gtx layout.Context) layout.Dimensions {
				size := sliderToMapSize(this.mapSizeSld.Value, this.chkExpSizes.Value)
				label := fmt.Sprintf("%d × %d  (%s)", size, size, mapSizeLabelInt(size))
				return sliderLabeled(gtx, theme, &this.mapSizeSld, label)
			}),
			checkRow(theme, &this.chkExpSizes, "Allow experimental large map sizes (>240)"),
		}),
		section(theme, "Topology", []layout.Widget{
			labeledRowW(theme, "Topology", 160, func(gtx layout.Context) layout.Dimensions {
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

func (this *State) tabGenerationOptions(theme *material.Theme) []layout.Widget {
	widgets := []layout.Widget{
		section(theme, "Connectivity", []layout.Widget{
			checkRow(theme, &this.chkRoads, "Generate roads between zones"),
			checkRow(theme, &this.chkPortals, "Random portals (instead of fixed connections)"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkPortals.Value {
					return layout.Dimensions{}
				}
				return labeledRowW(theme, "Max portal connections", 200, func(gtx layout.Context) layout.Dimensions {
					number := roundedRange(this.sldMaxPortals.Value, 1, 32)
					return sliderLabeled(gtx, theme, &this.sldMaxPortals, fmt.Sprintf("%d", number))
				})(gtx)
			},
			checkRow(theme, &this.chkFootholds, "Spawn remote footholds"),
			checkRow(theme, &this.chkBalancedZones, "Experimental balanced zone placement"),
			checkRow(theme, &this.chkPlayerIsolation, "Disallow direct player-to-player connections"),
			checkRow(theme, &this.chkMatchPlayerFactions, "Match player castle factions"),
			labeledRowW(theme, "Min neutrals between players", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldMinNeutralBetween.Value, 0, 8)
				return sliderLabeled(gtx, theme, &this.sldMinNeutralBetween, fmt.Sprintf("%d", number))
			}),
		}),
		section(theme, "Zone sizes", []layout.Widget{
			labeledRowW(theme, "Player zone size", 200, func(gtx layout.Context) layout.Dimensions {
				size := 0.5 + float64(this.sldPlayerZoneSize.Value)*1.5
				return sliderLabeled(gtx, theme, &this.sldPlayerZoneSize, fmt.Sprintf("× %.2f", size))
			}),
			labeledRowW(theme, "Neutral zone size", 200, func(gtx layout.Context) layout.Dimensions {
				size := 0.5 + float64(this.sldNeutralZoneSize.Value)*1.5
				return sliderLabeled(gtx, theme, &this.sldNeutralZoneSize, fmt.Sprintf("× %.2f", size))
			}),
			func(gtx layout.Context) layout.Dimensions {
				if topologyValues[this.topology.Selected] != generator.TopologyHubAndSpoke {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(labeledRowW(theme, "Hub zone size", 200, func(gtx layout.Context) layout.Dimensions {
						size := 0.5 + float64(this.sldHubSize.Value)*1.5
						return sliderLabeled(gtx, theme, &this.sldHubSize, fmt.Sprintf("× %.2f", size))
					})),
					layout.Rigid(labeledRowW(theme, "Hub zone castles", 200, func(gtx layout.Context) layout.Dimensions {
						number := roundedRange(this.sldHubCastles.Value, 0, 4)
						return sliderLabeled(gtx, theme, &this.sldHubCastles, fmt.Sprintf("%d", number))
					})),
				)
			},
		}),
		section(theme, "Difficulty & Density", []layout.Widget{
			labeledRowW(theme, "Resource density %", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldResourceDensity.Value, 25, 200)
				return sliderLabeled(gtx, theme, &this.sldResourceDensity, fmt.Sprintf("%d%%", number))
			}),
			labeledRowW(theme, "Structure density %", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldStructureDensity.Value, 25, 200)
				return sliderLabeled(gtx, theme, &this.sldStructureDensity, fmt.Sprintf("%d%%", number))
			}),
			labeledRowW(theme, "Neutral stack strength %", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralStack.Value, 25, 200)
				return sliderLabeled(gtx, theme, &this.sldNeutralStack, fmt.Sprintf("%d%%", number))
			}),
			labeledRowW(theme, "Border guard strength %", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldBorderGuard.Value, 25, 200)
				return sliderLabeled(gtx, theme, &this.sldBorderGuard, fmt.Sprintf("%d%%", number))
			}),
			labeledRowW(theme, "Guard randomization", 200, func(gtx layout.Context) layout.Dimensions {
				value := mapRange(this.sldGuardRandom.Value, 0, 0.5)
				return sliderLabeled(gtx, theme, &this.sldGuardRandom, fmt.Sprintf("± %.2f", value))
			}),
		}),
		checkRow(theme, &this.chkAdvancedZones, "Advanced zone control (split low / medium / high tiers)"),
	}

	if this.chkAdvancedZones.Value {
		widgets = append(widgets, section(theme, "Zones (advanced)", []layout.Widget{
			labeledRowW(theme, "Total neutral zones", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralCount.Value, 0, 16)
				return sliderLabeled(gtx, theme, &this.sldNeutralCount, fmt.Sprintf("%d", number))
			}),
			labeledRowW(theme, "Player castles per zone", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldPlayerCastles.Value, 0, 4)
				return sliderLabeled(gtx, theme, &this.sldPlayerCastles, fmt.Sprintf("%d", number))
			}),
			labeledRowW(theme, "Neutral castles per zone", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralCastles.Value, 0, 4)
				return sliderLabeled(gtx, theme, &this.sldNeutralCastles, fmt.Sprintf("%d", number))
			}),
			dimLabelW(theme, "Low tier"),
			labeledRowW(theme, "  no castle", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralLowNoCastle.Value, 0, 8)
				return sliderLabeled(gtx, theme, &this.sldNeutralLowNoCastle, fmt.Sprintf("%d", number))
			}),
			labeledRowW(theme, "  with castle", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralLowCastle.Value, 0, 8)
				return sliderLabeled(gtx, theme, &this.sldNeutralLowCastle, fmt.Sprintf("%d", number))
			}),
			dimLabelW(theme, "Medium tier"),
			labeledRowW(theme, "  no castle", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralMedNoCastle.Value, 0, 8)
				return sliderLabeled(gtx, theme, &this.sldNeutralMedNoCastle, fmt.Sprintf("%d", number))
			}),
			labeledRowW(theme, "  with castle", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralMedCastle.Value, 0, 8)
				return sliderLabeled(gtx, theme, &this.sldNeutralMedCastle, fmt.Sprintf("%d", number))
			}),
			dimLabelW(theme, "High tier"),
			labeledRowW(theme, "  no castle", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralHighNoCastle.Value, 0, 8)
				return sliderLabeled(gtx, theme, &this.sldNeutralHighNoCastle, fmt.Sprintf("%d", number))
			}),
			labeledRowW(theme, "  with castle", 220, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldNeutralHighCastle.Value, 0, 8)
				return sliderLabeled(gtx, theme, &this.sldNeutralHighCastle, fmt.Sprintf("%d", number))
			}),
		}))
	}
	return widgets
}

// ——————————————————————————————————————————————
// Tab 3: Game Rules
// ——————————————————————————————————————————————

func (this *State) tabGameRules(theme *material.Theme) []layout.Widget {
	widgets := []layout.Widget{
		section(theme, "Victory Condition", []layout.Widget{
			labeledRowW(theme, "Victory", 160, func(gtx layout.Context) layout.Dimensions {
				return this.victory.Layout(gtx, theme)
			}),
		}),
		section(theme, "Loss Conditions", []layout.Widget{
			checkRow(theme, &this.chkLostStartCity, "Lose if start city is captured"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkLostStartCity.Value {
					return layout.Dimensions{}
				}
				return labeledRowW(theme, "Grace period (days)", 200, func(gtx layout.Context) layout.Dimensions {
					number := roundedRange(this.sldLostCityDay.Value, 1, 30)
					return sliderLabeled(gtx, theme, &this.sldLostCityDay, fmt.Sprintf("%d", number))
				})(gtx)
			},
			checkRow(theme, &this.chkLostStartHero, "Lose if start hero is killed"),
		}),
		section(theme, "City Hold", []layout.Widget{
			checkRow(theme, &this.chkCityHold, "Win if you hold a target city"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkCityHold.Value && this.victory.Selected != 2 {
					return layout.Dimensions{}
				}
				return labeledRowW(theme, "Days to hold", 200, func(gtx layout.Context) layout.Dimensions {
					number := roundedRange(this.sldCityHoldDays.Value, 1, 30)
					return sliderLabeled(gtx, theme, &this.sldCityHoldDays, fmt.Sprintf("%d", number))
				})(gtx)
			},
		}),
		section(theme, "Gladiator Arena", []layout.Widget{
			checkRow(theme, &this.chkGladiatorArena, "Enable gladiator arena"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkGladiatorArena.Value {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(labeledRowW(theme, "Days delay start", 200, func(gtx layout.Context) layout.Dimensions {
						number := roundedRange(this.sldGladiatorDelay.Value, 1, 90)
						return sliderLabeled(gtx, theme, &this.sldGladiatorDelay, fmt.Sprintf("%d", number))
					})),
					layout.Rigid(labeledRowW(theme, "Count days", 200, func(gtx layout.Context) layout.Dimensions {
						number := roundedRange(this.sldGladiatorCountDay.Value, 1, 14)
						return sliderLabeled(gtx, theme, &this.sldGladiatorCountDay, fmt.Sprintf("%d", number))
					})),
				)
			},
		}),
		section(theme, "Tournament", []layout.Widget{
			checkRow(theme, &this.chkTournament, "Enable tournament"),
			func(gtx layout.Context) layout.Dimensions {
				if !this.chkTournament.Value && this.victory.Selected != 3 {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(labeledRowW(theme, "First tournament day", 200, func(gtx layout.Context) layout.Dimensions {
						number := roundedRange(this.sldTournamentDay.Value, 1, 60)
						return sliderLabeled(gtx, theme, &this.sldTournamentDay, fmt.Sprintf("%d", number))
					})),
					layout.Rigid(labeledRowW(theme, "Interval (days)", 200, func(gtx layout.Context) layout.Dimensions {
						number := roundedRange(this.sldTournamentInterval.Value, 1, 30)
						return sliderLabeled(gtx, theme, &this.sldTournamentInterval, fmt.Sprintf("%d", number))
					})),
					layout.Rigid(labeledRowW(theme, "Points to win", 200, func(gtx layout.Context) layout.Dimensions {
						number := roundedRange(this.sldTournamentPoints.Value, 1, 10)
						return sliderLabeled(gtx, theme, &this.sldTournamentPoints, fmt.Sprintf("%d", number))
					})),
					layout.Rigid(checkRow(theme, &this.chkTournamentSaveArmy, "Save army between rounds")),
				)
			},
		}),
		section(theme, "Heroes", []layout.Widget{
			labeledRowW(theme, "Hero count min", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldHeroMin.Value, 1, 16)
				return sliderLabeled(gtx, theme, &this.sldHeroMin, fmt.Sprintf("%d", number))
			}),
			labeledRowW(theme, "Hero count max", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldHeroMax.Value, 1, 16)
				return sliderLabeled(gtx, theme, &this.sldHeroMax, fmt.Sprintf("%d", number))
			}),
			labeledRowW(theme, "Increment", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldHeroIncr.Value, 1, 5)
				return sliderLabeled(gtx, theme, &this.sldHeroIncr, fmt.Sprintf("%d", number))
			}),
		}),
		section(theme, "Experience modifiers", []layout.Widget{
			labeledRowW(theme, "Faction laws exp %", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldFactionLawsExp.Value, 25, 200)
				return sliderLabeled(gtx, theme, &this.sldFactionLawsExp, fmt.Sprintf("%d%%", number))
			}),
			labeledRowW(theme, "Astrology exp %", 200, func(gtx layout.Context) layout.Dimensions {
				number := roundedRange(this.sldAstrologyExp.Value, 25, 200)
				return sliderLabeled(gtx, theme, &this.sldAstrologyExp, fmt.Sprintf("%d%%", number))
			}),
		}),
	}
	return widgets
}

// ——————————————————————————————————————————————
// Tab 4: Zone Content
// ——————————————————————————————————————————————

func (this *State) tabZoneContent(theme *material.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return warnBannerW(theme, "EXPERIMENTAL — Player zone mandatory content. Effects only apply on generation.")(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return toolbarButton{Text: "↺  Reset to defaults", Click: &this.btnZoneReset}.Layout(gtx, theme)
				}),
			)
		},
		this.zcMines.Layout(theme),
		this.zcTreasures.Layout(theme),
		this.zcHires.Layout(theme),
		this.zcBanks.Layout(theme),
	}
}

// section wraps a group of widgets in a bordered panel under a header.
func section(theme *material.Theme, title string, rows []layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(sectionHeaderW(theme, title)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return borderedPanel(gtx, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
						children := make([]layout.FlexChild, 0, len(rows)*2)
						for i, rowWidget := range rows {
							rowWidget := rowWidget
							if i > 0 {
								children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
							}
							children = append(children, layout.Rigid(rowWidget))
						}
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
					})
				}),
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

// (font import kept for future use)
var _ = font.Font{}
var _ = text.Start
