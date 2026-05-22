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
	templateName      widget.Editor
	playerCount       widget.Float
	mapSize           *content.DropdownSelector
	checkMoreMapSizes widget.Bool

	gameMode    *content.SegmentButtonGroup
	sldHeroMin  widget.Float
	sldHeroMax  widget.Float
	sldHeroIncr widget.Float

	topology *content.DropdownSelector

	scroll widget.List

	state *State
}

func NewBasicSetupPanel(state *State) *BasicSetupPanel {
	panel := &BasicSetupPanel{
		templateName: widget.Editor{SingleLine: true},
		gameMode:     content.NewSegmentButtonGroup([]string{"Classic", "SingleHero"}),
		topology: content.NewDropdownSelector(func() []string {
			labels := make([]string, 0)
			for _, topology := range constants.Topologies {
				labels = append(labels, topology.Label)
			}
			return labels
		}()),
		mapSize: content.NewDropdownSelector(func() []string {
			labels := make([]string, 0)
			if state.GetSettingsFile().ExperimentalMapSizes {
				for _, mapSize := range constants.AllMapSizes {
					labels = append(labels, mapSize.Label)
				}
			} else {
				for _, mapSize := range constants.BaseMapSizes {
					labels = append(labels, mapSize.Label)
				}
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
		func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
				layout.Flexed(0.5, this.getTemplateSectionWidget(theme)),
				layout.Rigid(widgets.NewHorizontalSpacerWidget(16)),
				layout.Flexed(0.5, this.getMapSectionWidget(theme)))
		},
		this.getTopologySectionWidget(theme),
	}
	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.scroll).
			Layout(gtx, len(widgetsList), func(gtx layout.Context, index int) layout.Dimensions { return widgetsList[index](gtx) })
	}
}

func (this *BasicSetupPanel) LoadFromState() {
	settings := this.state.GetSettingsFile()

	this.templateName.SetText(settings.TemplateName)
	this.playerCount.Value = utils.Normalize(float32(settings.PlayerCount), 2, 8)
	this.mapSize.SelectByName(constants.GetMapSize(settings.MapSize).Label)
	this.checkMoreMapSizes.Value = settings.ExperimentalMapSizes

	this.gameMode.SetSelectedIndex(0)
	this.sldHeroMin.Value = utils.Normalize(float32(settings.HeroCountMin), 1, 12)
	this.sldHeroMax.Value = utils.Normalize(float32(settings.HeroCountMax), 1, 12)
	this.sldHeroIncr.Value = utils.Normalize(float32(settings.HeroCountIncrement), 1, 10)

	this.topology.SelectByName(constants.GetTopologyDescriptor(settings.Topology).Label)
}

func (this *BasicSetupPanel) SaveToState() {
	// TODO: check `.Update(gtx)` and on true update the value
	this.state.UpdateState(func(settings *models.SettingsFile) {
		settings.TemplateName = strings.TrimSpace(this.templateName.Text())
		settings.PlayerCount = int(utils.RoundHalfAway(float64(utils.Denormalize(this.playerCount.Value, 2, 8))))

		settings.GameMode = constants.GameModes[this.gameMode.GetSelectedIndex()]
		settings.HeroCountMin = utils.RoundedRange(this.sldHeroMin.Value, 1, 12)
		settings.HeroCountMax = max(utils.RoundedRange(this.sldHeroMax.Value, 1, 12), settings.HeroCountMin)
		settings.HeroCountIncrement = utils.RoundedRange(this.sldHeroIncr.Value, 1, 10)

		settings.MapSize = this.getCurrentMapSize().Size
		settings.ExperimentalMapSizes = this.checkMoreMapSizes.Value
		settings.Topology = this.getCurrentTopology().Type
	})
}

func (this *BasicSetupPanel) getTemplateSectionWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Template", []layout.Widget{
		widgets.NewLabeledRowWidget(
			theme, "Template name", 160, widgets.NewTextboxWidget(theme, &this.templateName, "Enter template name")),
		widgets.NewLabeledRowWidget(
			theme, "Players", 160,
			widgets.NewLabeledSlider(theme, &this.playerCount, fmt.Sprintf("%d", utils.RoundedRange(this.playerCount.Value, 2, 8)))),
		widgets.NewLabeledRowWidget(theme, "Map size", 160, func(gtx layout.Context) layout.Dimensions {
			return this.updateMapSizeSelector(gtx).Layout(gtx, theme)
		}),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.checkMoreMapSizes, "Allow non official larger map sizes (>240)"),
	})
}

func (this *BasicSetupPanel) getMapSectionWidget(theme *material.Theme) layout.Widget {
	labelWidth := 150
	widgetList := []layout.Widget{
		widgets.NewLabeledRowWidget(theme, "Game mode", labelWidth, func(gtx layout.Context) layout.Dimensions {
			return this.gameMode.Layout(gtx, theme)
		}),
	}
	if !this.isSingleHero() {
		widgetList = append(widgetList,
			widgets.NewLabeledRowWidget(
				theme, "Hero count min", labelWidth,
				widgets.NewLabeledSlider(theme, &this.sldHeroMin, fmt.Sprintf("%d", utils.RoundedRange(this.sldHeroMin.Value, 1, 12)))),
			widgets.NewLabeledRowWidget(
				theme, "Hero count max", labelWidth,
				widgets.NewLabeledSlider(theme, &this.sldHeroMax, fmt.Sprintf("%d", utils.RoundedRange(this.sldHeroMax.Value, 1, 12)))),
			widgets.NewLabeledRowWidget(
				theme, "Increment", labelWidth,
				widgets.NewLabeledSlider(theme, &this.sldHeroIncr, fmt.Sprintf("%d", utils.RoundedRange(this.sldHeroIncr.Value, 1, 10)))))
	}
	return widgets.NewSectionWidget(theme, "Hero Restrictions", widgetList)
}

func (this *BasicSetupPanel) getTopologySectionWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Topology", []layout.Widget{
		widgets.NewLabeledRowWidget(theme, "Topology", 160, func(gtx layout.Context) layout.Dimensions {
			return this.topology.Layout(gtx, theme)
		}),
		func(gtx layout.Context) layout.Dimensions {
			label := material.Body2(theme, this.getCurrentTopology().Description)
			label.Color = themes.ColorTextDim
			label.TextSize = unit.Sp(12)
			return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(168)}.Layout(gtx, label.Layout)
		},
	})
}

func (this *BasicSetupPanel) updateMapSizeSelector(gtx layout.Context) *content.DropdownSelector {
	if !this.checkMoreMapSizes.Update(gtx) {
		return this.mapSize
	}

	labels := make([]string, 0)
	for _, mapSize := range constants.GetMapSizes(this.checkMoreMapSizes.Value) {
		labels = append(labels, mapSize.Label)
	}
	this.mapSize.SetItems(labels)
	this.mapSize.SelectByName(constants.GetMapSize(this.state.GetSettingsFile().MapSize).Label)

	return this.mapSize
}

func (this *BasicSetupPanel) getCurrentTopology() constants.TopologyDescriptor {
	return constants.Topologies[this.topology.GetSelectedIndex()]
}

func (this *BasicSetupPanel) getCurrentMapSize() constants.MapSize {
	return constants.AllMapSizes[this.mapSize.GetSelectedIndex()]
}

func (this *BasicSetupPanel) isSingleHero() bool {
	return this.gameMode.GetSelectedIndex() == 1
}
