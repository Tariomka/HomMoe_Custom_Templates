package gui

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/internal/generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// Domain knowledge ported from KnownValues.cs and TopologyOptions in MainWindow.xaml.cs.
var (
	gameModes      = []string{"Classic", "SingleHero"}
	mapSizes       = []int{64, 80, 96, 112, 128, 144, 160, 176, 192, 208, 240}
	topologyLabels = []string{"Random", "Ring", "Hub", "Chain"}
	topologyValues = []models.MapTopology{
		models.TopologyRandom,
		models.TopologyDefault,
		models.TopologyHubAndSpoke,
		models.TopologyChain,
	}
	victoryLabels = []string{"Standard", "Lost Starting City", "Hold City", "Tournament"}
)

func mapSizeLabel(size int) string {
	switch {
	case size == 64:
		return "S"
	case size == 80 || size == 96:
		return "M"
	case size == 112 || size == 128:
		return "L"
	case size == 144 || size == 160:
		return "XL"
	case size == 176 || size == 192:
		return "H"
	case size >= 208 && size <= 256:
		return "G"
	default:
		return "C"
	}
}

// State holds all interactive widget state and the underlying settings model.
type State struct {
	// Header / file controls.
	templateName    widget.Editor
	outputPath      widget.Editor
	btnGenerate     widget.Clickable
	btnSave         widget.Clickable
	btnBrowseOutput widget.Clickable
	btnRevealOutput widget.Clickable

	// Settings widgets.
	gameMode  *segmentGroup
	topology  *segmentGroup
	victory   *segmentGroup
	mapSize   widget.Float // 0..1, mapped to mapSizes index
	playerCnt widget.Float // 2..8

	// Toggles.
	chkRoads              widget.Bool
	chkPortals            widget.Bool
	chkFootholds          widget.Bool
	chkPlayerIsolation    widget.Bool
	chkBalancedZones      widget.Bool
	chkCityHold           widget.Bool
	chkShowDescription    widget.Bool
	chkIncludeOptionsDesc widget.Bool

	// Advanced sliders.
	sldGuardRandom widget.Float // 0..0.5
	sldConnPerZone widget.Float // 1..4
	sldNeutralLow  widget.Float // 0..8
	sldNeutralMed  widget.Float // 0..8
	sldNeutralHigh widget.Float // 0..8

	// Scroll lists.
	mainList widget.List

	// Status.
	lastTemplate *models.RmgTemplate
	statusMsg    string
	statusErr    bool
	statusTime   time.Time
}

func newState() *State {
	s := &State{
		gameMode: newSegmentGroup(gameModes),
		topology: newSegmentGroup(topologyLabels),
		victory:  newSegmentGroup(victoryLabels),
	}
	s.templateName.SingleLine = true
	s.templateName.SetText("Generated Template")
	s.outputPath.SingleLine = true
	if wd, err := os.Getwd(); err == nil {
		s.outputPath.SetText(wd)
	}
	s.mainList.Axis = layout.Vertical

	// Defaults matching MainWindow.xaml.cs initial values.
	s.mapSize.Value = float32(indexOf(mapSizes, 160)) / float32(len(mapSizes)-1)
	s.playerCnt.Value = mapRangeInv(4, 2, 8)
	s.chkRoads.Value = true
	s.chkPortals.Value = true
	s.chkBalancedZones.Value = true
	s.chkShowDescription.Value = true
	s.chkIncludeOptionsDesc.Value = true
	s.topology.Selected = 0 // Random
	s.victory.Selected = 0  // Standard
	s.gameMode.Selected = 0 // Classic

	s.sldGuardRandom.Value = mapRangeInv(0.05, 0, 0.5)
	s.sldConnPerZone.Value = mapRangeInv(2, 1, 4)
	s.sldNeutralLow.Value = mapRangeInv(1, 0, 8)
	s.sldNeutralMed.Value = mapRangeInv(1, 0, 8)
	s.sldNeutralHigh.Value = mapRangeInv(0, 0, 8)
	return s
}

func indexOf(xs []int, v int) int {
	for i, x := range xs {
		if x == v {
			return i
		}
	}
	return 0
}

func mapRange(v, lo, hi float32) float32 { return lo + v*(hi-lo) }

