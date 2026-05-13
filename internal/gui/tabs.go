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

func (s *State) tabMapSetup(th *material.Theme) []layout.Widget {
	return []layout.Widget{
		section(th, "Template", []layout.Widget{
			labeledRowW(th, "Template name", 160, func(gtx layout.Context) layout.Dimensions {
				return drawEditor(gtx, th, &s.templateName, "Enter template name")
			}),
			labeledRowW(th, "Game mode", 160, func(gtx layout.Context) layout.Dimensions {
				return s.gameMode.Layout(gtx, th)
			}),
		}),
		section(th, "Map", []layout.Widget{
			labeledRowW(th, "Players", 160, func(gtx layout.Context) layout.Dimensions {
				p := roundedRange(s.playerCnt.Value, 2, 8)
				return sliderLabeled(gtx, th, &s.playerCnt, fmt.Sprintf("%d", p))
			}),
			labeledRowW(th, "Map size", 160, func(gtx layout.Context) layout.Dimensions {
				size := sliderToMapSize(s.mapSizeSld.Value, s.chkExpSizes.Value)
				lbl := fmt.Sprintf("%d × %d  (%s)", size, size, mapSizeLabelInt(size))
				return sliderLabeled(gtx, th, &s.mapSizeSld, lbl)
			}),
			checkRow(th, &s.chkExpSizes, "Allow experimental large map sizes (>240)"),
		}),
		section(th, "Topology", []layout.Widget{
			labeledRowW(th, "Topology", 160, func(gtx layout.Context) layout.Dimensions {
				return s.topology.Layout(gtx, th)
			}),
			func(gtx layout.Context) layout.Dimensions {
				lbl := material.Body2(th, topologyDescription(topologyValues[s.topology.Selected]))
				lbl.Color = colTextDim
				lbl.TextSize = unit.Sp(12)
				return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(168)}.Layout(gtx, lbl.Layout)
			},
		}),
	}
}

// ——————————————————————————————————————————————
// Tab 2: Generation Options (advanced map gen)
// ——————————————————————————————————————————————

