package dialogs

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
)

// loadDir reads dir and replaces the cached listing. It is the single
// navigation primitive and is called only on open, on explicit navigation and
// on toggling hidden files - never per frame. On failure the previous listing
// is kept and the error is surfaced inline so navigation into an unreadable
// folder leaves the user where they were.
func (this *FileExplorerDialog) loadDir(dir string) {
	this.newFolderActive = false
	this.newFolderErr = ""
	this.overwriteActive = false
	this.listErr = ""
	this.saveErr = ""
	this.selectedPath = ""

	if dir == "" {
		this.currentDir = ""
		this.entries = this.fileSystem.ListRoots()
		this.rowClicks = map[string]*widget.Clickable{}
		this.resetScroll()
		return
	}

	entries, err := this.fileSystem.ListEntries(dir, this.filterSuffixes, this.showHidden)
	if err != nil {
		// Keep the current location; only adopt dir on the very first load so
		// the path bar is coherent when the initial directory is unreadable.
		if this.currentDir == "" {
			this.currentDir = dir
		}
		this.listErr = err.Error()
		return
	}

	this.currentDir = dir
	this.entries = entries
	this.rowClicks = map[string]*widget.Clickable{}
	this.resetScroll()
}

func (this *FileExplorerDialog) getListWidget(theme *material.Theme) layout.Widget {
	if len(this.entries) == 0 {
		message := "(empty folder)"
		if this.listErr != "" {
			message = "(unable to read folder)"
		}
		return func(gtx layout.Context) layout.Dimensions {
			return layout.Center.Layout(gtx, widgets.NewLabelBigWidget(theme, message, themes.ColorsBase.TextDim))
		}
	}

	return func(gtx layout.Context) layout.Dimensions {
		return material.List(theme, &this.list).Layout(gtx, len(this.entries),
			func(gtx layout.Context, index int) layout.Dimensions {
				return this.getEntryRowWidget(theme, this.entries[index])(gtx)
			})
	}
}

func (this *FileExplorerDialog) getEntryRowWidget(
	theme *material.Theme,
	entry models.DirectoryEntry) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		clk := this.clickFor(entry.Path)
		if clk.Clicked(gtx) {
			this.onEntryClicked(entry)
		}
		selected := !entry.IsDir && entry.Path == this.selectedPath
		badgeText := ""
		textColor := themes.ColorsBase.Text
		if entry.IsDir {
			badgeText = "DIR"
			textColor = themes.ColorsBase.AccentBright
		}
		return material.Clickable(gtx, clk, func(gtx layout.Context) layout.Dimensions {
			macro := op.Record(gtx.Ops)
			dims := layout.Inset{Top: unit.Dp(4), Bottom: unit.Dp(4), Left: unit.Dp(6), Right: unit.Dp(6)}.
				Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							gtx.Constraints.Min.X = gtx.Dp(unit.Dp(38))
							return widgets.NewStyledLabelWidget(
								theme, badgeText, themes.ColorsBase.Accent, font.Font{Weight: font.SemiBold})(gtx)
						}),
						layout.Flexed(1, widgets.NewLabelBuilder(theme).WithSizeBig().
							WithText(entry.Name).WithColor(textColor).WithMaxLines(1).Build))
				})
			call := macro.Stop()
			if selected {
				paint.FillShape(gtx.Ops, themes.ColorsBase.Selection, clip.Rect{Max: dims.Size}.Op())
			} else if clk.Hovered() {
				paint.FillShape(gtx.Ops, themes.ColorsBase.Hover, clip.Rect{Max: dims.Size}.Op())
			}
			call.Add(gtx.Ops)
			return dims
		})
	}
}

func (this *FileExplorerDialog) getErrorLineWidget(theme *material.Theme) layout.Widget {
	message := this.listErr
	if message == "" {
		message = this.saveErr
	}

	if message == "" {
		return widgets.NewEmptyWidget()
	}

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: constants.DefaultPaddingSmall}.
			Layout(gtx, widgets.NewLabelBuilder(theme).WithSizeDefault().
				WithText(message).WithColor(themes.ColorsBase.Error).WithMaxLines(2).Build)
	}
}
