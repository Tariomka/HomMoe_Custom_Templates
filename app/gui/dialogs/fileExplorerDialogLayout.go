package dialogs

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
)

// layoutContent stacks the header, the scrollable listing, an optional error
// line, the mode-specific input rows and the footer.
func (this *FileExplorerDialog) layoutContent(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return this.layoutHeader(gtx, theme)
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(8)),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return this.layoutList(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return this.layoutErrorLine(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return this.layoutSaveRow(gtx, theme)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return this.layoutNewFolderRow(gtx, theme)
		}),
		layout.Rigid(widgets.NewVerticalSpacerWidget(10)),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return this.layoutFooter(gtx, theme)
		}),
	)
}

func (this *FileExplorerDialog) layoutHeader(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	upDisabled := this.parentDir() == this.currentDir
	if this.pathEd.Text() != this.currentDir {
		if this.currentDir == "" {
			this.pathEd.SetText("This PC")
		} else {
			this.pathEd.SetText(this.currentDir)
		}
	}

	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
		layout.Rigid(widgets.NewButtonWidget(theme, "Up", &this.upBtn, upDisabled)),
		layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
		layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.pathEd, "Current directory", true)),
		layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
		layout.Rigid(widgets.NewToggleButtonWidget(theme, "Show hidden", &this.hiddenToggle, this.showHidden)),
	)
}

func (this *FileExplorerDialog) layoutList(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if len(this.entries) == 0 {
		message := "(empty folder)"
		if this.listErr != "" {
			message = "(unable to read folder)"
		}
		label := material.Body2(theme, message)
		label.Color = themes.ColorTextDim
		return layout.Center.Layout(gtx, label.Layout)
	}
	return material.List(theme, &this.list).Layout(gtx, len(this.entries), func(gtx layout.Context, index int) layout.Dimensions {
		return this.entryRow(gtx, theme, this.entries[index])
	})
}

func (this *FileExplorerDialog) entryRow(gtx layout.Context, theme *material.Theme, entry fileEntry) layout.Dimensions {
	clk := this.clickFor(entry.path)
	if clk.Clicked(gtx) {
		this.onEntryClicked(entry)
	}
	selected := !entry.isDir && entry.path == this.selectedPath
	return material.Clickable(gtx, clk, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		dims := layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(6), Right: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min.X = gtx.Dp(unit.Dp(38))
					badge := ""
					if entry.isDir {
						badge = "DIR"
					}
					label := material.Caption(theme, badge)
					label.Color = themes.ColorAccent
					label.Font = font.Font{Weight: font.SemiBold}
					return label.Layout(gtx)
				}),
				layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					label := material.Body2(theme, entry.name)
					label.Color = themes.ColorText
					if entry.isDir {
						label.Color = themes.ColorAccentBright
					}
					label.MaxLines = 1
					label.Truncator = "..."
					return label.Layout(gtx)
				}),
			)
		})
		call := macro.Stop()
		if selected {
			paint.FillShape(gtx.Ops, themes.ColorSelection, clip.Rect{Max: dims.Size}.Op())
		}
		call.Add(gtx.Ops)
		return dims
	})
}

func (this *FileExplorerDialog) layoutErrorLine(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if this.listErr == "" {
		return layout.Dimensions{}
	}
	return layout.Inset{Top: unit.Dp(6)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		label := material.Caption(theme, this.listErr)
		label.Color = themes.ColorError
		label.MaxLines = 2
		label.Truncator = "..."
		return label.Layout(gtx)
	})
}

func (this *FileExplorerDialog) layoutSaveRow(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if this.mode != modeSaveFile {
		return layout.Dimensions{}
	}
	return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, "Save as:")
				label.Color = themes.ColorTextDim
				return label.Layout(gtx)
			}),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
			layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.filenameEd, "filename"+saveFileSuffix, false)),
		)
	})
}

func (this *FileExplorerDialog) layoutNewFolderRow(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if !this.canModify() || !this.newFolderActive {
		return layout.Dimensions{}
	}
	return layout.Inset{Top: unit.Dp(8)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						label := material.Body2(theme, "New folder:")
						label.Color = themes.ColorTextDim
						return label.Layout(gtx)
					}),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
					layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.newFolderEd, "folder name", false)),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
					layout.Rigid(widgets.NewButtonWidget(theme, "Create", &this.createFolderBtn, false)),
					layout.Rigid(widgets.NewHorizontalSpacerWidget(6)),
					layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelFolderBtn, false)),
				)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				if this.newFolderErr == "" {
					return layout.Dimensions{}
				}
				return layout.Inset{Top: unit.Dp(4)}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					label := material.Caption(theme, this.newFolderErr)
					label.Color = themes.ColorError
					label.MaxLines = 2
					label.Truncator = "..."
					return label.Layout(gtx)
				})
			}),
		)
	})
}

func (this *FileExplorerDialog) layoutFooter(gtx layout.Context, theme *material.Theme) layout.Dimensions {
	if this.mode == modeSaveFile && this.overwriteActive {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
				label := material.Body2(theme, "File already exists. Overwrite?")
				label.Color = themes.ColorWarnText
				label.MaxLines = 1
				label.Truncator = "..."
				return label.Layout(gtx)
			}),
			layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.overwriteCancelBtn, false)),
			layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
			layout.Rigid(widgets.NewGoldButtonWidget(theme, "Overwrite", &this.overwriteConfirmBtn, false)),
		)
	}

	confirmLabel, showConfirm, confirmDisabled := this.confirmButtonState()
	cancelLabel := "Cancel"
	if this.mode == modeBrowse {
		cancelLabel = "Close"
	}

	children := make([]layout.FlexChild, 0, 5)
	if this.canModify() {
		children = append(children, layout.Rigid(widgets.NewButtonWidget(theme, "New Folder", &this.newFolderBtn, false)))
	}
	children = append(children, layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
		return layout.Dimensions{Size: image.Pt(gtx.Constraints.Min.X, 0)}
	}))
	children = append(children, layout.Rigid(widgets.NewButtonWidget(theme, cancelLabel, &this.cancelBtn, false)))
	if showConfirm {
		children = append(children,
			layout.Rigid(widgets.NewHorizontalSpacerWidget(8)),
			layout.Rigid(widgets.NewGoldButtonWidget(theme, confirmLabel, &this.confirmBtn, confirmDisabled)),
		)
	}
	return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx, children...)
}

// confirmButtonState returns the primary button's label, whether it is shown at
// all (browse mode has none) and whether it is currently disabled.
func (this *FileExplorerDialog) confirmButtonState() (label string, show bool, disabled bool) {
	switch this.mode {
	case modeOpenFile:
		return "Open", true, this.selectedPath == ""
	case modeSaveFile:
		return "Save", true, len(this.filenameEd.Text()) == 0
	case modePickFolder:
		return "Select This Folder", true, this.currentDir == ""
	default:
		return "", false, false
	}
}