func (s *State) tabGenerationOptions(th *material.Theme) []layout.Widget {
	widgets := []layout.Widget{
		section(th, "Connectivity", []layout.Widget{
			checkRow(th, &s.chkRoads, "Generate roads between zones"),
			checkRow(th, &s.chkPortals, "Random portals (instead of fixed connections)"),
			func(gtx layout.Context) layout.Dimensions {
				if !s.chkPortals.Value {
					return layout.Dimensions{}
				}
				return labeledRowW(th, "Max portal connections", 200, func(gtx layout.Context) layout.Dimensions {
					n := roundedRange(s.sldMaxPortals.Value, 1, 32)
					return sliderLabeled(gtx, th, &s.sldMaxPortals, fmt.Sprintf("%d", n))
				})(gtx)
			},
			checkRow(th, &s.chkFootholds, "Spawn remote footholds"),
			checkRow(th, &s.chkBalancedZones, "Experimental balanced zone placement"),
			checkRow(th, &s.chkPlayerIsolation, "Disallow direct player-to-player connections"),
			checkRow(th, &s.chkMatchPlayerFactions, "Match player castle factions"),
			labeledRowW(th, "Min neutrals between players", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldMinNeutralBetween.Value, 0, 8)
				return sliderLabeled(gtx, th, &s.sldMinNeutralBetween, fmt.Sprintf("%d", n))
			}),
		}),
		section(th, "Zone sizes", []layout.Widget{
			labeledRowW(th, "Player zone size", 200, func(gtx layout.Context) layout.Dimensions {
				size := 0.5 + float64(s.sldPlayerZoneSize.Value)*1.5
				return sliderLabeled(gtx, th, &s.sldPlayerZoneSize, fmt.Sprintf("× %.2f", size))
			}),
			labeledRowW(th, "Neutral zone size", 200, func(gtx layout.Context) layout.Dimensions {
				size := 0.5 + float64(s.sldNeutralZoneSize.Value)*1.5
				return sliderLabeled(gtx, th, &s.sldNeutralZoneSize, fmt.Sprintf("× %.2f", size))
			}),
			func(gtx layout.Context) layout.Dimensions {
				if topologyValues[s.topology.Selected] != generator.TopologyHubAndSpoke {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(labeledRowW(th, "Hub zone size", 200, func(gtx layout.Context) layout.Dimensions {
						size := 0.5 + float64(s.sldHubSize.Value)*1.5
						return sliderLabeled(gtx, th, &s.sldHubSize, fmt.Sprintf("× %.2f", size))
					})),
					layout.Rigid(labeledRowW(th, "Hub zone castles", 200, func(gtx layout.Context) layout.Dimensions {
						n := roundedRange(s.sldHubCastles.Value, 0, 4)
						return sliderLabeled(gtx, th, &s.sldHubCastles, fmt.Sprintf("%d", n))
					})),
				)
			},
		}),
		section(th, "Difficulty & Density", []layout.Widget{
			labeledRowW(th, "Resource density %", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldResourceDensity.Value, 25, 200)
				return sliderLabeled(gtx, th, &s.sldResourceDensity, fmt.Sprintf("%d%%", n))
			}),
			labeledRowW(th, "Structure density %", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldStructureDensity.Value, 25, 200)
				return sliderLabeled(gtx, th, &s.sldStructureDensity, fmt.Sprintf("%d%%", n))
			}),
			labeledRowW(th, "Neutral stack strength %", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralStack.Value, 25, 200)
				return sliderLabeled(gtx, th, &s.sldNeutralStack, fmt.Sprintf("%d%%", n))
			}),
			labeledRowW(th, "Border guard strength %", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldBorderGuard.Value, 25, 200)
				return sliderLabeled(gtx, th, &s.sldBorderGuard, fmt.Sprintf("%d%%", n))
			}),
			labeledRowW(th, "Guard randomization", 200, func(gtx layout.Context) layout.Dimensions {
				v := mapRange(s.sldGuardRandom.Value, 0, 0.5)
				return sliderLabeled(gtx, th, &s.sldGuardRandom, fmt.Sprintf("± %.2f", v))
			}),
		}),
		checkRow(th, &s.chkAdvancedZones, "Advanced zone control (split low / medium / high tiers)"),
	}

	if s.chkAdvancedZones.Value {
		widgets = append(widgets, section(th, "Zones (advanced)", []layout.Widget{
			labeledRowW(th, "Total neutral zones", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralCount.Value, 0, 16)
				return sliderLabeled(gtx, th, &s.sldNeutralCount, fmt.Sprintf("%d", n))
			}),
			labeledRowW(th, "Player castles per zone", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldPlayerCastles.Value, 0, 4)
				return sliderLabeled(gtx, th, &s.sldPlayerCastles, fmt.Sprintf("%d", n))
			}),
			labeledRowW(th, "Neutral castles per zone", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralCastles.Value, 0, 4)
				return sliderLabeled(gtx, th, &s.sldNeutralCastles, fmt.Sprintf("%d", n))
			}),
			dimLabelW(th, "Low tier"),
			labeledRowW(th, "  no castle", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralLowNoCastle.Value, 0, 8)
				return sliderLabeled(gtx, th, &s.sldNeutralLowNoCastle, fmt.Sprintf("%d", n))
			}),
			labeledRowW(th, "  with castle", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralLowCastle.Value, 0, 8)
				return sliderLabeled(gtx, th, &s.sldNeutralLowCastle, fmt.Sprintf("%d", n))
			}),
			dimLabelW(th, "Medium tier"),
			labeledRowW(th, "  no castle", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralMedNoCastle.Value, 0, 8)
				return sliderLabeled(gtx, th, &s.sldNeutralMedNoCastle, fmt.Sprintf("%d", n))
			}),
			labeledRowW(th, "  with castle", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralMedCastle.Value, 0, 8)
				return sliderLabeled(gtx, th, &s.sldNeutralMedCastle, fmt.Sprintf("%d", n))
			}),
			dimLabelW(th, "High tier"),
			labeledRowW(th, "  no castle", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralHighNoCastle.Value, 0, 8)
				return sliderLabeled(gtx, th, &s.sldNeutralHighNoCastle, fmt.Sprintf("%d", n))
			}),
			labeledRowW(th, "  with castle", 220, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldNeutralHighCastle.Value, 0, 8)
				return sliderLabeled(gtx, th, &s.sldNeutralHighCastle, fmt.Sprintf("%d", n))
			}),
		}))
	}
	return widgets
}

// ——————————————————————————————————————————————
// Tab 3: Game Rules
// ——————————————————————————————————————————————

