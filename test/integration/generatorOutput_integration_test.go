package integration_test

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_generator"
	"github.com/Tariomka/hommoe_custom_templates/internal/composition"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// updateGeneratorHashes rewrites testdata/generatorOutputHashes.json for the
// current generator version instead of checking against it.
//
//nolint:gochecknoglobals // command-line flags are the idiomatic exception.
var updateGeneratorHashes = flag.Bool("update-generator-hashes", false,
	"record the generator output hashes for the current GeneratorOutputVersion")

const generatorHashesPath = "testdata/generatorOutputHashes.json"

// generatorOutputHashes pins one generation: the exported .rmg.json and the
// editor model of variant 0, which also carries the generator positions, rings
// and tiers the export drops.
type generatorOutputHashes struct {
	Rmg   string `json:"rmg"`
	Model string `json:"model"`
}

// The derived values of stored manual edits are only reproducible while the
// generator output for the same settings is. Any change to it must come with
// a GeneratorOutputVersion bump, which the editor reports on load. Random
// layouts pin only their per-zone content.
func TestWhenALayoutIsGenerated_TheOutputMatchesTheRecordedVersion(t *testing.T) {
	t.Parallel()
	version := strconv.Itoa(common_generator.GeneratorOutputVersion)
	current := make(map[string]generatorOutputHashes)
	for name, state := range deterministicGeneratorCases() {
		current[name] = hashGeneration(t, state)
	}
	for name, state := range nonDeterministicGeneratorCases() {
		current[name] = generatorOutputHashes{Model: hashJSON(t, zoneContentByName(generateVariant(t, state)))}
	}
	recorded := readGeneratorHashes(t)
	if *updateGeneratorHashes {
		recorded[version] = current
		writeGeneratorHashes(t, recorded)
		return
	}
	require.Containsf(t, recorded, version,
		"no hashes are recorded for generator version %s; record them with -update-generator-hashes", version)

	for name, hashes := range current {
		t.Run(name+"_MatchesTheRecordedHashes", func(t *testing.T) {
			t.Parallel()
			// Arrange
			expected := recorded[version][name]

			// Act
			actual := hashes

			// Assert
			assert.Equalf(t, expected, actual,
				"the generator output changed at version %s; bump common_generator.GeneratorOutputVersion "+
					"and record the new hashes with -update-generator-hashes", version)
		})
	}
}

func TestWhenADeterministicLayoutIsGeneratedTwice_TheOutputIsIdentical(t *testing.T) {
	t.Parallel()
	for name, state := range deterministicGeneratorCases() {
		t.Run(name+"_IsIdentical", func(t *testing.T) {
			t.Parallel()
			// Arrange
			first := hashGeneration(t, state)

			// Act
			second := hashGeneration(t, state)

			// Assert
			assert.Equal(t, first, second)
		})
	}
}

// Random layouts differ per generation in their placement, roads and
// connections, but every named zone must come back with the same content.
func TestWhenANonDeterministicLayoutIsGeneratedTwice_TheZoneContentIsIdentical(t *testing.T) {
	t.Parallel()
	for name, state := range nonDeterministicGeneratorCases() {
		t.Run(name+"_ZoneContentIsIdentical", func(t *testing.T) {
			t.Parallel()
			// Arrange
			first := zoneContentByName(generateVariant(t, state))

			// Act
			second := zoneContentByName(generateVariant(t, state))

			// Assert
			assert.Equal(t, first, second)
		})
	}
}

// Manual-edit capture tells generated zones from editor-made ones by the
// generator position every topology stamps.
func TestWhenAnyLayoutIsGenerated_EveryZoneCarriesAGeneratorPosition(t *testing.T) {
	t.Parallel()
	cases := deterministicGeneratorCases()
	maps.Copy(cases, nonDeterministicGeneratorCases())
	for name, state := range cases {
		t.Run(name+"_StampsEveryZone", func(t *testing.T) {
			t.Parallel()
			// Arrange
			template := generateVariant(t, state)

			// Act
			unstamped := unstampedZoneNames(template)

			// Assert
			assert.Empty(t, unstamped)
		})
	}
}

