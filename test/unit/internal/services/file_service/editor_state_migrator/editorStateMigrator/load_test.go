package editorStateMigrator_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template_entity"
	"github.com/Tariomka/hommoe_custom_templates/internal/helpers/data"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service/editor_state_migrator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenTheFileIsMissing_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	missingPath := filepath.Join(t.TempDir(), "missing.gen.json")
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(missingPath, &target)

	// Assert
	assert.Error(t, err)
}

// The probe fails before either repository is picked, so a malformed file never
// reaches a decode.
func TestWhenTheFileIsNotJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, "{not json")
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenTheFileIsCurrentVersion_TheDecodedKeysReachTheTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"templateName":"X","playerCount":4,"mapSize":192}`)
	expected := editor_state.EditorState{
		TemplateName:  "X",
		PlayerCount:   4,
		MapSize:       192,
		SchemaVersion: editor_state.CurrentEditorStateSchemaVersion,
	}
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, target)
}

// The load writes into whatever it is handed, so a key the file omits keeps
// what the caller seeded. FileService relies on this to seed the defaults.
func TestWhenACurrentFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"playerCount":4}`)
	target := editor_state.EditorState{TemplateName: "Seeded"}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Seeded", target.TemplateName)
}

// The same contract has to hold on the legacy path, which decodes into a
// different struct entirely before migrating.
func TestWhenALegacyFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"playerCount":4}`)
	target := editor_state.EditorState{TemplateName: "Seeded"}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, "Seeded", target.TemplateName)
}

func TestWhenALegacyFileOmitsTheContentRows_TheRowsSeededByTheCallerSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"templateName":"X"}`)
	seededRules := []editor_state.ContentRuleRow{{Name: "Guarded"}}
	seededRows := []editor_state.ZoneContentRow{{Sid: "seeded", Count: 1, Rules: seededRules}}
	target := editor_state.EditorState{PlayerZoneContentRows: seededRows}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, seededRows, target.PlayerZoneContentRows)
}

func TestWhenALegacyFileOmitsTheBonuses_TheBonusesSeededByTheCallerSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"templateName":"X"}`)
	seededBonuses := []editor_state.BonusEntry{{
		PresetType:     editor_state.BonusStartingGold,
		ReceiverFilter: "all_heroes",
		Param:          "9500",
	}}
	target := editor_state.EditorState{Bonuses: seededBonuses}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, seededBonuses, target.Bonuses)
}

// The legacy path decodes into a different struct than the current one, so it
// needs its own proof that a bad file surfaces rather than half-loads.
func TestWhenALegacyFileHasAMistypedValue_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"playerCount":"seven"}`)
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

// v1 wrote the manual position as [x, y]; encoding/json/v2 refuses to read that
// into the current struct, so without the migration the file simply will not
// load.
func TestWhenALegacyFileCarriesAnArrayPosition_ItBecomesAVector(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(
		t,
		`{"schemaVersion":1,"manualZones":[{"zone":{"name":"A"},"manualPosition":[0.25,0.75]}]}`)
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, data.NewVec2(0.25, 0.75), *target.ManualZones[0].ManualPosition)
}

// v0 predates the version key entirely and probes as 0, which is why the gate is
// "below 2" rather than "equal to 1".
func TestWhenTheFileHasNoVersionKey_ItStillTakesTheLegacyPath(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"manualZones":[{"zone":{"name":"A"},"manualPosition":[0.25,0.75]}]}`)
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, data.NewVec2(0.25, 0.75), *target.ManualZones[0].ManualPosition)
}

func TestWhenALegacyFileIsMigrated_TheStateLandsAtTheCurrentVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"templateName":"X"}`)
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, editor_state.CurrentEditorStateSchemaVersion, target.SchemaVersion)
}

// v1 never persisted the generator stamps, so a migrated zone has to read as
// "never stamped" rather than as a zone sitting at the origin.
func TestWhenALegacyFileIsMigrated_TheZoneCarriesNoGeneratorPosition(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"manualZones":[{"zone":{"name":"A"}}]}`)
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, target.ManualZones[0].GeneratorPosition)
}

func TestWhenTheFileIsNewerThanThisBuild_RefusesToLoadIt(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":99,"templateName":"X"}`)
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	assert.ErrorAs(t, err, new(*editor_state_migrator.UnsupportedSchemaVersionError))
}

// A refused load must leave the caller's state exactly as it was: the user
// keeps working on what they had rather than on a half-decoded file.
func TestWhenTheFileIsNewerThanThisBuild_LeavesTheTargetUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":99,"templateName":"X"}`)
	target := editor_state.EditorState{TemplateName: "Seeded"}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.Error(t, err)
	assert.Equal(t, editor_state.EditorState{TemplateName: "Seeded"}, target)
}

func TestWhenACurrentFileCarriesTheGeneratorStamps_TheyReachTheTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"manualZones":[{"zone":{"name":"A"},`+
		`"generatorPosition":{"X":0.4,"Y":0.6},"generatorRing":0,"manualPosition":{"X":0.25,"Y":0.75}}]}`)
	expected := editor_state.ManualZoneSave{
		Zone:              template_entity.Zone{Name: "A"},
		GeneratorPosition: new(data.NewVec2(0.4, 0.6)),
		GeneratorRing:     new(0),
		ManualPosition:    new(data.NewVec2(0.25, 0.75)),
	}
	target := editor_state.EditorState{}

	// Act
	err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, target.ManualZones[0])
}
