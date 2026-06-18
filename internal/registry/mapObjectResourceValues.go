package registry

type objectResources struct {
	Gold       string
	Dust       string
	Wood       string
	Ore        string
	Crystal    string
	Mercury    string
	Gemstone   string
	Chest      string
	PandoraBox string
	CampFire   string
}

var objectResourceValues = objectResources{
	Gold:       "resource_gold",
	Dust:       "resource_dust",
	Wood:       "resource_wood",
	Ore:        "resource_ore",
	Crystal:    "resource_crystals",
	Mercury:    "resource_mercury",
	Gemstone:   "resource_gemstones",
	Chest:      "chest",
	PandoraBox: "pandora_box",
	CampFire:   "camp_fire",
}

func GetMapObjectResourceValues() objectResources {
	return objectResourceValues
}
