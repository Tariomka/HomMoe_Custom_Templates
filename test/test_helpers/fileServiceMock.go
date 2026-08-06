package test_helpers

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/mock"
)

// FileServiceMock is a testify mock of file_service.IFileService, used to
// unit-test collaborators without touching the file system.
type FileServiceMock struct {
	mock.Mock
}

func (this *FileServiceMock) LoadSettingsFile(filePath string) (*dtos.EditorStateDto, error) {
	arguments := this.Called(filePath)
	state, _ := arguments.Get(0).(*dtos.EditorStateDto)
	return state, arguments.Error(1)
}

func (this *FileServiceMock) SaveSettings(filePath string, editorState *dtos.EditorStateDto) (string, error) {
	arguments := this.Called(filePath, editorState)
	return arguments.String(0), arguments.Error(1)
}

func (this *FileServiceMock) SaveTemplateWithPreview(
	directory string,
	template *entities.RmgTemplate,
	previewImage *image.RGBA) (string, error) {
	arguments := this.Called(directory, template, previewImage)
	return arguments.String(0), arguments.Error(1)
}
