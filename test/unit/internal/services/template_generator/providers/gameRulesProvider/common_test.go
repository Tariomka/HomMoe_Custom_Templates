package gameRulesProvider_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/require"
)

// createGameRules builds a default configuration, applies the optional mutation
// and runs it through the provider's CreateGameRules.
func createGameRules(configure func(configuration *config.GeneratorConfig)) entities.GameRules {
	configuration := config.NewGeneratorConfig()
	if configure != nil {
		configure(configuration)
	}
	return providers.NewGameRulesProvider().CreateGameRules(*configuration)
}

// bonusesFor runs the provider's bonus expansion (via CreateGameRules) for the
// given UI bonus entries and returns the produced raw bonuses.
func bonusesFor(entries ...config.BonusEntry) entities.BonusList {
	return createGameRules(func(configuration *config.GeneratorConfig) {
		configuration.Bonuses = entries
	}).Bonuses
}

// loadExampleTemplate reads one of the real game templates shipped under
// data/ExampleTemplates for functional-equivalence checks.
func loadExampleTemplate(t *testing.T, name string) entities.RmgTemplate {
	t.Helper()
	path, err := filepath.Abs(filepath.Join(
		"..", "..", "..", "..", "..", "..", "..", "data", "ExampleTemplates", name))
	require.NoError(t, err)
	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	var parsedTemplate entities.RmgTemplate
	require.NoError(t, json.Unmarshal(raw, &parsedTemplate))
	return parsedTemplate
}
