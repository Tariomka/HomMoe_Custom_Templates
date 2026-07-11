package registry

type gameModes struct {
	Classic    string
	SingleHero string
}

// GetGameModeValues returns the available game mode types used for
//
//	gameMode
func GetGameModeValues() gameModes {
	return gameModes{
		Classic:    "Classic",
		SingleHero: "SingleHero",
	}
}

func GetGameModeList() []string {
	gameModeValues := GetGameModeValues()
	return []string{gameModeValues.Classic, gameModeValues.SingleHero}
}