func (s *State) tabGameRules(th *material.Theme) []layout.Widget {
	widgets := []layout.Widget{
		section(th, "Victory Condition", []layout.Widget{
			labeledRowW(th, "Victory", 160, func(gtx layout.Context) layout.Dimensions {
				return s.victory.Layout(gtx, th)
			}),
		}),
		section(th, "Loss Conditions", []layout.Widget{
			checkRow(th, &s.chkLostStartCity, "Lose if start city is captured"),
			func(gtx layout.Context) layout.Dimensions {
				if !s.chkLostStartCity.Value {
					return layout.Dimensions{}
				}
				return labeledRowW(th, "Grace period (days)", 200, func(gtx layout.Context) layout.Dimensions {
					n := roundedRange(s.sldLostCityDay.Value, 1, 30)
					return sliderLabeled(gtx, th, &s.sldLostCityDay, fmt.Sprintf("%d", n))
				})(gtx)
			},
			checkRow(th, &s.chkLostStartHero, "Lose if start hero is killed"),
		}),
		section(th, "City Hold", []layout.Widget{
			checkRow(th, &s.chkCityHold, "Win if you hold a target city"),
			func(gtx layout.Context) layout.Dimensions {
				if !s.chkCityHold.Value && s.victory.Selected != 2 {
					return layout.Dimensions{}
				}
				return labeledRowW(th, "Days to hold", 200, func(gtx layout.Context) layout.Dimensions {
					n := roundedRange(s.sldCityHoldDays.Value, 1, 30)
					return sliderLabeled(gtx, th, &s.sldCityHoldDays, fmt.Sprintf("%d", n))
				})(gtx)
			},
		}),
		section(th, "Gladiator Arena", []layout.Widget{
			checkRow(th, &s.chkGladiatorArena, "Enable gladiator arena"),
			func(gtx layout.Context) layout.Dimensions {
				if !s.chkGladiatorArena.Value {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(labeledRowW(th, "Days delay start", 200, func(gtx layout.Context) layout.Dimensions {
						n := roundedRange(s.sldGladiatorDelay.Value, 1, 90)
						return sliderLabeled(gtx, th, &s.sldGladiatorDelay, fmt.Sprintf("%d", n))
					})),
					layout.Rigid(labeledRowW(th, "Count days", 200, func(gtx layout.Context) layout.Dimensions {
						n := roundedRange(s.sldGladiatorCountDay.Value, 1, 14)
						return sliderLabeled(gtx, th, &s.sldGladiatorCountDay, fmt.Sprintf("%d", n))
					})),
				)
			},
		}),
		section(th, "Tournament", []layout.Widget{
			checkRow(th, &s.chkTournament, "Enable tournament"),
			func(gtx layout.Context) layout.Dimensions {
				if !s.chkTournament.Value && s.victory.Selected != 3 {
					return layout.Dimensions{}
				}
				return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(labeledRowW(th, "First tournament day", 200, func(gtx layout.Context) layout.Dimensions {
						n := roundedRange(s.sldTournamentDay.Value, 1, 60)
						return sliderLabeled(gtx, th, &s.sldTournamentDay, fmt.Sprintf("%d", n))
					})),
					layout.Rigid(labeledRowW(th, "Interval (days)", 200, func(gtx layout.Context) layout.Dimensions {
						n := roundedRange(s.sldTournamentInterval.Value, 1, 30)
						return sliderLabeled(gtx, th, &s.sldTournamentInterval, fmt.Sprintf("%d", n))
					})),
					layout.Rigid(labeledRowW(th, "Points to win", 200, func(gtx layout.Context) layout.Dimensions {
						n := roundedRange(s.sldTournamentPoints.Value, 1, 10)
						return sliderLabeled(gtx, th, &s.sldTournamentPoints, fmt.Sprintf("%d", n))
					})),
					layout.Rigid(checkRow(th, &s.chkTournamentSaveArmy, "Save army between rounds")),
				)
			},
		}),
		section(th, "Heroes", []layout.Widget{
			labeledRowW(th, "Hero count min", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldHeroMin.Value, 1, 16)
				return sliderLabeled(gtx, th, &s.sldHeroMin, fmt.Sprintf("%d", n))
			}),
			labeledRowW(th, "Hero count max", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldHeroMax.Value, 1, 16)
				return sliderLabeled(gtx, th, &s.sldHeroMax, fmt.Sprintf("%d", n))
			}),
			labeledRowW(th, "Increment", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldHeroIncr.Value, 1, 5)
				return sliderLabeled(gtx, th, &s.sldHeroIncr, fmt.Sprintf("%d", n))
			}),
		}),
		section(th, "Experience modifiers", []layout.Widget{
			labeledRowW(th, "Faction laws exp %", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldFactionLawsExp.Value, 25, 200)
				return sliderLabeled(gtx, th, &s.sldFactionLawsExp, fmt.Sprintf("%d%%", n))
			}),
			labeledRowW(th, "Astrology exp %", 200, func(gtx layout.Context) layout.Dimensions {
				n := roundedRange(s.sldAstrologyExp.Value, 25, 200)
				return sliderLabeled(gtx, th, &s.sldAstrologyExp, fmt.Sprintf("%d%%", n))
			}),
		}),
	}
	return widgets
}

// ——————————————————————————————————————————————
// Tab 4: Zone Content
// ——————————————————————————————————————————————

func (s *State) tabZoneContent(th *material.Theme) []layout.Widget {
	return []layout.Widget{
		func(gtx layout.Context) layout.Dimensions {
			return warnBannerW(th, "EXPERIMENTAL — Player zone mandatory content. Effects only apply on generation.")(gtx)
		},
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return toolbarButton{Text: "↺  Reset to defaults", Click: &s.btnZoneReset}.Layout(gtx, th)
				}),
			)
		},
		s.zcMines.Layout(th),
		s.zcTreasures.Layout(th),
		s.zcHires.Layout(th),
		s.zcBanks.Layout(th),
	}
}

// section wraps a group of widgets in a bordered panel under a header.
func section(th *material.Theme, title string, rows []layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Bottom: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(sectionHeaderW(th, title)),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return borderedPanel(gtx, unit.Dp(8), func(gtx layout.Context) layout.Dimensions {
						children := make([]layout.FlexChild, 0, len(rows)*2)
						for i, w := range rows {
							w := w
							if i > 0 {
								children = append(children, layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout))
							}
							children = append(children, layout.Rigid(w))
						}
						return layout.Flex{Axis: layout.Vertical}.Layout(gtx, children...)
					})
				}),
			)
		})
	}
}

func topologyDescription(t models.MapTopology) string {
	switch t {
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
