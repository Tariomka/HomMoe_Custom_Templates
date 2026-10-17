package test_helpers

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// FileServiceMock is a testify mock of file_service.IFileService, used to
// unit-test collaborators without touching the file system.
type FileServiceMock struct {
	mock.Mock
}

// LoadSettingsFile returns the configured state and error, plus warnings when a
// third return value is configured.
func (this *FileServiceMock) LoadSettingsFile(filePath string) (*editor_state_model.EditorState, []string, error) {
	arguments := this.Called(filePath)
	state, _ := arguments.Get(0).(*editor_state_model.EditorState)
	var warnings []string
	if len(arguments) > 2 {
		warnings, _ = arguments.Get(2).([]string)
	}
	return state, warnings, arguments.Error(1)
}

func (this *FileServiceMock) SaveSettings(
	filePath string,
	editorState *editor_state_model.EditorState,
) (string, error) {
	arguments := this.Called(filePath, editorState)
	return arguments.String(0), arguments.Error(1)
}

func (this *FileServiceMock) SaveTemplateWithPreview(
	directory string,
	generated *template_model.Template,
	previewImage *image.RGBA) (string, error) {
	arguments := this.Called(directory, generated, previewImage)
	return arguments.String(0), arguments.Error(1)
}
