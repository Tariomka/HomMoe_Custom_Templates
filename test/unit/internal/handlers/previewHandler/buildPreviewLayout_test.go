package previewHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/test/test_helpers"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestWhenRequestIsEmpty_LaysOutNoTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	_, _ = handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{})

	// Assert
	layoutService.AssertCalled(t, "BuildPreviewLayout", (*entities.RmgTemplate)(nil), mock.Anything, mock.Anything)
}

func TestWhenOnlyZonesAreProvided_LaysOutASynthesizedTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	zones := []entities.Zone{{Name: gofakeit.Word()}}
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	_, _ = handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{Zones: zones})

	// Assert
	assert.Equal(t, zones, laidOutTemplate(t, layoutService).Variants[0].Zones)
}

func TestWhenOnlyConnectionsAreProvided_LaysOutASynthesizedTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	connections := []entities.Connection{{Name: gofakeit.Word()}}
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	_, _ = handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{Connections: connections})

	// Assert
	assert.Equal(t, connections, laidOutTemplate(t, layoutService).Variants[0].Connections)
}

func TestWhenTemplateIsProvided_LaysOutThatTemplate(t *testing.T) {
	t.Parallel()
	// Arrange
	template := &entities.RmgTemplate{Name: gofakeit.Word()}
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	_, _ = handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{
		Template: template,
		Zones:    []entities.Zone{{Name: gofakeit.Word()}},
	})

	// Assert
	assert.Same(t, template, laidOutTemplate(t, layoutService))
}

func TestWhenTopologyAndCanvasSideAreProvided_ForwardsThemToTheLayoutService(t *testing.T) {
	t.Parallel()
	// Arrange
	canvasSide := gofakeit.Float64Range(100, 900)
	layoutService := newLayoutServiceReturning(preview.Layout{})
	handler := handlers.NewPreviewHandler(layoutService)

	// Act
	_, _ = handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{
		Topology:   config.TopologyChain,
		CanvasSide: canvasSide,
	})

	// Assert
	layoutService.AssertCalled(t, "BuildPreviewLayout", mock.Anything, config.TopologyChain, canvasSide)
}

func TestWhenLayoutIsComputed_ReturnsIt(t *testing.T) {
	t.Parallel()
	// Arrange
	expected := preview.Layout{ZoneRadius: gofakeit.IntRange(10, 40)}
	handler := handlers.NewPreviewHandler(newLayoutServiceReturning(expected))

	// Act
	layoutDto, err := handler.BuildPreviewLayout(dtos.PreviewLayoutRequestDto{})

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, layoutDto.Layout)
}

func newLayoutServiceReturning(layout preview.Layout) *test_helpers.PreviewLayoutServiceMock {
	layoutService := &test_helpers.PreviewLayoutServiceMock{}
	layoutService.On("BuildPreviewLayout", mock.Anything, mock.Anything, mock.Anything).Return(layout)
	return layoutService
}

// laidOutTemplate returns the template the handler handed to the layout service.
func laidOutTemplate(t *testing.T, layoutService *test_helpers.PreviewLayoutServiceMock) *entities.RmgTemplate {
	t.Helper()

	template, _ := layoutService.Calls[0].Arguments.Get(0).(*entities.RmgTemplate)
	return template
}