func mapRangeInv(actual, lo, hi float32) float32 {
	if hi == lo {
		return 0
	}
	r := (actual - lo) / (hi - lo)
	if r < 0 {
		return 0
	}
	if r > 1 {
		return 1
	}
	return r
}

// playerCount returns the integer player count from the slider.
func (s *State) playerCount() int {
	return int(math.Round(float64(mapRange(s.playerCnt.Value, 2, 8))))
}

// mapSizeValue returns the discrete map-size in the array.
func (s *State) mapSizeValue() int {
	idx := int(math.Round(float64(s.mapSize.Value) * float64(len(mapSizes)-1)))
	if idx < 0 {
		idx = 0
	}
	if idx >= len(mapSizes) {
		idx = len(mapSizes) - 1
	}
	return mapSizes[idx]
}

func (s *State) guardRandom() float64 { return float64(mapRange(s.sldGuardRandom.Value, 0, 0.5)) }
func (s *State) connectionsPerZone() int {
	return int(math.Round(float64(mapRange(s.sldConnPerZone.Value, 1, 4))))
}
func (s *State) neutralLow() int {
	return int(math.Round(float64(mapRange(s.sldNeutralLow.Value, 0, 8))))
}
func (s *State) neutralMed() int {
	return int(math.Round(float64(mapRange(s.sldNeutralMed.Value, 0, 8))))
}
func (s *State) neutralHigh() int {
	return int(math.Round(float64(mapRange(s.sldNeutralHigh.Value, 0, 8))))
}

// buildSettings converts the current widget state into a GeneratorSettings.
func (s *State) buildSettings() *models.GeneratorSettings {
	return &models.GeneratorSettings{
		TemplateName:                strings.TrimSpace(s.templateName.Text()),
		GameMode:                    gameModes[s.gameMode.Selected],
		PlayerCount:                 s.playerCount(),
		MapSize:                     mapSizeLabel(s.mapSizeValue()),
		Topology:                    topologyValues[s.topology.Selected],
		AllowRoads:                  s.chkRoads.Value,
		AllowPortals:                s.chkPortals.Value,
		AllowFootholds:              s.chkFootholds.Value,
		EnablePlayerIsolation:       s.chkPlayerIsolation.Value,
		EnableCityHold:              s.chkCityHold.Value,
		ShowDescription:             s.chkShowDescription.Value,
		IncludeOptionsInDescription: s.chkIncludeOptionsDesc.Value,
		AdvancedSettings: &models.AdvancedSettings{
			GuardRandomization:     s.guardRandom(),
			ConnectionCountPerZone: s.connectionsPerZone(),
			NeutralZoneLowCount:    s.neutralLow(),
			NeutralZoneMediumCount: s.neutralMed(),
			NeutralZoneHighCount:   s.neutralHigh(),
		},
		GameEndConditions: &models.GameEndConditions{
			EnableClassicVictory: s.victory.Selected == 0,
			EnableCityHold:       s.chkCityHold.Value || s.victory.Selected == 2,
			EnableTournaments:    s.victory.Selected == 3,
		},
	}
}

// generate runs the template generator and stores the result.
func (s *State) generate() {
	settings := s.buildSettings()
	if settings.TemplateName == "" {
		s.setStatus("Template name is required.", true)
		return
	}
	tmpl, err := generator.Generate(settings)
	if err != nil {
		s.setStatus(fmt.Sprintf("Generation failed: %v", err), true)
		s.lastTemplate = nil
		return
	}
	s.lastTemplate = tmpl
	zones := 0
	conns := 0
	if len(tmpl.Variants) > 0 {
		zones = len(tmpl.Variants[0].Zones)
		conns = len(tmpl.Variants[0].Connections)
	}
	s.setStatus(fmt.Sprintf("Generated '%s' — %d zones, %d connections.", tmpl.Name, zones, conns), false)
}

// save writes the most recently generated template to disk.
func (s *State) save() {
	if s.lastTemplate == nil {
		s.setStatus("Nothing to save — generate a template first.", true)
		return
	}
	dir := strings.TrimSpace(s.outputPath.Text())
	if dir == "" {
		s.setStatus("Output directory is empty.", true)
		return
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		s.setStatus(fmt.Sprintf("Cannot create directory: %v", err), true)
		return
	}
	safeName := sanitizeFilename(s.lastTemplate.Name)
	if safeName == "" {
		safeName = "Generated_Template"
	}
	out := filepath.Join(dir, safeName+".rmg.json")
	data, err := json.MarshalIndent(s.lastTemplate, "", "  ")
	if err != nil {
		s.setStatus(fmt.Sprintf("Marshal failed: %v", err), true)
		return
	}
	if err := os.WriteFile(out, data, 0o644); err != nil {
		s.setStatus(fmt.Sprintf("Write failed: %v", err), true)
		return
	}
	s.setStatus("Saved to "+out, false)
}

