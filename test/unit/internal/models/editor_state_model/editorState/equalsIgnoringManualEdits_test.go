package editorState_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/editor_state_model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStatesAreFullyIdentical_ReportsEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := editor_state_model.NewDefaultEditorStateModel()
	right := left

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.True(t, equal)
}

func TestWhenOnlyManualEditFieldsDiffer_ReportsEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := editor_state_model.NewDefaultEditorStateModel()
	right := left
	right.ManualZones = []editor_state_model.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}
	right.ManualConnections = []editor_state_model.ManualConnectionSave{
		{Connection: entities.Connection{Name: "A-B"}, IsUserAdded: true},
	}

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.True(t, equal)
}

func TestWhenNonManualFieldDiffers_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := editor_state_model.NewDefaultEditorStateModel()
	right := left
	right.TemplateName = "Different Name"

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.False(t, equal)
}

func TestWhenStatesAreDeepClonesWithDistinctRulePointers_ReportsEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := fuzzedEditorState()
	right := deepCloneEditorState(left)

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.True(t, equal)
}

func TestWhenRuleGuardedValueDiffers_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := fuzzedEditorState()
	right := deepCloneEditorState(left)
	*right.HighNeutralContentRows[0].Rules[0].IsGuarded = !*left.HighNeutralContentRows[0].Rules[0].IsGuarded

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.False(t, equal)
}

func TestWhenRuleGuardedIsNilOnOneSide_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := fuzzedEditorState()
	right := deepCloneEditorState(left)
	right.HighNeutralContentRows[0].Rules[0].IsGuarded = nil

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.False(t, equal)
}

func TestWhenContentRowSliceIsNilVersusEmpty_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := fuzzedEditorState()
	right := deepCloneEditorState(left)
	right.MediumNeutralContentRows = []editor_state_model.ZoneContentRow{}

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.False(t, equal)
}

func TestWhenBonusEntryDiffers_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := fuzzedEditorState()
	right := deepCloneEditorState(left)
	right.Bonuses[0].Param = left.Bonuses[0].Param + "0"

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.False(t, equal)
}

// TestWhenFuzzedStatePairsCompared_MatchesReflectDeepEqual pins the
// hand-rolled comparison to the [reflect.DeepEqual] reference it replaced, so
// the two implementations cannot drift apart.
func TestWhenFuzzedStatePairsCompared_MatchesReflectDeepEqual(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		subtestName string
		mutate      func(state *editor_state_model.EditorState)
	}{
		{"StatesAreDeepClones_MatchesDeepEqual", func(_ *editor_state_model.EditorState) {}},
		{"IntFieldDiffers_MatchesDeepEqual", func(state *editor_state_model.EditorState) { state.PlayerCount++ }},
		{
			"FloatFieldDiffers_MatchesDeepEqual",
			func(state *editor_state_model.EditorState) { state.PlayerZoneSize += 0.25 },
		},
		{"BoolFieldDiffers_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.AdvancedMode = !state.AdvancedMode
		}},
		{"TopologyDiffers_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.Topology = config.TopologyHubAndSpoke
		}},
		{
			"BonusEntryDiffers_MatchesDeepEqual",
			func(state *editor_state_model.EditorState) { state.Bonuses[0].Param += "0" },
		},
		{
			"BonusesNilOnOneSide_MatchesDeepEqual",
			func(state *editor_state_model.EditorState) { state.Bonuses = nil },
		},
		{"RowAppended_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.HighNeutralContentRows = append(
				state.HighNeutralContentRows,
				editor_state_model.ZoneContentRow{Sid: "extra"},
			)
		}},
		{
			"RowCountFieldDiffers_MatchesDeepEqual",
			func(state *editor_state_model.EditorState) { state.HighNeutralContentRows[0].Count++ },
		},
		{
			"RowSidDiffers_MatchesDeepEqual",
			func(state *editor_state_model.EditorState) { state.HighNeutralContentRows[0].Sid += "x" },
		},
		{"RuleNameDiffers_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.HighNeutralContentRows[0].Rules[0].Name += "x"
		}},
		{"RuleGuardedValueDiffers_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			*state.HighNeutralContentRows[0].Rules[0].IsGuarded = !*state.HighNeutralContentRows[0].Rules[0].IsGuarded
		}},
		{"RuleGuardedNilOnOneSide_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.HighNeutralContentRows[0].Rules[0].IsGuarded = nil
		}},
		{"RuleSoloEncounterValueDiffers_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			*state.HighNeutralContentRows[0].Rules[0].IsSoloEncounter = !*state.HighNeutralContentRows[0].Rules[0].IsSoloEncounter
		}},
		{"RuleVariantIDValueDiffers_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			*state.HighNeutralContentRows[0].Rules[0].VariantID++
		}},
		{"RulesNilVersusEmpty_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.HighNeutralContentRows[1].Rules = []editor_state_model.ContentRuleRow{}
		}},
		{"RowsNilVersusEmpty_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.MediumNeutralContentRows = []editor_state_model.ZoneContentRow{}
		}},
		{"ManualZonesDiffer_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.ManualZones = []editor_state_model.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}
		}},
		{"ManualConnectionsDiffer_MatchesDeepEqual", func(state *editor_state_model.EditorState) {
			state.ManualConnections = []editor_state_model.ManualConnectionSave{
				{Connection: entities.Connection{Name: "A-B"}, IsUserAdded: true},
			}
		}},
	}
	for _, testCase := range testCases {
		t.Run(testCase.subtestName, func(t *testing.T) {
			t.Parallel()
			// Arrange
			left := fuzzedEditorState()
			right := deepCloneEditorState(left)
			testCase.mutate(&right)

			// Act
			equal := left.EqualsIgnoringManualEdits(&right)

			// Assert
			assert.Equal(t, reflectDeepEqualIgnoringManualEdits(left, right), equal)
		})
	}
}

