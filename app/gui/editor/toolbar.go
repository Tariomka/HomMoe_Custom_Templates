package editor

import (
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/drivers"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

type Toolbar struct {
	buttonReset  widget.Clickable
	buttonOpen   widget.Clickable
	buttonSave   widget.Clickable
	buttonSaveAs widget.Clickable
	buttonExit   widget.Clickable

	resetCallback func()
	state         *drivers.State
}

func NewToolbar(state *drivers.State, resetCallback func()) *Toolbar {
	return &Toolbar{
		state:         state,
		resetCallback: resetCallback,
	}
}

func (this *Toolbar) GetWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Flexed(0.5, this.getButtonsWidget(theme)),
			widgets.NewDefaultWidgetSpacer(),
			layout.Rigid(widgets.NewTitleBarWidget(theme, "Heroes: Olden Era - Custom Templates")),
			widgets.NewDefaultWidgetSpacer(),
			layout.Flexed(0.5, this.getStateStatusWidget(theme)),
		)
	}
}

func (this *Toolbar) HandleClicks(gtx layout.Context) {
	if this.buttonReset.Clicked(gtx) {
		this.state.Reset()
		this.resetCallback()
	}
	if this.buttonOpen.Clicked(gtx) {
		this.state.Load(this.resetCallback)
	}
	if this.buttonSave.Clicked(gtx) {
		this.state.Save()
	}
	if this.buttonSaveAs.Clicked(gtx) {
		this.state.SaveAs(this.state.GetTemplateName())
	}
	if this.buttonExit.Clicked(gtx) {
		this.state.Exit()
	}
}

func (this *Toolbar) getButtonsWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
			layout.Rigid(widgets.NewButtonWidget(theme, "New", &this.buttonReset, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Load", &this.buttonOpen, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Save", &this.buttonSave, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Save As", &this.buttonSaveAs, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "Exit", &this.buttonExit, false)))
	}
}

func (this *Toolbar) getStateStatusWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		path := this.state.GetCurrentPath()
		if path == "" {
			path = "(unsaved)"
		}
		if this.state.IsUnsaved() {
			path += " *"
		}
		return layout.Inset{Left: constants.DefaultPaddingLarge + 2}.Layout(gtx,
			widgets.NewLabelBuilder(theme).WithSizeDefault().WithText("File: "+path).WithColorDim().
				WithMaxLines(1).WithAlignment(text.End).Build)
	}
}
