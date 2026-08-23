//go:build integration_test && gui

package gui_test

import (
	"image"
	"testing"
	"time"

	"gioui.org/io/input"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/dialogs"
	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenManageRulesDialogHasVariantRule_RendersContent(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 2
	dialog := dialogs.NewManageRulesDialog(
		constants.ContentIDs.DragonUtopia,
		[]models.ContentRuleRow{{Name: "Variant", VariantID: &variantID}},
		composition.InitializeGuiHandler(),
		nil)
	gtx, frameRouter := newDialogContext(image.Pt(540, 500))

	// Act
	dimensions, _ := dialog.Body(gtx, themes.NewTheme())
	frameRouter.Frame(gtx.Ops)

	// Assert
	assert.NotEqual(t, image.Point{}, dimensions.Size)
}

func TestWhenZoneContentDialogRenders_PreservesSavedRules(t *testing.T) {
	t.Parallel()
	// Arrange
	variantID := 2
	guarded := false
	expected := []models.ZoneContentRow{{
		Sid:   constants.ContentIDs.DragonUtopia.Sid,
		Count: 2,
		Rules: []models.ContentRuleRow{
			{Name: "Variant", VariantID: &variantID},
			{Name: "Guarded", IsGuarded: &guarded},
		},
	}}
	var persisted []models.ZoneContentRow
	dialog := dialogs.NewZoneContentDialog(
		"Zone Content: High Neutral",
		false,
		expected,
		composition.InitializeGuiHandler(),
		nil,
		func(rows []models.ZoneContentRow) { persisted = rows })
	gtx, frameRouter := newDialogContext(image.Pt(640, 560))

	// Act
	dialog.Body(gtx, themes.NewTheme())
	frameRouter.Frame(gtx.Ops)

	// Assert
	assert.Equal(t, expected, persisted)
}

func newDialogContext(size image.Point) (layout.Context, *input.Router) {
	var operations op.Ops
	frameRouter := new(input.Router)
	return layout.Context{
		Ops:         &operations,
		Constraints: layout.Exact(size),
		Metric:      unit.Metric{PxPerDp: 1, PxPerSp: 1},
		Source:      frameRouter.Source(),
		Now:         time.Now(),
	}, frameRouter
}
