package drivers

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/interfaces"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/models"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/widgets"
	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
)

const (
	autoRegenDebounce   = 300 * time.Millisecond
	configFileExtension = ".gen.json"
)

type State struct {
	handler interfaces.ITemplateHandler
	mapper  *mappers.GeneratorConfigMapper

	innerState *models.EditorState

	currentPath string
	unsaved     bool

	outputPath   widget.Editor
	lastTemplate *entities.RmgTemplate
	statusMsg    string
	statusErr    bool

	confirmExit bool
	// onExit closes the application window (injected via SetOnExit) so exit
	// flows through the normal Gio app.DestroyEvent path.
	onExit func()

	applyNextStateAt time.Time

	// dialogs renders modal dialogs (rule editors, pickers, the connection
	// editor) over the main UI.
	dialogs *DialogHost
}

func NewUIState() *State {
	state := NewUIStateWithHandler(handlers.NewGuiHandler())

	templateDir, err := helpers.FindOldenEraTemplatesDir(false)
	if templateDir == "" {
		if errors.Is(err, common_errors.ErrTemplatesDirNotFound) {
			state.SetStatus("Game template directory not found, using fallback directory.", false)
		} else {
			state.SetStatus(fmt.Sprintf("Failed to find game template directory: %v", err), true)
		}

		if workingDir, err := os.Getwd(); err == nil {
			templateDir = workingDir
		}
	}
	state.outputPath.SetText(templateDir)
	return state
}

// NewUIStateWithHandler builds a State around the given template handler
// without probing the disk for the game templates directory. Production code
// uses NewUIState; tests inject a mock handler here.
func NewUIStateWithHandler(handler interfaces.ITemplateHandler) *State {
	state := &State{
		handler:    handler,
		mapper:     mappers.NewConfigMapper(),
		innerState: models.NewEditorState(),
	}
	state.outputPath.SingleLine = true
	state.dialogs = &DialogHost{}
	return state
}

func (this *State) GetStatus() (msg string, isErr bool) { return this.statusMsg, this.statusErr }

func (this *State) GetDialogHost() *DialogHost { return this.dialogs }

func (this *State) GetStateData() dtos.EditorStateDto { return this.innerState.GetCurrentState() }

func (this *State) GetGeneratorConfig() *config.GeneratorConfig {
	return this.mapper.FromEditorState(this.innerState.GetCurrentState())
}

func (this *State) GetCurrentPath() string { return this.currentPath }

func (this *State) IsUnsaved() bool { return this.unsaved }

func (this *State) GetLastTemplate() *entities.RmgTemplate { return this.lastTemplate }

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
