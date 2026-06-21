package constants

import "github.com/Tariomka/hommoe_custom_templates/internal/registry"

var GameModeValues = registry.GetGameModeValues()
var GameModes = []string{
	GameModeValues.Classic,
	GameModeValues.SingleHero,
}
