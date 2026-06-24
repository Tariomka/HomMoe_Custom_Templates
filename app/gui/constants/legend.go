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
		{Label: "Player", Color: themes.ColorPreviewSpawnEdge},
		{Label: "Bronze", Color: themes.ColorPreviewBronzeEdge},
		{Label: "Silver", Color: themes.ColorPreviewSilverEdge},
		{Label: "Gold", Color: themes.ColorPreviewGoldEdge},
		{Label: "Hub", Color: themes.ColorPreviewHubEdge},
	},
	{
		{Label: "Road", Color: themes.ColorPreviewDirectLine, Line: true},
		{Label: "Portal", Color: themes.ColorPreviewPortalLine, Line: true},
	},
}