// TestWhenAnyNonManualFieldIsMutated_ReportsNotEqual walks every leaf field of
// the model via reflection - descending through the embedded entity groups -
// mutates it on a deep clone and expects inequality (except for the ignored
// manual-edit leaves). It trips when a new field is added to any group without
// extending the hand-rolled comparison.
func TestWhenAnyNonManualFieldIsMutated_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	ignoredLeaves := map[string]bool{
		// Persistence metadata, not a setting: every in-memory state is at the
		// current version because the mapper stamps it on the way through.
		"SchemaOptions.SchemaVersion":          true,
		"ManualEditSettings.ManualZones":       true,
		"ManualEditSettings.ManualConnections": true,
	}
	for _, leaf := range modelLeafFields() {
		t.Run(settingPath(leaf)+"Mutated_ReportsNotEqual", func(t *testing.T) {
			t.Parallel()
			// Arrange
			left := fuzzedEditorState()
			right := deepCloneEditorState(left)
			mutateFieldValue(t, reflect.ValueOf(&right).Elem().FieldByIndex(leaf.indexes))

			// Act
			equal := left.EqualsIgnoringManualEdits(&right)

			// Assert
			assert.Equal(t, ignoredLeaves[settingPath(leaf)], equal)
		})
	}
}

// TestWhenTheLeafWalkIsBuilt_ReachesEveryFieldOfEveryGroup pins what the guard
// above actually covers. It fails if the walk stops at a group instead of
// descending into it, and it fails when a group gains or loses a field - which
// is the moment the comparison has to be extended.
func TestWhenTheLeafWalkIsBuilt_ReachesEveryFieldOfEveryGroup(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedFieldsPerGroup := map[string]int{
		"SchemaOptions":       1,
		"TemplateIdentity":    2,
		"MapSettings":         2,
		"PlayerSettings":      4,
		"NeutralZoneSettings": 11,
		"CastleSettings":      10,
		"GenerationSettings":  15,
		"GameRuleSettings":    16,
		"ContentSettings":     10,
		"ManualEditSettings":  2,
	}

	// Act
	fieldsPerGroup := map[string]int{}
	for _, leaf := range modelLeafFields() {
		groupName, _, _ := strings.Cut(settingPath(leaf), ".")
		fieldsPerGroup[groupName]++
	}

	// Assert
	assert.Equal(t, expectedFieldsPerGroup, fieldsPerGroup)
}

// leafField is one mutable field of the model, addressed through the embedded
// group that declares it.
type leafField struct {
	path    string
	indexes []int
}

// modelEmbedPrefix is the model's own embed of the persisted entity. Stripping
// it keeps every path below named after the settings group that declares it.
const modelEmbedPrefix = "EditorState."

