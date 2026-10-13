package editorStateMigrator_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/editor_state"
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
	_, err := newMigrator().Load(missingPath, &target)

	// Assert
	assert.Error(t, err)
}

// The probe fails before any repository is picked, so a malformed file never
// reaches a decode.
func TestWhenTheFileIsNotJson_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, "{not json")
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenTheFileIsCurrentVersion_TheDecodedKeysReachTheTarget(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t,
		`{"schemaVersion":3,"generatorVersion":1,"templateName":"X","mapSize":192,"players":{"playerCount":4}}`)
	expected := editor_state.EditorState{
		SchemaVersion: 3, GeneratorVersion: 1,
		TemplateName: "X",
		MapSize:      192,
		Players:      editor_state.PlayerSettings{PlayerCount: 4},
	}
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, target)
}

func TestWhenTheFileIsCurrentVersion_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":3,"templateName":"X"}`)
	target := editor_state.EditorState{}

	// Act
	warnings, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, warnings)
}

// The load writes into whatever it is handed, so a key the file omits keeps
// what the caller seeded. FileService relies on this to seed the defaults.
func TestWhenACurrentFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":3,"players":{"playerCount":4}}`)
	target := editor_state.EditorState{Players: editor_state.PlayerSettings{HeroCountMax: 8}}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 8, target.Players.HeroCountMax)
}

func TestWhenTheContentDefaultsVersionIsNewerThanThisBuild_RefusesToLoadIt(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":3,"content":{"contentDefaultsVersion":2}}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	assert.ErrorAs(t, err, new(*editor_state_migrator.UnsupportedContentDefaultsVersionError))
}

func TestWhenTheContentDefaultsVersionIsNewerThanThisBuild_LeavesTheTargetUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":3,"templateName":"X","content":{"contentDefaultsVersion":2}}`)
	seeded := editor_state.EditorState{TemplateName: "Seeded"}
	target := seeded

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.Error(t, err)
	assert.Equal(t, seeded, target)
}

func TestWhenAVersionTwoFileIsMigrated_TheFlatSettingsMoveIntoTheirSections(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"templateName":"X","playerCount":4,"topology":"Square"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t,
		[]any{"X", 4, "Square"},
		[]any{target.TemplateName, target.Players.PlayerCount, string(target.Generation.Topology)})
}

func TestWhenAVersionTwoFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"playerCount":4}`)
	target := editor_state.EditorState{Players: editor_state.PlayerSettings{HeroCountMax: 8}}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 8, target.Players.HeroCountMax)
}

func TestWhenAVersionTwoFileCarriesBanText_ItBecomesBanEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"bannedItems":"  item_a  \n\n\t\nitem_b\nitem_a"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t,
		[]editor_state.BanEntry{{Sid: "item_a"}, {Sid: "item_b"}, {Sid: "item_a"}},
		target.Content.BannedItems)
}

func TestWhenAVersionTwoFileCarriesOverrideText_TheValidLinesBecomeEntries(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t,
		`{"schemaVersion":2,"valueOverrides":"watchtower=25000\n\n  =5 \nbad_line\ngold_mine = 12000 \nnonnum=abc"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(
		t,
		[]editor_state.ValueOverrideEntry{
			{Sid: "watchtower", GuardValue: 25000},
			{Sid: "gold_mine", GuardValue: 12000},
		},
		target.Content.ValueOverrides,
	)
}

func TestWhenAVersionTwoOverrideLineIsMalformed_ReturnsAWarningNamingTheLine(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name            string
		overridesText   string
		expectedWarning string
	}{
		{"LineHasNoSeparator_ReturnsWarning", `bad_line`,
			"value override line 1 'bad_line' is not sid=value and was dropped"},
		{"LineStartsWithSeparator_ReturnsWarning", `=5`,
			"value override line 1 '=5' is not sid=value and was dropped"},
		{"SeparatorFollowsOnlyWhitespace_ReturnsWarningForTheTrimmedLine", `  = 5`,
			"value override line 1 '= 5' is not sid=value and was dropped"},
		{"ValueIsNotANumber_ReturnsWarning", `nonnum=abc`,
			"value override line 1 'nonnum=abc' has a non-numeric value and was dropped"},
		{"RejectedLineFollowsBlankLines_ReturnsWarningWithTheSourceLineNumber", `watchtower=25000\n\nbad_line`,
			"value override line 3 'bad_line' is not sid=value and was dropped"},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			statePath := writeStateFile(t, `{"schemaVersion":2,"valueOverrides":"`+testCase.overridesText+`"}`)
			target := editor_state.EditorState{}

			// Act
			warnings, err := newMigrator().Load(statePath, &target)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, []string{testCase.expectedWarning}, warnings)
		})
	}
}