func sanitizeFilename(name string) string {
	bad := []rune{'/', '\\', ':', '*', '?', '"', '<', '>', '|'}
	out := []rune(strings.TrimSpace(name))
	for i, r := range out {
		for _, b := range bad {
			if r == b {
				out[i] = '_'
			}
		}
	}
	return string(out)
}

func (s *State) setStatus(msg string, isErr bool) {
	s.statusMsg = msg
	s.statusErr = isErr
	s.statusTime = time.Now()
}

// Layout draws the entire main window contents.
func (s *State) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	// Process button clicks first.
	if s.btnGenerate.Clicked(gtx) {
		s.generate()
	}
	if s.btnSave.Clicked(gtx) {
		s.save()
	}

	fillBackground(gtx, colBackground)
	return layout.UniformInset(unit.Dp(12)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceEnd}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return s.layoutHeader(gtx, th) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions { return s.layoutBody(gtx, th) }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return s.layoutFooter(gtx, th) }),
		)
	})
}

func (s *State) layoutHeader(gtx layout.Context, th *material.Theme) layout.Dimensions {
	title := material.H6(th, "⚔  Olden Era — Template Generator")
	title.Color = colGold
	subtitle := material.Body2(th, "⚠  Work in progress — generated templates may contain bugs.")
	subtitle.Color = colWarnText
	subtitle.TextSize = unit.Sp(11)
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(title.Layout),
		layout.Rigid(layout.Spacer{Height: unit.Dp(2)}.Layout),
		layout.Rigid(subtitle.Layout),
	)
}

func (s *State) layoutBody(gtx layout.Context, th *material.Theme) layout.Dimensions {
	list := material.List(th, &s.mainList)
	sections := []layout.Widget{
		func(gtx layout.Context) layout.Dimensions { return s.sectionGeneral(gtx, th) },
		spacer(unit.Dp(10)),
		func(gtx layout.Context) layout.Dimensions { return s.sectionTopologyVictory(gtx, th) },
		spacer(unit.Dp(10)),
		func(gtx layout.Context) layout.Dimensions { return s.sectionOptions(gtx, th) },
		spacer(unit.Dp(10)),
		func(gtx layout.Context) layout.Dimensions { return s.sectionAdvanced(gtx, th) },
		spacer(unit.Dp(10)),
		func(gtx layout.Context) layout.Dimensions { return s.sectionStatus(gtx, th) },
	}
	return list.Layout(gtx, len(sections), func(gtx layout.Context, i int) layout.Dimensions {
		return sections[i](gtx)
	})
}

func spacer(h unit.Dp) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Spacer{Height: h}.Layout(gtx)
	}
}

func (s *State) sectionGeneral(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return borderedPanel(gtx, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return sectionHeader(gtx, th, "General") }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Template name", 180, func(gtx layout.Context) layout.Dimensions {
					return drawEditor(gtx, th, &s.templateName, "Generated Template")
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Game mode", 180, func(gtx layout.Context) layout.Dimensions {
					return s.gameMode.Layout(gtx, th)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Player count", 180, func(gtx layout.Context) layout.Dimensions {
					return sliderWithValue(gtx, th, &s.playerCnt, fmt.Sprintf("%d", s.playerCount()))
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				size := s.mapSizeValue()
				return labeledRow(gtx, th, "Map size", 180, func(gtx layout.Context) layout.Dimensions {
					return sliderWithValue(gtx, th, &s.mapSize, fmt.Sprintf("%d  (%s)", size, mapSizeLabel(size)))
				})
			}),
		)
	})
}

func (s *State) sectionTopologyVictory(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return borderedPanel(gtx, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return sectionHeader(gtx, th, "Layout & Victory") }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Topology", 180, func(gtx layout.Context) layout.Dimensions {
					return s.topology.Layout(gtx, th)
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return dimLabel(gtx, th, topologyDescription(s.topology.Selected))
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Victory condition", 180, func(gtx layout.Context) layout.Dimensions {
					return s.victory.Layout(gtx, th)
				})
			}),
		)
	})
}

