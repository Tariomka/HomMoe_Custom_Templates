package templateHandler_test

import (
	"errors"
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenThereIsNoTemplateToSave_ReturnsNothingToSaveError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()

	// Act
	_, err := fixture.handler.SaveTemplate(dtos.TemplateSaveDto{OutputPath: gofakeit.Word()})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNothingToSave)
}

func TestWhenTemplateOutputPathIsEmpty_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()

	// Act
	_, err := fixture.handler.SaveTemplate(dtos.TemplateSaveDto{Template: &template_model.Template{}})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenTemplateOutputPathIsWhitespaceOnly_ReturnsNoOutputPathError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()

	// Act
	_, err := fixture.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   &template_model.Template{},
		OutputPath: "  \t ",
	})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNoOutputPath)
}

func TestWhenTemplateOutputPathIsPadded_SavesToTheTrimmedPath(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	outputPath := gofakeit.Word()
	template := &template_model.Template{}
	templateEntity := new(mappers.NewTemplateMapper().ToEntity(*template))
	fixture.previewGenerator.On("CreatePreviewImage", mock.Anything, mock.Anything).Return(nil)
	fixture.fileService.On("SaveTemplateWithPreview", outputPath, templateEntity, mock.Anything).
		Return(gofakeit.Word(), nil)

	// Act
	_, _ = fixture.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   template,
		OutputPath: " " + outputPath + " ",
	})

	// Assert
	fixture.fileService.AssertCalled(t, "SaveTemplateWithPreview", outputPath, templateEntity, (*image.RGBA)(nil))
}

func TestWhenPreviewIsRendered_SavesItAlongsideTheTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	template := &template_model.Template{}
	templateEntity := new(mappers.NewTemplateMapper().ToEntity(*template))
	previewImage := image.NewRGBA(image.Rect(0, 0, 1, 1))
	fixture.previewGenerator.On("CreatePreviewImage", template, config.TopologyChain).Return(previewImage)
	fixture.fileService.On("SaveTemplateWithPreview", mock.Anything, mock.Anything, mock.Anything).
		Return(gofakeit.Word(), nil)

	// Act
	_, _ = fixture.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   template,
		Topology:   config.TopologyChain,
		OutputPath: gofakeit.Word(),
	})

	// Assert
	fixture.fileService.AssertCalled(t, "SaveTemplateWithPreview", mock.Anything, templateEntity, previewImage)
}

func TestWhenTemplateIsSaved_ReturnsTheWrittenPath(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	expectedPath := gofakeit.Word() + ".rmg.json"
	fixture.previewGenerator.On("CreatePreviewImage", mock.Anything, mock.Anything).Return(nil)
	fixture.fileService.On("SaveTemplateWithPreview", mock.Anything, mock.Anything, mock.Anything).
		Return(expectedPath, nil)

	// Act
	writtenPath, err := fixture.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   &template_model.Template{},
		OutputPath: gofakeit.Word(),
	})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedPath, writtenPath)
}

func TestWhenTemplateCannotBeSaved_PropagatesTheError(t *testing.T) {
	t.Parallel()
	// Arrange
	fixture := newTemplateHandlerFixture()
	expectedError := errors.New(gofakeit.Sentence(3))
	fixture.previewGenerator.On("CreatePreviewImage", mock.Anything, mock.Anything).Return(nil)
	fixture.fileService.On("SaveTemplateWithPreview", mock.Anything, mock.Anything, mock.Anything).
		Return("", expectedError)

	// Act
	_, err := fixture.handler.SaveTemplate(dtos.TemplateSaveDto{
		Template:   &template_model.Template{},
		OutputPath: gofakeit.Word(),
	})

	// Assert
	assert.ErrorIs(t, err, expectedError)
}
