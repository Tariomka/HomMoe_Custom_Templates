package constants

import (
	"image/color"

	"github.com/Tariomka/hommoe_custom_templates/internal/gui/themes"
)

type legendItem struct {
	Label string
	Color color.NRGBA
}

var LegendItems = []legendItem{
	{Label: "Player", Color: themes.ColorPreviewSpawnEdge},
	{Label: "Bronze", Color: themes.ColorPreviewBronzeEdge},
	{Label: "Silver", Color: themes.ColorPreviewSilverEdge},
	{Label: "Gold", Color: themes.ColorPreviewGoldEdge},
	{Label: "Hub", Color: themes.ColorPreviewHubEdge},
}
