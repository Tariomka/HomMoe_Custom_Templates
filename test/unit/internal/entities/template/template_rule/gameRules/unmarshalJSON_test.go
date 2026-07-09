package gameRules_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_rule"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenWinConditionsAreNested_DecodesNestedBlock(t *testing.T) {
	// Arrange
	data := []byte(`{"heroCountMin":2,"winConditions":{"classic":true,"cityHold":true,"cityHoldDays":6}}`)
	expected := template_rule.WinConditions{Classic: true, CityHold: true, CityHoldDays: 6}
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, rules.WinConditions)
}

func TestWhenWinConditionsAreFlatSiblings_PopulatesWinConditions(t *testing.T) {
	// Arrange - Zookeeper-style template: win-condition keys flat inside gameRules.
	data := []byte(`{"heroCountMin":2,"classic":true,"cityHold":true,"cityHoldDays":9}`)
	expected := template_rule.WinConditions{Classic: true, CityHold: true, CityHoldDays: 9}
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, rules.WinConditions)
}

func TestWhenNestedAndFlatDeclareSameField_NestedValueWins(t *testing.T) {
	// Arrange
	data := []byte(`{"cityHoldDays":9,"winConditions":{"cityHoldDays":6}}`)
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 6, rules.WinConditions.CityHoldDays)
}

func TestWhenFlatFieldIsMissingFromNestedBlock_FillsItFromFlatSibling(t *testing.T) {
	// Arrange
	data := []byte(`{"desertionDay":3,"winConditions":{"classic":true}}`)
	expected := template_rule.WinConditions{Classic: true, DesertionDay: 3}
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, rules.WinConditions)
}

func TestWhenGladiatorArenaIsFlatSibling_MirrorsIntoWinConditions(t *testing.T) {
	// Arrange - Symmetry/Jebus Outcast style: gladiator fields as siblings of winConditions.
	data := []byte(`{"gladiatorArena":true,"gladiatorArenaCountDay":3,"winConditions":{}}`)
	expected := template_rule.WinConditions{GladiatorArena: true, GladiatorArenaCountDay: 3}
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, rules.WinConditions)
}

func TestWhenScalarRuleFieldsProvided_DecodesWholeGameRules(t *testing.T) {
	// Arrange
	data := []byte(
		`{"heroCountMin":4,"heroCountMax":8,"heroCountIncrement":1,"heroHireBan":true,"encounterHoles":true,"tournamentRules":true}`,
	)
	expected := template_rule.GameRules{
		HeroCountMin:       4,
		HeroCountMax:       8,
		HeroCountIncrement: 1,
		HeroHireBan:        true,
		EncounterHoles:     true,
		TournamentRules:    true,
	}
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, rules)
}

func TestWhenGameRuleFieldHasInvalidType_ReturnsError(t *testing.T) {
	// Arrange
	data := []byte(`{"heroCountMin":"notANumber"}`)
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	assert.Error(t, err)
}

func TestWhenFlatWinConditionDecodeFails_ReturnsNoError(t *testing.T) {
	// Arrange - desertionDay is unknown to GameRules but ill-typed for WinConditions.
	data := []byte(`{"heroCountMin":4,"desertionDay":"notANumber"}`)
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	assert.NoError(t, err)
}

func TestWhenFlatWinConditionDecodeFails_LeavesWinConditionsEmpty(t *testing.T) {
	// Arrange
	data := []byte(`{"heroCountMin":4,"desertionDay":"notANumber"}`)
	var rules template_rule.GameRules

	// Act
	err := json.Unmarshal(data, &rules)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, template_rule.WinConditions{}, rules.WinConditions)
}
