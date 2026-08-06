package drivers

import (
	"errors"
	"fmt"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers/handler_interfaces"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
)

const (
	autoRegenDebounce   = 300 * time.Millisecond
	configFileExtension = ".gen.json"
)

type State struct {
	handler    handler_interfaces.IGuiHandler
	fileSystem handler_interfaces.IFileSystemHandler

	innerState *models.EditorState

	currentPath string
	unsaved     bool

	outputPath   widget.Editor
	lastTemplate *entities.RmgTemplate
	// templateRevision counts every replacement of lastTemplate, letting the
	// preview cache detect a new template without comparing its contents.
	templateRevision uint64
	statusMsg        string
	statusErr        bool

	confirmExit bool
	// onExit closes the application window (injected via SetOnExit) so exit
	// flows through the normal Gio app.DestroyEvent path.
	onExit func()

	applyNextStateAt time.Time

	// dialogs renders modal dialogs (rule editors, pickers, the connection
	// editor) over the main UI.
	dialogs *DialogHost
}

func NewUIState(
	handler handler_interfaces.IGuiHandler,
	fileSystem handler_interfaces.IFileSystemHandler,
	findTemplateDir bool) *State {
	state := &State{
		handler:    handler,
		fileSystem: fileSystem,
		innerState: models.NewEditorState(handler),
	}
	state.outputPath.SingleLine = true
	state.dialogs = &DialogHost{}
	if !findTemplateDir {
		return state
	}

	templateDir, err := helpers.FindOldenEraTemplatesDir(false)
	if templateDir == "" {
		if errors.Is(err, common_errors.ErrTemplatesDirNotFound) {
			state.SetStatus("Game template directory not found, using fallback directory.", false)
		} else {
			state.SetStatus(fmt.Sprintf("Failed to find game template directory: %v", err), true)
		}

		templateDir = state.workingDirectory()
	}
	state.outputPath.SetText(templateDir)
	return state
}

func (this *State) GetStatus() (msg string, isErr bool) { return this.statusMsg, this.statusErr }

func (this *State) GetDialogHost() *DialogHost { return this.dialogs }

func (this *State) GetStateData() dtos.EditorStateDto { return this.innerState.GetCurrentState() }

func (this *State) GetCurrentPath() string { return this.currentPath }

func (this *State) IsUnsaved() bool { return this.unsaved }

func (this *State) GetLastTemplate() *entities.RmgTemplate { return this.lastTemplate }

func (this *State) GetTemplateRevision() uint64 { return this.templateRevision }

func (this *State) GetOutputPath() string { return this.outputPath.Text() }

func (this *State) GetOutputPathWidget(theme *material.Theme) layout.Widget {
	return widgets.NewTextboxWidget(theme, &this.outputPath, "Choose folder", true)
}

func (this *State) Reset() {
	this.innerState.ResetState()
	this.currentPath = ""
	this.unsaved = false
	this.clearGeneratedState()
	this.SetStatus("New settings file.", false)
}

func (this *State) UpdateState(updateFunc func(*dtos.EditorStateDto)) {
	this.innerState.UpdateCurrentState(updateFunc)
	if this.innerState.WasStateChanged() {
		this.unsaved = true
		// New edits invalidate a pending exit confirmation.
		this.confirmExit = false
	}
}

func (this *State) SetStatus(msg string, isErr bool) {
	this.statusMsg = msg
	this.statusErr = isErr
}

func (this *State) hasTemplateVariants() bool {
	return this.lastTemplate != nil && len(this.lastTemplate.Variants) > 0
}

// setLastTemplate is the only writer of lastTemplate, so templateRevision
// cannot drift away from the template the preview is showing.
func (this *State) setLastTemplate(template *entities.RmgTemplate) {
	this.lastTemplate = template
	this.templateRevision++
}
