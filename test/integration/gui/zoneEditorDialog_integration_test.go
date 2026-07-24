//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/handlers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenZoneEditorDialogRenders_UsesHandlerProvidedOptions(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := handlers.NewGuiHandler()
	state := dtos.NewDefaultEditorStateDto()
	generated, err := handler.GenerateTemplate(state)
	require.NoError(t, err)
	require.NotNil(t, generated.Template)
	require.NotEmpty(t, generated.Template.Variants)
	variant := generated.Template.Variants[0]
	options := handler.GetZoneEditorOptions(state, len(variant.Zones))
	dialog := dialogs.NewZoneEditorDialog(
		variant.Zones,
		variant.Connections,
		options.Topology,
		options.Tuning,
		options.GenerateRoads,
		handler,
		handler,
		nil,
	)
	gtx, frameRouter := newDialogContext(image.Pt(1000, 720))

	// Act
	dimensions, closed := dialog.Body(gtx, themes.NewTheme())
	frameRouter.Frame(gtx.Ops)

	// Assert
	assert.Equal(t, image.Pt(1000, 720), dimensions.Size)
	assert.False(t, closed)
}
