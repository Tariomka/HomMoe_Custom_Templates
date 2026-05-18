package components

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/content"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/generator"
)

type BasicSetupPanel struct {
	templateName           widget.Editor
	gameMode               *content.SegmentButtonGroup
	playerCount            widget.Float
	mapSizeSlider          widget.Float
	checkExperimentalSizes widget.Bool
	topology               *content.DropdownSelector

	scroll widget.List

	state *State
}

func NewBasicSetupPanel(state *State) *BasicSetupPanel {
	panel := &BasicSetupPanel{
		templateName: widget.Editor{SingleLine: true},
		gameMode:     content.NewSegmentButtonGroup(constants.GameModes),
		topology:     content.NewDropdownSelector(constants.TopologyLabels),
		state:        state,
	}
	panel.scroll.Axis = layout.Vertical
	panel.loadFromState()
	return panel
}

func (this *BasicSetupPanel) GetPanelWidget(theme *material.Theme) layout.Widget {
	widgetsList := []layout.Widget{
		widgets.NewSectionWidget(theme, "Template", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Template name", 160, widgets.NewTextboxWidget(theme, &this.templateName, "Enter template name")),
			widgets.NewLabeledRowWidget(theme, "Game mode", 160, func(gtx layout.Context) layout.Dimensions {
				return this.gameMode.Layout(gtx, theme)
			}),
		}),
		widgets.NewSectionWidget(theme, "Map", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Players", 160, widgets.NewLabeledSlider(theme, &this.playerCount, fmt.Sprintf("%d", utils.RoundedRange(this.playerCount.Value, 2, 8)))),
			widgets.NewLabeledRowWidget(theme, "Map size", 160, func(gtx layout.Context) layout.Dimensions {
				size := utils.SliderToMapSize(this.mapSizeSlider.Value, this.checkExperimentalSizes.Value)
				label := fmt.Sprintf("%d × %d  (%s)", size, size, mapSizeLabelInt(size))
				return widgets.NewLabeledSlider(theme, &this.mapSizeSlider, label)(gtx)
			}),
			widgets.NewLabeledCheckboxRowWidget(theme, &this.checkExperimentalSizes, "Allow experimental large map sizes (>240)"),
		}),
		widgets.NewSectionWidget(theme, "Topology", []layout.Widget{
			widgets.NewLabeledRowWidget(theme, "Topology", 160, func(gtx layout.Context) layout.Dimensions {
				return this.topology.Layout(gtx, theme)
			}),
			func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, topologyDescription(constants.TopologyValues[this.topology.GetSelectedIndex()]))
				label.Color = themes.ColorTextDim
				label.TextSize = unit.Sp(12)
				return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(168)}.Layout(gtx, label.Layout)
			},
		}),
	}
	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.scroll).Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions {
			return widgetsList[index](gtx)
		})
	}
}

func (this *BasicSetupPanel) loadFromState() {
	settings := this.state.GetSettingsFile()
	this.templateName.SetText(settings.TemplateName)
	this.gameMode.SetSelectedIndex(0)
	this.playerCount.Value = utils.Normalize(float32(settings.PlayerCount), 2, 8)
	this.checkExperimentalSizes.Value = settings.ExperimentalMapSizes
	this.mapSizeSlider.Value = utils.MapSizeToSlider(settings.MapSize, this.checkExperimentalSizes.Value)
	this.topology.SelectByName(topologyLabelFor(settings.Topology))
}

// TODO: check `.Update(gtx)` and on true update the value
func (this *BasicSetupPanel) saveToState() {
	this.state.UpdateState(func(settings *models.SettingsFile) {
		settings.TemplateName = strings.TrimSpace(this.templateName.Text())
		settings.PlayerCount = int(utils.RoundHalfAway(float64(utils.Denormalize(this.playerCount.Value, 2, 8))))
		settings.MapSize = utils.SliderToMapSize(this.mapSizeSlider.Value, this.checkExperimentalSizes.Value)
		settings.ExperimentalMapSizes = this.checkExperimentalSizes.Value
		settings.Topology = constants.TopologyValues[this.topology.GetSelectedIndex()]
	})
}

func topologyLabelFor(topology models.MapTopology) string {
	for i, value := range constants.TopologyValues {
		if value == topology {
			return constants.TopologyLabels[i]
		}
	}
	return constants.TopologyLabels[0]
}

// mapSizeLabelInt returns the short S/M/L/XL/H/G/C label for an integer size.
func mapSizeLabelInt(size int) string {
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
