package editor

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/panels"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
)

type Window struct {
	state *drivers.State

	tabs        []*drivers.Tab
	selectedTab int
	// tabChildren caches the FlexChild wrappers for the static tab strip;
	// built on the first frame (the theme does not exist yet in NewWindow)
	// and reused - each child reads live tab state at layout time.
	tabChildren []layout.FlexChild

	toolbar      *Toolbar
	previewPanel *panels.PreviewPanel
}

func NewWindow(backend handler_interfaces.IGuiHandler) *Window {
	window := Window{state: drivers.NewUIStateWithBackend(backend)}
	window.toolbar = NewToolbar(window.state, window.load)
	window.tabs = []*drivers.Tab{
		drivers.NewTab("General", panels.NewGeneralPanel(window.state)),
		drivers.NewTab("Layout & Zones", panels.NewLayoutPanel(window.state, backend, backend, backend)),
		drivers.NewTab("Bonuses & Bans", panels.NewBonusesPanel(window.state)),
	}
	window.previewPanel = panels.NewPreviewPanel(window.state, backend)
	window.tabs[0].SetSelected(true)
	return &window
}

// SetOnExit installs the callback the editor uses to close the application
// window when the user exits (see drivers.State.Exit).
func (this *Window) SetOnExit(onExit func()) { this.state.SetOnExit(onExit) }

func (this *Window) Layout(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	this.save()
	if redrawAt, scheduleRedraw := this.state.AutoRegenerate(gtx.Now); scheduleRedraw {
		gtx.Execute(op.InvalidateCmd{At: redrawAt})
	}
	this.handleClicks(gtx)

	paint.FillShape(gtx.Ops, themes.ColorsBase.Background, clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Op())

	return layout.Stack{}.Layout(gtx,
		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(10)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(this.toolbar.GetWidget(theme)),
					layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
					layout.Rigid(this.getTabsWidget(gtx, theme)),
					layout.Flexed(1, this.getPanelsWidget(theme)))
			})
		}),
		layout.Expanded(this.state.GetDialogHost().GetActiveDialogWidget(theme)),
	)
}

func (this *Window) getTabsWidget(gtx layout.Context, theme *material.Theme) layout.Widget {
	this.updateTabs(gtx)

	if this.tabChildren == nil {
		this.tabChildren = make([]layout.FlexChild, 0, len(this.tabs))
		for _, tab := range this.tabs {
			this.tabChildren = append(this.tabChildren, layout.Rigid(tab.GetWidget(theme)))
		}
	}
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, this.tabChildren...)
	}
}

func (this *Window) getPanelsWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(1, widgets.NewPanelWidget(unit.Dp(0), func(gtx layout.Context) layout.Dimensions {
				return layout.UniformInset(constants.DefaultPaddingLarge).
					Layout(gtx, this.getSelectedPanelWidget(theme))
			})),
			widgets.NewDefaultWidgetSpacer(),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				gtx.Constraints.Min.X = gtx.Dp(constants.DefaultPreviewWidthMinimum)
				gtx.Constraints.Max.X = gtx.Dp(constants.DefaultPreviewWidthMaximum)
				return this.previewPanel.GetPanelWidget(theme)(gtx)
			}))
	}
}

func (this *Window) getSelectedPanelWidget(theme *material.Theme) layout.Widget {
	return this.tabs[this.selectedTab].GetPanelWidget(theme)
}

func (this *Window) updateTabs(gtx layout.Context) {
	for i, tab := range this.tabs {
		if tab.IsTabClicked(gtx) && this.selectedTab != i {
			this.selectedTab = i
			this.updateTabSelection()
		}
	}
}

func (this *Window) updateTabSelection() {
	for i, tab := range this.tabs {
		tab.SetSelected(this.selectedTab == i)
	}
}

func (this *Window) handleClicks(gtx layout.Context) {
	this.toolbar.HandleClicks(gtx)
	this.previewPanel.HandleClicks(gtx)
}

func (this *Window) save() {
	for _, tab := range this.tabs {
		tab.SaveToState()
	}
}

func (this *Window) load() {
	for _, tab := range this.tabs {
		tab.LoadFromState()
	}
}
