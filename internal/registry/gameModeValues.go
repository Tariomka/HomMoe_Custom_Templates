package registry

type gameModes struct {
	Classic    string
	SingleHero string
}

var gameModeValues = gameModes{
	Classic:    "Classic",
	SingleHero: "SingleHero",
}

// GetGameModeValues returns the available game mode types used for
//
//	gameMode
func GetGameModeValues() gameModes {
	return gameModeValues
}
