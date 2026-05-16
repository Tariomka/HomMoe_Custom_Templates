package components

import (
	"strings"

	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/components/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/services"
)

type Toolbar struct {
	buttonReset  widget.Clickable
	buttonOpen   widget.Clickable
	buttonSave   widget.Clickable
	buttonSaveAs widget.Clickable

	state *State
}

func NewToolbar(state *State) *Toolbar {
	return &Toolbar{
		state: state,
	}
}

func (this *Toolbar) GetWidget(theme *material.Theme) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		row := layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}
		return row.Layout(gtx,
			layout.Rigid(widgets.NewButtonWidget(theme, "📄 New", &this.buttonReset, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "📂 Open…", &this.buttonOpen, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "💾 Save", &this.buttonSave, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
			layout.Rigid(widgets.NewButtonWidget(theme, "💾 Save As…", &this.buttonSaveAs, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(12)),
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{Left: unit.Dp(12)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					path := this.currentPath
					if path == "" {
						path = "(unsaved)"
					}
					if this.dirty {
						path += " *"
					}
					label := material.Body2(theme, "File: "+path)
					label.Color = colTextDim
					label.TextSize = unit.Sp(11)
					label.MaxLines = 1
					label.Truncator = "…"
					label.Alignment = text.End
					return label.Layout(gtx)
				})
			}),
		)
	}
}

func (this *Toolbar) HandleClicks(gtx layout.Context) {
	if this.buttonReset.Clicked(gtx) {
		this.fileNew()
	}
	if this.buttonOpen.Clicked(gtx) {
		this.fileOpen()
	}
	if this.buttonSave.Clicked(gtx) {
		this.fileSave()
	}
	if this.buttonSaveAs.Clicked(gtx) {
		this.fileSaveAs()
	}
	// if this.buttonTemplates.Clicked(gtx) {
	// 	this.openTemplatesFolder()
	// }
}

// fileNew clears the in-memory model.
func (this *Toolbar) fileNew() {
	this.state.Reset()
	this.seedDefaultPlayerZoneContent()
	this.applyFromSettingsFile()
	this.setStatus("New settings file.", false)
}

// fileOpen presents a dialog and loads the chosen .gen.json file.
func (this *Toolbar) fileOpen() {
	path, err := utils.PickOpenFile("Open settings", "Settings (*.gen.json)|*.gen.json|All files|*.*", this.suggestDir())
	if err != nil {
		this.setStatus("Open dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	loaded, err := services.LoadSettingsFile(path)
	if err != nil {
		this.setStatus("Load failed: "+err.Error(), true)
		return
	}
	this.settingsFile = loaded
	this.currentPath = path
	this.dirty = false
	this.applyFromSettingsFile()
	this.setStatus("Loaded "+path, false)
}

// fileSave writes to the current path or prompts via Save As if none.
func (this *Toolbar) fileSave() {
	if this.currentPath == "" {
		this.fileSaveAs()
		return
	}
	if err := services.SaveSettingsFile(this.currentPath, this.captureToSettingsFile()); err != nil {
		this.setStatus("Save failed: "+err.Error(), true)
		return
	}
	this.dirty = false
	this.setStatus("Saved "+this.currentPath, false)
}

// fileSaveAs prompts for a destination path then writes the settings file.
func (this *Toolbar) fileSaveAs() {
	defaultName := services.SanitizeFilename(strings.TrimSpace(this.templateName.Text())) + ".gen.json"
	path, err := utils.PickSaveFile("Save settings as", "Settings (*.gen.json)|*.gen.json", this.suggestDir(), defaultName)
	if err != nil {
		this.setStatus("Save dialog failed: "+err.Error(), true)
		return
	}
	if path == "" {
		return
	}
	if err := services.SaveSettingsFile(path, this.captureToSettingsFile()); err != nil {
		this.setStatus("Save failed: "+err.Error(), true)
		return
	}
	this.currentPath = path
	this.dirty = false
	this.setStatus("Saved "+path, false)
}
