package constants

import (
	"image/color"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/themes"
)

type LegendItem struct {
	Label string
	Color color.NRGBA
	Line  bool
}

var LegendRows = [][]LegendItem{
	{
		{Label: "Player", Color: themes.ColorsPreview.SpawnEdge},
		{Label: "Bronze", Color: themes.ColorsPreview.BronzeEdge},
		{Label: "Silver", Color: themes.ColorsPreview.SilverEdge},
		{Label: "Gold", Color: themes.ColorsPreview.GoldEdge},
		{Label: "Hub", Color: themes.ColorsPreview.HubEdge},
	},
	{
		{Label: "Road", Color: themes.ColorsPreview.DirectLine, Line: true},
		{Label: "Portal", Color: themes.ColorsPreview.PortalLine, Line: true},
	},
}