func deterministicGeneratorCases() map[string]editor_state_model.EditorState {
	cases := map[string]editor_state_model.EditorState{}
	for _, topology := range []config.MapTopology{
		config.TopologyCircles, config.TopologySquare, config.TopologyGeometric,
		config.TopologyCross, config.TopologyFractal, config.TopologyGeometricHub,
	} {
		cases[string(topology)] = generatorCase(topology, false)
	}
	tournament := generatorCase(config.TopologySquare, false)
	tournament.PlayerCount = 2
	tournament.VictoryCondition = registry.GetWinningConditionValues().Tournament
	cases["Tournament2p"] = tournament
	return cases
}

func nonDeterministicGeneratorCases() map[string]editor_state_model.EditorState {
	return map[string]editor_state_model.EditorState{
		"Random":                 generatorCase(config.TopologyRandom, false),
		"SquareRandomPortals":    generatorCase(config.TopologySquare, true),
		"GeometricRandomPortals": generatorCase(config.TopologyGeometric, true),
	}
}

func generatorCase(topology config.MapTopology, randomPortals bool) editor_state_model.EditorState {
	state := editor_state_model.NewDefaultEditorStateModel()
	state.TemplateName = "Generator Guard"
	state.Topology = topology
	state.RandomPortals = randomPortals
	state.PlayerCount = 4
	state.NeutralZoneCount = 6
	state.MapSize = 256
	return state
}

func generateVariant(t *testing.T, state editor_state_model.EditorState) *template_model.Template {
	t.Helper()
	generated, err := composition.InitializeGuiHandler().
		GenerateTemplate(editor_state_dto.EditorStateDto{EditorState: state})
	require.NoError(t, err)
	require.NotEmpty(t, generated.Template.Variants)

	return generated.Template
}

func hashGeneration(t *testing.T, state editor_state_model.EditorState) generatorOutputHashes {
	t.Helper()
	template := generateVariant(t, state)

	return generatorOutputHashes{
		Rmg:   hashJSON(t, mappers.NewTemplateMapper().ToEntity(*template)),
		Model: hashJSON(t, template.Variants[0]),
	}
}

func hashJSON(t *testing.T, value any) string {
	t.Helper()
	encoded, err := json.Marshal(value)
	require.NoError(t, err)
	sum := sha256.Sum256(encoded)

	return hex.EncodeToString(sum[:])
}

func unstampedZoneNames(template *template_model.Template) []string {
	var names []string
	for _, zone := range template.Variants[0].Zones {
		if zone.GeneratorPosition == nil {
			names = append(names, zone.Name)
		}
	}

	return names
}

func zoneContentByName(template *template_model.Template) map[string]template_model.Zone {
	zones := make(map[string]template_model.Zone, len(template.Variants[0].Zones))
	for _, zone := range template.Variants[0].Zones {
		zone.GeneratorPosition = nil
		zone.GeneratorRing = nil
		zone.Roads = nil
		zones[zone.Name] = zone
	}

	return zones
}

func readGeneratorHashes(t *testing.T) map[string]map[string]generatorOutputHashes {
	t.Helper()
	recorded := map[string]map[string]generatorOutputHashes{}
	content, err := os.ReadFile(filepath.FromSlash(generatorHashesPath))
	if os.IsNotExist(err) {
		return recorded
	}
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(content, &recorded))

	return recorded
}

func writeGeneratorHashes(t *testing.T, recorded map[string]map[string]generatorOutputHashes) {
	t.Helper()
	content, err := json.MarshalIndent(recorded, "", "  ")
	require.NoError(t, err)
	path := filepath.FromSlash(generatorHashesPath)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
	require.NoError(t, os.WriteFile(path, append(content, '\n'), 0o600))
}
