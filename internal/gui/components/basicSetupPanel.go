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
		gameMode:     content.NewSegmentButtonGroup([]string{"Classic", "SingleHero"}),
		topology: content.NewDropdownSelector(func() []string {
			labels := make([]string, len(constants.Topologies))
			for _, topology := range constants.Topologies {
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
				label := material.Body2(theme, this.getCurrentTopology().Description)
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

func (this *BasicSetupPanel) LoadFromState() {
	settings := this.state.GetSettingsFile()
	this.templateName.SetText(settings.TemplateName)
	this.gameMode.SetSelectedIndex(0)
	this.playerCount.Value = utils.Normalize(float32(settings.PlayerCount), 2, 8)
	this.checkExperimentalSizes.Value = settings.ExperimentalMapSizes
	this.mapSizeSlider.Value = utils.MapSizeToSlider(settings.MapSize, this.checkExperimentalSizes.Value)
	this.topology.SelectByName(constants.GetTopologyDescriptor(settings.Topology).Label)
}

func (this *BasicSetupPanel) SaveToState() {
	// TODO: check `.Update(gtx)` and on true update the value
	this.state.UpdateState(func(settings *models.SettingsFile) {
		settings.TemplateName = strings.TrimSpace(this.templateName.Text())
		settings.PlayerCount = int(utils.RoundHalfAway(float64(utils.Denormalize(this.playerCount.Value, 2, 8))))
		settings.MapSize = utils.SliderToMapSize(this.mapSizeSlider.Value, this.checkExperimentalSizes.Value)
		settings.ExperimentalMapSizes = this.checkExperimentalSizes.Value
		settings.Topology = this.getCurrentTopology().Type
	})
}

func (this *BasicSetupPanel) getCurrentTopology() models.TopologyDescriptor {
	return constants.Topologies[this.topology.GetSelectedIndex()]
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
