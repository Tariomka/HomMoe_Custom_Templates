package dialogs

import (
	"errors"
	"fmt"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
)

func (this *FileExplorerDialog) getHeaderWidget(theme *material.Theme) layout.Widget {
	upDisabled := this.parentDir() == this.currentDir
	if this.pathEd.Text() != this.currentDir {
		if this.currentDir == "" {
			this.pathEd.SetText("This PC")
		} else {
			this.pathEd.SetText(this.currentDir)
		}
	}

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
			layout.Rigid(widgets.NewButtonWidget(theme, "← Back", &this.upBtn, upDisabled)),
			widgets.NewDefaultComponentSpacer(),
			layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.pathEd, "Current directory", true)),
			widgets.NewDefaultComponentSpacer(),
			layout.Rigid(widgets.NewToggleButtonWidget(theme, "Show hidden", &this.hiddenToggle, this.showHidden)),
		)
	}
}

func (this *FileExplorerDialog) getSaveRowWidget(theme *material.Theme) layout.Widget {
	if this.mode != modeSaveFile {
		return widgets.NewEmptyWidget()
	}

	hint := fmt.Sprintf("filename%s", saveFileSuffix)
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: constants.DefaultPadding}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
				layout.Rigid(widgets.NewLabelBigWidget(theme, "Save as:", themes.ColorsBase.TextDim)),
				widgets.NewDefaultComponentSpacer(),
				layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.filenameEd, hint, false)),
			)
		})
	}
}

func (this *FileExplorerDialog) getNewFolderRowWidget(theme *material.Theme) layout.Widget {
	if !this.canModify() || !this.newFolderActive {
		return widgets.NewEmptyWidget()
	}

	return func(gtx layout.Context) layout.Dimensions {
		return layout.Inset{Top: constants.DefaultPadding}.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return layout.Flex{Axis: layout.Horizontal, Alignment: layout.Middle}.Layout(gtx,
						layout.Rigid(widgets.NewLabelBigWidget(theme, "New folder:", themes.ColorsBase.TextDim)),
						widgets.NewDefaultComponentSpacer(),
						layout.Flexed(1, widgets.NewTextboxWidget(theme, &this.newFolderEd, "folder name", false)),
						widgets.NewDefaultComponentSpacer(),
						layout.Rigid(widgets.NewButtonWidget(theme, "Create", &this.createFolderBtn, false)),
						widgets.NewDefaultComponentSpacer(),
						layout.Rigid(widgets.NewButtonWidget(theme, "Cancel", &this.cancelFolderBtn, false)),
					)
				}),
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					if this.newFolderErr == "" {
						return layout.Dimensions{}
					}

					return layout.Inset{Top: constants.DefaultPaddingSmall - 2}.
						Layout(gtx, widgets.NewLabelBuilder(theme).WithSizeDefault().
							WithText(this.newFolderErr).WithColor(themes.ColorsBase.Error).WithMaxLines(2).Build)
				}),
			)
		})
	}
}

func (this *FileExplorerDialog) tryCreateFolder() {
	target, err := this.fileSystem.CreateDirectory(this.currentDir, this.newFolderEd.Text())
	switch {
	case errors.Is(err, common_errors.ErrDirectoryNameEmpty):
		this.newFolderErr = "Enter a folder name."
		return
	case errors.Is(err, common_errors.ErrDirectoryNameInvalid):
		this.newFolderErr = "Invalid folder name."
		return
	case err != nil:
		this.newFolderErr = err.Error()
		return
	}

	this.newFolderActive = false
	this.newFolderErr = ""
	this.loadDir(target)
}
