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

//nolint:gochecknoglobals // Allocated once and reused on each render
var LegendRows = [][]LegendItem{
	{
		{Label: "Player", Color: themes.ColorsPreview.SpawnEdge},
		{Label: "Plastic (T1)", Color: themes.ColorsPreview.PlasticEdge},
		{Label: "Bronze (T2)", Color: themes.ColorsPreview.BronzeEdge},
		{Label: "Silver (T3)", Color: themes.ColorsPreview.SilverEdge},
		{Label: "Gold (T4)", Color: themes.ColorsPreview.GoldEdge},
		{Label: "Hub (T5)", Color: themes.ColorsPreview.HubEdge},
	},
	ConnectionLegendRow,
}

//nolint:gochecknoglobals // Allocated once and reused on each render
var ConnectionLegendRow = []LegendItem{
	{Label: "Road", Color: themes.ColorsPreview.DirectLine, Line: true},
	{Label: "No road", Color: themes.ColorsPreview.NoRoadLine, Line: true},
	{Label: "Portal", Color: themes.ColorsPreview.PortalLine, Line: true},
	{Label: "Portal without road", Color: themes.ColorsPreview.PortalNoRoadLine, Line: true},
}