func settingPath(leaf leafField) string {
	return strings.TrimPrefix(leaf.path, modelEmbedPrefix)
}

// modelLeafFields lists every field of every embedded entity group, as
// "Group.Field" paths. Only anonymous fields are descended into, so the walk
// stops at the leaves the hand-rolled comparison is written against.
func modelLeafFields() []leafField {
	return appendLeafFields(nil, reflect.TypeFor[editor_state_model.EditorState](), nil, "")
}

func appendLeafFields(leaves []leafField, structType reflect.Type, indexes []int, pathPrefix string) []leafField {
	for fieldIndex := range structType.NumField() {
		field := structType.Field(fieldIndex)
		fieldIndexes := append(slices.Clone(indexes), fieldIndex)
		if field.Anonymous {
			leaves = appendLeafFields(leaves, field.Type, fieldIndexes, pathPrefix+field.Name+".")
			continue
		}
		leaves = append(leaves, leafField{path: pathPrefix + field.Name, indexes: fieldIndexes})
	}
	return leaves
}

// fuzzedEditorState builds a state with fuzzed scalar fields, bonuses and
// content rows (including pointer rule fields) on top of the defaults.
// MediumNeutralContentRows stays nil and the second high-neutral row keeps
// nil Rules so nil-versus-empty mutations exercise both directions.
func fuzzedEditorState() editor_state_model.EditorState {
	state := editor_state_model.NewDefaultEditorStateModel()
	state.TemplateName = gofakeit.ProductName()
	state.PlayerCount = gofakeit.Number(2, 8)
	state.NeutralZoneCount = gofakeit.Number(0, 16)
	state.PlayerZoneSize = gofakeit.Float64Range(0.5, 2)
	state.AdvancedMode = gofakeit.Bool()
	state.BannedItems = gofakeit.LetterN(20)
	state.Bonuses = []config.BonusEntry{{
		PresetType:     config.BonusStartingGold,
		ReceiverFilter: "start_hero",
		Param:          gofakeit.DigitN(4),
	}}
	state.HighNeutralContentRows = []editor_state_model.ZoneContentRow{
		{
			Sid:   gofakeit.LetterN(12),
			Count: gofakeit.Number(1, 5),
			Rules: []editor_state_model.ContentRuleRow{{
				Name:            "Guarded",
				IsGuarded:       new(gofakeit.Bool()),
				IsSoloEncounter: new(gofakeit.Bool()),
				VariantID:       new(gofakeit.Number(0, 3)),
			}},
		},
		{Sid: gofakeit.LetterN(12), Count: 1, IsGroup: true},
	}
	return state
}

// deepCloneEditorState is the production deep copy, behind a local name so the
// comparison tests read as "these two states are independent copies". Its own
// isolation guarantees are pinned in clone_test.go.
func deepCloneEditorState(source editor_state_model.EditorState) editor_state_model.EditorState {
	return source.Clone()
}

// mutateFieldValue changes a struct field to a value that differs from the
// current one, regardless of the field's kind.
func mutateFieldValue(t *testing.T, fieldValue reflect.Value) {
	t.Helper()
	//nolint:exhaustive // kinds not present in the DTO fail via the default case
	switch fieldValue.Kind() {
	case reflect.Bool:
		fieldValue.SetBool(!fieldValue.Bool())
	case reflect.Int:
		fieldValue.SetInt(fieldValue.Int() + 1)
	case reflect.Float64:
		fieldValue.SetFloat(fieldValue.Float() + 1)
	case reflect.String:
		fieldValue.SetString(fieldValue.String() + "x")
	case reflect.Slice:
		zeroElement := reflect.Zero(fieldValue.Type().Elem())
		fieldValue.Set(reflect.Append(fieldValue, zeroElement))
	default:
		t.Fatalf("unsupported field kind %s - extend mutateFieldValue", fieldValue.Kind())
	}
}

// reflectDeepEqualIgnoringManualEdits is the pre-optimization reference
// implementation the hand-rolled comparison must stay equivalent to.
func reflectDeepEqualIgnoringManualEdits(left, right editor_state_model.EditorState) bool {
	left.ManualZones, left.ManualConnections = nil, nil
	right.ManualZones, right.ManualConnections = nil, nil
	return reflect.DeepEqual(left, right)
}