func topologyDescription(i int) string {
	switch i {
	case 0:
		return "Zones are placed at random positions. Each zone connects to all zones that border it — no fixed structure."
	case 1:
		return "All zones are arranged in a circle. Each zone connects to the two zones next to it."
	case 2:
		return "All zones connect to a shared central hub. Players never border each other directly."
	case 3:
		return "Zones are connected in a straight line from one end to the other, with no wrap-around."
	}
	return ""
}

func (s *State) sectionOptions(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return borderedPanel(gtx, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return sectionHeader(gtx, th, "Map Generation Options") }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(checkRow(th, &s.chkRoads, "Generate roads")),
			layout.Rigid(checkRow(th, &s.chkPortals, "Always spawn portals")),
			layout.Rigid(checkRow(th, &s.chkFootholds, "Spawn remote footholds")),
			layout.Rigid(checkRow(th, &s.chkBalancedZones, "Balance zone placement (recommended)")),
			layout.Rigid(checkRow(th, &s.chkPlayerIsolation, "Connect players via neutrals only when possible")),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(checkRow(th, &s.chkCityHold, "Enable City Hold win condition")),
			layout.Rigid(layout.Spacer{Height: unit.Dp(4)}.Layout),
			layout.Rigid(checkRow(th, &s.chkShowDescription, "Include description in template")),
			layout.Rigid(checkRow(th, &s.chkIncludeOptionsDesc, "Include current options in description")),
		)
	})
}

func (s *State) sectionAdvanced(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return borderedPanel(gtx, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return sectionHeader(gtx, th, "Advanced") }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(10)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Guard randomization", 220, func(gtx layout.Context) layout.Dimensions {
					return sliderWithValue(gtx, th, &s.sldGuardRandom, fmt.Sprintf("%.2f", s.guardRandom()))
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Connections per zone", 220, func(gtx layout.Context) layout.Dimensions {
					return sliderWithValue(gtx, th, &s.sldConnPerZone, fmt.Sprintf("%d", s.connectionsPerZone()))
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Neutral zones — Low", 220, func(gtx layout.Context) layout.Dimensions {
					return sliderWithValue(gtx, th, &s.sldNeutralLow, fmt.Sprintf("%d", s.neutralLow()))
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Neutral zones — Medium", 220, func(gtx layout.Context) layout.Dimensions {
					return sliderWithValue(gtx, th, &s.sldNeutralMed, fmt.Sprintf("%d", s.neutralMed()))
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(6)}.Layout),

			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Neutral zones — High", 220, func(gtx layout.Context) layout.Dimensions {
					return sliderWithValue(gtx, th, &s.sldNeutralHigh, fmt.Sprintf("%d", s.neutralHigh()))
				})
			}),
		)
	})
}

func (s *State) sectionStatus(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return borderedPanel(gtx, unit.Dp(12), func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions { return sectionHeader(gtx, th, "Output") }),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return labeledRow(gtx, th, "Output directory", 180, func(gtx layout.Context) layout.Dimensions {
					return drawEditor(gtx, th, &s.outputPath, "Choose folder")
				})
			}),
			layout.Rigid(layout.Spacer{Height: unit.Dp(8)}.Layout),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				msg := s.statusMsg
				if msg == "" {
					msg = "Ready."
				}
				color := colTextDim
				if s.statusErr {
					color = colError
				} else if s.lastTemplate != nil {
					color = colGoldBright
				}
				lbl := material.Body2(th, msg)
				lbl.Color = color
				lbl.TextSize = unit.Sp(12)
				return lbl.Layout(gtx)
			}),
		)
	})
}

func (s *State) layoutFooter(gtx layout.Context, th *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceBetween, Alignment: layout.Middle}.Layout(gtx,
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return goldButton{Text: "⚔  Generate Template", Click: &s.btnGenerate}.Layout(gtx, th)
			})
		}),
		layout.Rigid(layout.Spacer{Width: unit.Dp(10)}.Layout),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(2)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return goldButton{Text: "💾  Save Template", Click: &s.btnSave, Disabled: s.lastTemplate == nil}.Layout(gtx, th)
			})
		}),
	)
}
