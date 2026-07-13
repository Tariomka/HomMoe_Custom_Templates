package editorStateDto_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStatesAreFullyIdentical_ReportsEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := dtos.NewDefaultEditorStateDto()
	right := left

	// Act
	equal := left.EqualsIgnoringManualEdits(&right)

	// Assert
	assert.True(t, equal)
}

func TestWhenOnlyManualEditFieldsDiffer_ReportsEqual(t *testing.T) {
	t.Parallel()
	// Arrange
	left := dtos.NewDefaultEditorStateDto()
	right := left
	right.ManualZones = []editor_state_dto.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}
	right.ManualConnections = []editor_state_dto.ManualConnectionSave{
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
	left := dtos.NewDefaultEditorStateDto()
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
	right.MediumNeutralContentRows = []models.ZoneContentRowSave{}

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
		mutate      func(state *dtos.EditorStateDto)
	}{
		{"StatesAreDeepClones_MatchesDeepEqual", func(_ *dtos.EditorStateDto) {}},
		{"IntFieldDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) { state.PlayerCount++ }},
		{"FloatFieldDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) { state.PlayerZoneSize += 0.25 }},
		{"BoolFieldDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.AdvancedMode = !state.AdvancedMode
		}},
		{"TopologyDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.Topology = config.TopologyHubAndSpoke
		}},
		{"BonusEntryDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) { state.Bonuses[0].Param += "0" }},
		{"BonusesNilOnOneSide_MatchesDeepEqual", func(state *dtos.EditorStateDto) { state.Bonuses = nil }},
		{"RowAppended_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.HighNeutralContentRows = append(state.HighNeutralContentRows, models.ZoneContentRowSave{Sid: "extra"})
		}},
		{
			"RowCountFieldDiffers_MatchesDeepEqual",
			func(state *dtos.EditorStateDto) { state.HighNeutralContentRows[0].Count++ },
		},
		{
			"RowSidDiffers_MatchesDeepEqual",
			func(state *dtos.EditorStateDto) { state.HighNeutralContentRows[0].Sid += "x" },
		},
		{"RuleNameDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.HighNeutralContentRows[0].Rules[0].Name += "x"
		}},
		{"RuleGuardedValueDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			*state.HighNeutralContentRows[0].Rules[0].IsGuarded = !*state.HighNeutralContentRows[0].Rules[0].IsGuarded
		}},
		{"RuleGuardedNilOnOneSide_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.HighNeutralContentRows[0].Rules[0].IsGuarded = nil
		}},
		{"RuleSoloEncounterValueDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			*state.HighNeutralContentRows[0].Rules[0].IsSoloEncounter = !*state.HighNeutralContentRows[0].Rules[0].IsSoloEncounter
		}},
		{"RuleVariantIDValueDiffers_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			*state.HighNeutralContentRows[0].Rules[0].VariantID++
		}},
		{"RulesNilVersusEmpty_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.HighNeutralContentRows[1].Rules = []models.ContentRuleRowSave{}
		}},
		{"RowsNilVersusEmpty_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.MediumNeutralContentRows = []models.ZoneContentRowSave{}
		}},
		{"ManualZonesDiffer_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.ManualZones = []editor_state_dto.ManualZoneSave{{Zone: entities.Zone{Name: "Zone A"}}}
		}},
		{"ManualConnectionsDiffer_MatchesDeepEqual", func(state *dtos.EditorStateDto) {
			state.ManualConnections = []editor_state_dto.ManualConnectionSave{
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

// TestWhenAnyNonManualFieldIsMutated_ReportsNotEqual walks every field of the
// DTO via reflection, mutates it on a deep clone and expects inequality
// (except for the ignored manual-edit fields). It trips when a new field is
// added to EditorStateDto without extending the hand-rolled comparison.
func TestWhenAnyNonManualFieldIsMutated_ReportsNotEqual(t *testing.T) {
	t.Parallel()
	ignoredFields := map[string]bool{"ManualZones": true, "ManualConnections": true}
	stateType := reflect.TypeFor[dtos.EditorStateDto]()
	for fieldIndex := range stateType.NumField() {
		field := stateType.Field(fieldIndex)
		t.Run(field.Name+"Mutated_ReportsNotEqual", func(t *testing.T) {
			t.Parallel()
			// Arrange
			left := fuzzedEditorState()
			right := deepCloneEditorState(left)
			mutateFieldValue(t, reflect.ValueOf(&right).Elem().Field(fieldIndex))

			// Act
			equal := left.EqualsIgnoringManualEdits(&right)

			// Assert
			assert.Equal(t, ignoredFields[field.Name], equal)
		})
	}
}

// fuzzedEditorState builds a state with fuzzed scalar fields, bonuses and
// content rows (including pointer rule fields) on top of the defaults.
// MediumNeutralContentRows stays nil and the second high-neutral row keeps
// nil Rules so nil-versus-empty mutations exercise both directions.
func fuzzedEditorState() dtos.EditorStateDto {
	state := dtos.NewDefaultEditorStateDto()
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
	state.HighNeutralContentRows = []models.ZoneContentRowSave{
		{
			Sid:   gofakeit.LetterN(12),
			Count: gofakeit.Number(1, 5),
			Rules: []models.ContentRuleRowSave{{
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

func deepCloneEditorState(source dtos.EditorStateDto) dtos.EditorStateDto {
	cloned := source
	cloned.Bonuses = slices.Clone(source.Bonuses)
	cloned.PlayerZoneContentRows = deepCloneContentRows(source.PlayerZoneContentRows)
	cloned.LowNeutralContentRows = deepCloneContentRows(source.LowNeutralContentRows)
	cloned.MediumNeutralContentRows = deepCloneContentRows(source.MediumNeutralContentRows)
	cloned.HighNeutralContentRows = deepCloneContentRows(source.HighNeutralContentRows)
	cloned.HubZoneContentRows = deepCloneContentRows(source.HubZoneContentRows)
	return cloned
}

func deepCloneContentRows(rows []models.ZoneContentRowSave) []models.ZoneContentRowSave {
	cloned := slices.Clone(rows)
	for rowIndex := range cloned {
		cloned[rowIndex].Rules = slices.Clone(cloned[rowIndex].Rules)
		for ruleIndex := range cloned[rowIndex].Rules {
			rule := &cloned[rowIndex].Rules[ruleIndex]
			rule.IsGuarded = clonePointer(rule.IsGuarded)
			rule.IsSoloEncounter = clonePointer(rule.IsSoloEncounter)
			rule.VariantID = clonePointer(rule.VariantID)
		}
	}
	return cloned
}

func clonePointer[Value any](source *Value) *Value {
	if source == nil {
		return nil
	}
	value := *source
	return &value
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
func reflectDeepEqualIgnoringManualEdits(left, right dtos.EditorStateDto) bool {
	left.ManualZones, left.ManualConnections = nil, nil
	right.ManualZones, right.ManualConnections = nil, nil
	return reflect.DeepEqual(left, right)
}
