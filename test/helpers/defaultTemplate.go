package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
)

func GetDefaultTemplate() entities.RmgTemplate {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		panic("helpers: unable to resolve the defaultTemplate.json fixture path")
	}

	data, err := os.ReadFile(filepath.Join(filepath.Dir(currentFile), "defaultTemplate.json"))
	if err != nil {
		panic(err)
	}

	var template entities.RmgTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		panic(err)
	}

	template.GameRules.Bonuses = entities.BonusList{}

	zones := template.Variants[0].Zones
	zones[0].GeneratorPosition = &[2]float64{0.8799878400648531, 0.4969600324265629} // Most likely this is random
	zones[0].GeneratorRing = new(0)
	zones[1].GeneratorPosition = &[2]float64{0.12, 0.5}
	zones[1].GeneratorRing = new(0)

	return template
}