func TestWhenAVersionTwoOverrideTextIsOnlyBlankLines_ReturnsNoWarnings(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"valueOverrides":"\n   \n\t\n"}`)
	target := editor_state.EditorState{}

	// Act
	warnings, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Empty(t, warnings)
}

func TestWhenAVersionTwoFileCarriesGeneratorStamps_TheyReachTheZoneStructure(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"manualZones":[{"zone":{"name":"A","size":1.5},`+
		`"generatorPosition":{"X":0.4,"Y":0.6},"generatorRing":0,"manualPosition":{"X":0.25,"Y":0.75}}]}`)
	expected := editor_state.ManualZoneEdit{
		Name:              "A",
		GeneratorPosition: new(data.NewVec2(0.4, 0.6)),
		GeneratorRing:     new(0),
		ManualPosition:    new(data.NewVec2(0.25, 0.75)),
		FieldOverrides:    editor_state.ZoneFieldOverrides{Size: new(1.5)},
	}
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, target.ManualEdits.Zones[0])
}

// A migrated snapshot keeps every field; generator version 0 marks it for
// conversion at the first regeneration.
func TestWhenAVersionTwoFileIsMigrated_ItLandsAtTheCurrentVersionWithNoGeneratorVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"templateName":"X"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t,
		editor_state.SchemaOptions{SchemaVersion: editor_state.CurrentEditorStateSchemaVersion, GeneratorVersion: 0},
		target.SchemaOptions)
}

func TestWhenAVersionTwoFileIsMigrated_ItUsesTheFirstContentDefaultsVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"templateName":"X"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 1, target.Content.ContentDefaultsVersion)
}

// The same seeding contract has to hold on the v1 path, which decodes into a
// different struct entirely before migrating twice.
func TestWhenALegacyFileOmitsAKey_TheValueSeededByTheCallerSurvives(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"playerCount":4}`)
	target := editor_state.EditorState{TemplateName: "Seeded"}

	// Act
	_, err := newMigrator().Load(statePath, &target)

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
	target := editor_state.EditorState{Content: editor_state.ContentSettings{PlayerZoneContentRows: seededRows}}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, seededRows, target.Content.PlayerZoneContentRows)
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
	target := editor_state.EditorState{Content: editor_state.ContentSettings{Bonuses: seededBonuses}}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, seededBonuses, target.Content.Bonuses)
}

func TestWhenALegacyFileOmitsTheBans_TheBansSeededByTheCallerSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"templateName":"X"}`)
	seededBans := []editor_state.BanEntry{{Sid: "seeded_a"}, {Sid: "seeded_b"}}
	target := editor_state.EditorState{Content: editor_state.ContentSettings{BannedMagics: seededBans}}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, seededBans, target.Content.BannedMagics)
}

func TestWhenALegacyFileOmitsTheOverrides_TheOverridesSeededByTheCallerSurvive(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"templateName":"X"}`)
	seededOverrides := []editor_state.ValueOverrideEntry{{Sid: "seeded", GuardValue: 700}}
	target := editor_state.EditorState{Content: editor_state.ContentSettings{ValueOverrides: seededOverrides}}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, seededOverrides, target.Content.ValueOverrides)
}

// The legacy path decodes into a different struct than the current one, so it
// needs its own proof that a bad file surfaces rather than half-loads.
func TestWhenALegacyFileHasAMistypedValue_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"playerCount":"seven"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenACurrentFileHasAMistypedValue_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":3,"players":{"playerCount":"seven"}}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

func TestWhenAVersionTwoFileHasAMistypedValue_ReturnsError(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":2,"playerCount":"seven"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	assert.Error(t, err)
}

// v1 wrote the manual position as [x, y]; encoding/json/v2 refuses to read that
// into a newer struct, so without the migration the file simply will not load.
func TestWhenALegacyFileCarriesAnArrayPosition_ItBecomesAVector(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(
		t,
		`{"schemaVersion":1,"manualZones":[{"zone":{"name":"A"},"manualPosition":[0.25,0.75]}]}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, data.NewVec2(0.25, 0.75), *target.ManualEdits.Zones[0].ManualPosition)
}

// v0 predates the version key entirely and probes as 0, which is why the gate is
// "below 2" rather than "equal to 1".
func TestWhenTheFileHasNoVersionKey_ItStillTakesTheLegacyPath(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"manualZones":[{"zone":{"name":"A"},"manualPosition":[0.25,0.75]}]}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, data.NewVec2(0.25, 0.75), *target.ManualEdits.Zones[0].ManualPosition)
}

func TestWhenALegacyFileIsMigrated_TheStateLandsAtTheCurrentVersion(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":1,"templateName":"X"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

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
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.NoError(t, err)
	assert.Nil(t, target.ManualEdits.Zones[0].GeneratorPosition)
}

func TestWhenTheFileIsNewerThanThisBuild_RefusesToLoadIt(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":99,"templateName":"X"}`)
	target := editor_state.EditorState{}

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	assert.ErrorAs(t, err, new(*editor_state_migrator.UnsupportedSchemaVersionError))
}

// A refused load must leave the caller's state exactly as it was: the user
// keeps working on what they had rather than on a half-decoded file.
func TestWhenTheFileIsNewerThanThisBuild_LeavesTheTargetUntouched(t *testing.T) {
	t.Parallel()
	// Arrange
	statePath := writeStateFile(t, `{"schemaVersion":99,"templateName":"X"}`)
	seeded := editor_state.EditorState{TemplateName: "Seeded"}
	target := seeded

	// Act
	_, err := newMigrator().Load(statePath, &target)

	// Assert
	require.Error(t, err)
	assert.Equal(t, seeded, target)
}
