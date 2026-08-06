package constants

import "github.com/Tariomka/hommoe_custom_templates/internal/registry"

func GetGameModes() []string {
	gameModes := registry.GetGameModeValues()
	return []string{
		gameModes.Classic,
		gameModes.SingleHero,
	}
}
