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

			return widgets.NewLabeledRowWidget(theme, "Max portal connections", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(theme, &this.sldMaxPortals,
					utils.RoundedRangeString(this.sldMaxPortals.Value, 1, 32)))(gtx)
		},
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkPlayerIsolation,
			"Disallow direct player-to-player connections"),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkMatchPlayerFactions, "Match player castle factions"),
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkFootholds, "Spawn remote footholds"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkFootholds.Value {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(theme, "Remote footholds", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(theme, &this.sldRemoteFootholds,
					utils.RoundedRangeString(this.sldRemoteFootholds.Value, 0, 4)))(gtx)
		},
		widgets.NewLabeledCheckboxRowWidget(theme, &this.chkAbandonedOutposts,
			"Spawn abandoned outposts alongside neutral castles"),
		func(gtx layout.Context) layout.Dimensions {
			if !this.chkAbandonedOutposts.Value {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(theme, "Abandoned outposts", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(theme, &this.sldAbandonedOutposts,
					utils.RoundedRangeString(this.sldAbandonedOutposts.Value, 0, 4)))(gtx)
		},
		// TODO: Fuck this shit, it doesn't actually honor anything, the backend logic is lost somewhere,
		// need to remove this option completely from UI, Internal and any dtos.
		// Fable tried but its not worth it to have this here, and I don't want spaghetti to honor it.
		// widgets.NewLabeledRowWidget(theme, "Min neutrals between players", constants.DefaultLabelWidthLong,
		// 	widgets.NewLabeledSliderWidget(theme, &this.sldMinNeutralBetween,
		// 		utils.RoundedRangeString(this.sldMinNeutralBetween.Value, 0, 8))),
		// func(gtx layout.Context) layout.Dimensions {
		// 	label := material.Caption(theme,
		// 		"Honored for Ring, Circles, Chain, Hub & Spoke and Shared Web topologies when enough "+
		// 			"neutral zones exist and random portals are off; ignored otherwise.")
		// 	label.Color = themes.ColorTextDim
		// 	return layout.Inset{Top: unit.Dp(2), Left: unit.Dp(constants.DefaultLabelWidthLong + 8)}.
		// 		Layout(gtx, label.Layout)
		// },
	})
}

func (this *LayoutPanel) getZoneSizesWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Zone sizes", []layout.Widget{
		widgets.NewLabeledRowWidget(theme, "Player zone size", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(theme, &this.sldPlayerZoneSize,
				utils.MultiplierString(this.sldPlayerZoneSize.Value, 0.5, 1.5))),
		widgets.NewLabeledRowWidget(theme, "Neutral zone size", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(theme, &this.sldNeutralZoneSize,
				utils.MultiplierString(this.sldNeutralZoneSize.Value, 0.5, 1.5))),
		func(gtx layout.Context) layout.Dimensions {
			if this.state.GetStateData().Topology != config.TopologyHubAndSpoke {
				return layout.Dimensions{}
			}

			return widgets.NewLabeledRowWidget(theme, "Hub zone size", constants.DefaultLabelWidth,
				widgets.NewLabeledSliderWidget(theme, &this.sldHubSize,
					utils.MultiplierString(this.sldHubSize.Value, 0.5, 1.5)))(gtx)
		},
	})
}

func (this *LayoutPanel) getDifficultyAndDensityWidget(theme *material.Theme) layout.Widget {
	return widgets.NewSectionWidget(theme, "Difficulty & Density", []layout.Widget{
		widgets.NewLabeledRowWidget(theme, "Resource density", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(theme, &this.sldResourceDensity,
				utils.RoundedRangePercentString(this.sldResourceDensity.Value, 25, 200))),
		widgets.NewLabeledRowWidget(theme, "Structure density", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(theme, &this.sldStructureDensity,
				utils.RoundedRangePercentString(this.sldStructureDensity.Value, 25, 200))),
		widgets.NewLabeledRowWidget(theme, "Neutral stack strength", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(theme, &this.sldNeutralStack,
				utils.RoundedRangePercentString(this.sldNeutralStack.Value, 25, 200))),
		widgets.NewLabeledRowWidget(theme, "Border guard strength", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(theme, &this.sldBorderGuard,
				utils.RoundedRangePercentString(this.sldBorderGuard.Value, 25, 200))),
		widgets.NewLabeledRowWidget(theme, "Guard randomization", constants.DefaultLabelWidth,
			widgets.NewLabeledSliderWidget(theme, &this.sldGuardRandom,
				utils.DenormalizeString(this.sldGuardRandom.Value, 0, 0.5))),
	})
}
