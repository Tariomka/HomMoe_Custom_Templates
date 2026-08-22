package fileService_test

import (
	"errors"
	"image"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenSavedTemplateNameNeedsSanitizing_ForwardsItUnchangedToTheRepository(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "a/b:c"}
	mocks.template.On("Save", "out", "a/b:c", rmgTemplate).Return("written", nil)

	// Act
	_, err := service.SaveTemplateWithPreview("out", &rmgTemplate, nil)

	// Assert
	require.NoError(t, err)
	mocks.template.AssertCalled(t, "Save", "out", "a/b:c", rmgTemplate)
}

func TestWhenTemplateIsSaved_ReturnsThePathTheRepositoryWrote(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	expectedPath := filepath.Join("out", "T.rmg.json")
	mocks.template.On("Save", "out", "T", rmgTemplate).Return(expectedPath, nil)

	// Act
	actualPath, err := service.SaveTemplateWithPreview("out", &rmgTemplate, nil)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)
}

func TestWhenPreviewImageIsNil_DoesNotSaveAPreview(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	mocks.template.On("Save", "out", "T", rmgTemplate).Return("written", nil)

	// Act
	_, err := service.SaveTemplateWithPreview("out", &rmgTemplate, nil)

	// Assert
	require.NoError(t, err)
	mocks.preview.AssertNotCalled(t, "Save")
}

func TestWhenPreviewImageIsGiven_SavesItBesideTheTemplateUnderTheSameName(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	previewImage := newPreviewImage()
	mocks.template.On("Save", "out", "T", rmgTemplate).Return("written", nil)
	mocks.preview.On("Save", "out", "T", *previewImage).Return("preview", nil)

	// Act
	_, err := service.SaveTemplateWithPreview("out", &rmgTemplate, previewImage)

	// Assert
	require.NoError(t, err)
	mocks.preview.AssertCalled(t, "Save", "out", "T", *previewImage)
}

func TestWhenTemplateCannotBeSaved_DoesNotSaveAPreview(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	mocks.template.On("Save", "out", "T", rmgTemplate).Return("", errors.New("disk full"))

	// Act
	_, _ = service.SaveTemplateWithPreview("out", &rmgTemplate, newPreviewImage())

	// Assert
	mocks.preview.AssertNotCalled(t, "Save")
}

func TestWhenTemplateCannotBeSaved_ReturnsNoPath(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	mocks.template.On("Save", "out", "T", rmgTemplate).Return("", errors.New("disk full"))

	// Act
	actualPath, _ := service.SaveTemplateWithPreview("out", &rmgTemplate, nil)

	// Assert
	assert.Empty(t, actualPath)
}

func TestWhenTemplateCannotBeSaved_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	expectedError := errors.New("disk full")
	mocks.template.On("Save", "out", "T", rmgTemplate).Return("", expectedError)

	// Act
	_, err := service.SaveTemplateWithPreview("out", &rmgTemplate, nil)

	// Assert
	assert.ErrorIs(t, err, expectedError)
}

func TestWhenPreviewCannotBeSaved_StillReturnsTheTemplatePath(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	previewImage := newPreviewImage()
	expectedPath := filepath.Join("out", "T.rmg.json")
	mocks.template.On("Save", "out", "T", rmgTemplate).Return(expectedPath, nil)
	mocks.preview.On("Save", "out", "T", *previewImage).Return("", errors.New("no space"))

	// Act
	actualPath, _ := service.SaveTemplateWithPreview("out", &rmgTemplate, previewImage)

	// Assert
	assert.Equal(t, expectedPath, actualPath)
}

func TestWhenPreviewCannotBeSaved_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	service, mocks := newServiceWithMocks()
	rmgTemplate := entities.RmgTemplate{Name: "T"}
	previewImage := newPreviewImage()
	expectedError := errors.New("no space")
	mocks.template.On("Save", "out", "T", rmgTemplate).Return("written", nil)
	mocks.preview.On("Save", "out", "T", *previewImage).Return("", expectedError)

	// Act
	_, err := service.SaveTemplateWithPreview("out", &rmgTemplate, previewImage)

	// Assert
	assert.ErrorIs(t, err, expectedError)
}

func newPreviewImage() *image.RGBA { return image.NewRGBA(image.Rect(0, 0, 16, 16)) }
