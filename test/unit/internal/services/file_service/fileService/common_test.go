package fileService_test

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/stretchr/testify/mock"
)

// mockFileRepository stands in for any IFileRepository so the service can be
// tested for the decisions it makes rather than for what lands on disk.
type mockFileRepository[T any] struct {
	mock.Mock
}

func (this *mockFileRepository[T]) Load(filePath string) (T, error) {
	arguments := this.Called(filePath)

	return arguments.Get(0).(T), arguments.Error(1)
}

func (this *mockFileRepository[T]) Save(directory string, filename string, entity T) (string, error) {
	arguments := this.Called(directory, filename, entity)

	return arguments.String(0), arguments.Error(1)
}

type serviceMocks struct {
	editorState *mockFileRepository[editor_state_model.EditorState]
	template    *mockFileRepository[entities.RmgTemplate]
	preview     *mockFileRepository[image.RGBA]
}

func newServiceWithMocks() (file_service.IFileService, serviceMocks) {
	mocks := serviceMocks{
		editorState: &mockFileRepository[editor_state_model.EditorState]{},
		template:    &mockFileRepository[entities.RmgTemplate]{},
		preview:     &mockFileRepository[image.RGBA]{},
	}

	return file_service.NewFileService(mocks.editorState, mocks.template, mocks.preview), mocks
}
