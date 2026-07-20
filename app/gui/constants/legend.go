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
		{Label: "Plastic (T1)", Color: themes.ColorsPreview.PlasticEdge},
		{Label: "Bronze (T2)", Color: themes.ColorsPreview.BronzeEdge},
		{Label: "Silver (T3)", Color: themes.ColorsPreview.SilverEdge},
		{Label: "Gold (T4)", Color: themes.ColorsPreview.GoldEdge},
		{Label: "Hub (T5)", Color: themes.ColorsPreview.HubEdge},
	},
	{
		{Label: "Road", Color: themes.ColorsPreview.DirectLine, Line: true},
		{Label: "Portal", Color: themes.ColorsPreview.PortalLine, Line: true},
	},
}
