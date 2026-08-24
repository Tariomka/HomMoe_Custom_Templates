package fileService_test

import (
	"image"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/stretchr/testify/mock"
)

// mockFileRepository stands in for any IFileRepository so the service can be
// tested for the decisions it makes rather than for what lands on disk.
type mockFileRepository[T any] struct {
	mock.Mock

	// seed is what the caller had put in target before the load, which is the
	// only place the defaults-before-decode contract is observable.
	seed T
}

func (this *mockFileRepository[T]) Load(filePath string, target *T) error {
	this.seed = *target
	arguments := this.Called(filePath)
	if loaded, ok := arguments.Get(0).(T); ok {
		*target = loaded
	}

	return arguments.Error(1)
}

func (this *mockFileRepository[T]) Save(directory string, filename string, entity T) (string, error) {
	arguments := this.Called(directory, filename, entity)

	return arguments.String(0), arguments.Error(1)
}

type serviceMocks struct {
	editorState *mockFileRepository[editor_state.EditorState]
	template    *mockFileRepository[template.RmgTemplate]
	preview     *mockFileRepository[image.RGBA]
	mapper      mappers.IEditorStateEntityMapper
}

func newServiceWithMocks() (file_service.IFileService, serviceMocks) {
	mocks := serviceMocks{
		editorState: &mockFileRepository[editor_state.EditorState]{},
		template:    &mockFileRepository[template.RmgTemplate]{},
		preview:     &mockFileRepository[image.RGBA]{},
		mapper:      mappers.NewEditorStateEntityMapper(),
	}

	service := file_service.NewFileService(mocks.editorState, mocks.template, mocks.preview, mocks.mapper)

	return service, mocks
}
