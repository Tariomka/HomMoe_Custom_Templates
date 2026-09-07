package previewHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWhenRequestIsEmpty_LaysOutNoTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{})

	// Assert
	layoutService.AssertCalled(t, "BuildPreviewLayout", (*template_model.Template)(nil), mock.Anything, mock.Anything)
}

func TestWhenTemplateIsProvided_LaysOutThatTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	template := &template_model.Template{Name: gofakeit.Word()}
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{Template: template})

	// Assert
	assert.Equal(t, template.Name, laidOutTemplate(t, layoutService).Name)
}

func TestWhenTopologyAndCanvasSideAreProvided_ForwardsThemToTheLayoutService(t *testing.T) {
	t.Parallel()
	// Arrange
	canvasSide := gofakeit.Float64Range(100, 900)
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{
		Topology:   config.TopologyChain,
		CanvasSide: canvasSide,
	})

	// Assert
	layoutService.AssertCalled(t, "BuildPreviewLayout", mock.Anything, config.TopologyChain, canvasSide)
}

func TestWhenLayoutIsComputed_ReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := preview.Layout{ZoneRadius: gofakeit.Float64Range(10, 40)}
	handler := handlers.NewPreviewHandler(newLayoutServiceReturning(expected))

	// Act
	layoutDto := handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{})

	// Assert
	assert.Equal(t, expected, layoutDto.Layout)
}

func newLayoutServiceReturning(layout preview.Layout) *test_helpers.PreviewLayoutServiceMock {
	layoutService := &test_helpers.PreviewLayoutServiceMock{}
	layoutService.On("BuildPreviewLayout", mock.Anything, mock.Anything, mock.Anything).Return(layout)
	return layoutService
}

// laidOutTemplate returns the template the handler handed to the layout service.
func laidOutTemplate(t *testing.T, layoutService *test_helpers.PreviewLayoutServiceMock) *template_model.Template {
	t.Helper()

	template, _ := layoutService.Calls[0].Arguments.Get(0).(*template_model.Template)
	return template
}
