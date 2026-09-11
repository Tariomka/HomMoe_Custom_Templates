package utils

import (
	"image/color"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/preview"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
)

type ConnectionLineStyle struct {
	HasRoad        bool
	ExplicitPortal bool
	DrawsAsPortal  bool
}

func NewPreviewConnectionLineStyle(connection preview.Connection) ConnectionLineStyle {
	return ConnectionLineStyle{
		HasRoad:        connection.HasRoad,
		ExplicitPortal: connection.ExplicitPortal,
		DrawsAsPortal:  connection.IsPortal(),
	}
}

func NewEditorConnectionLineStyle(connection template_model.Connection) ConnectionLineStyle {
	return ConnectionLineStyle{
		HasRoad:        connection.HasRoad(),
		ExplicitPortal: connection.IsExplicitPortal(),
		DrawsAsPortal:  connection.IsExplicitPortal(),
	}
}

func (this ConnectionLineStyle) Color() color.NRGBA {
	if !this.HasRoad {
		if this.ExplicitPortal {
			return themes.ColorsPreview.PortalNoRoadLine
		}

		return themes.ColorsPreview.NoRoadLine
	}

	if this.DrawsAsPortal {
		return themes.ColorsPreview.PortalLine
	}

	return themes.ColorsPreview.DirectLine
}
