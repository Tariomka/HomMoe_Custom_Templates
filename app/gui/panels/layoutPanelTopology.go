package panels

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

func (this *LayoutPanel) getTopologySectionWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Topology", []layout.Widget{
		widgets.NewLabeledRowWidget(theme, "Topology", constants.DefaultLabelWidth, this.topology.GetWidget(theme)),
		func(gtx layout.Context) layout.Dimensions {
			label := material.Caption(theme, this.getCurrentTopology().Description)
			label.Color = themes.ColorTextDim
			return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(constants.DefaultLabelWidth + 8)}.
				Layout(gtx, label.Layout)
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

			return widgets.NewSliderRowWidget(theme, "Max portal connections", constants.DefaultLabelWidth,
				&this.sldMaxPortals, utils.RoundedRangeFormatter(1, 32))(gtx)
		},
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPlayerIsolation,
			"Disallow direct player-to-player connections"),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkMatchPlayerFactions, "Match player castle factions"),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkFootholds, "Spawn remote footholds"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkFootholds.Value {
				return layout.Dimensions{}
			}

			return widgets.NewSliderRowWidget(theme, "Remote footholds", constants.DefaultLabelWidth,
				&this.sldRemoteFootholds, utils.RoundedRangeFormatter(0, 4))(gtx)
		},
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkAbandonedOutposts,
			"Spawn abandoned outposts alongside neutral castles"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkAbandonedOutposts.Value {
				return layout.Dimensions{}
			}

			return widgets.NewSliderRowWidget(theme, "Abandoned outposts", constants.DefaultLabelWidth,
				&this.sldAbandonedOutposts, utils.RoundedRangeFormatter(0, 4))(gtx)
		},
	})
}

func (this *LayoutPanel) getZoneSizesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Zone sizes", []layout.Widget{
		widgets.NewSliderRowWidget(theme, "Player zone size", constants.DefaultLabelWidth,
			&this.sldPlayerZoneSize, utils.MultiplierFormatter(0.5, 1.5)),
		widgets.NewSliderRowWidget(theme, "Neutral zone size", constants.DefaultLabelWidth,
			&this.sldNeutralZoneSize, utils.MultiplierFormatter(0.5, 1.5)),
		func(gtx layout.Context) layout.Dimensions {
			if this.state.GetStateData().Topology != config.TopologyHubAndSpoke {
				return layout.Dimensions{}
			}

			return widgets.NewSliderRowWidget(theme, "Hub zone size", constants.DefaultLabelWidth,
				&this.sldHubSize, utils.MultiplierFormatter(0.5, 1.5))(gtx)
		},
	})
}

func (this *LayoutPanel) getDifficultyAndDensityWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Difficulty & Density", []layout.Widget{
		widgets.NewSliderRowWidget(theme, "Resource density", constants.DefaultLabelWidth,
			&this.sldResourceDensity, utils.RoundedRangePercentFormatter(25, 200)),
		widgets.NewSliderRowWidget(theme, "Structure density", constants.DefaultLabelWidth,
			&this.sldStructureDensity, utils.RoundedRangePercentFormatter(25, 200)),
		widgets.NewSliderRowWidget(theme, "Neutral stack strength", constants.DefaultLabelWidth,
			&this.sldNeutralStack, utils.RoundedRangePercentFormatter(25, 200)),
		widgets.NewSliderRowWidget(theme, "Border guard strength", constants.DefaultLabelWidth,
			&this.sldBorderGuard, utils.RoundedRangePercentFormatter(25, 200)),
		widgets.NewSliderRowWidget(theme, "Guard randomization", constants.DefaultLabelWidth,
			&this.sldGuardRandom, utils.DenormalizeFormatter(0, 0.5)),
	})
}
