package test_helpers

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

	return template
}
