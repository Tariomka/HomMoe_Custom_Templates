package test_helpers

import (
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/stretchr/testify/mock"
)

// PreviewLayoutServiceMock is a testify mock of
// preview_service.IPreviewLayoutService, used to unit-test collaborators
// without computing a real layout.
type PreviewLayoutServiceMock struct {
	mock.Mock
}

func (this *PreviewLayoutServiceMock) BuildPreviewLayout(
	template *template_model.Template,
	topology config.MapTopology,
	side float64) preview.Layout {
	arguments := this.Called(template, topology, side)
	layout, _ := arguments.Get(0).(preview.Layout)
	return layout
}
