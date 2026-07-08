package winConditions_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_rule"
	"github.com/stretchr/testify/assert"
)

func TestWhenDestinationIsEmpty_CopiesEverySourceField(t *testing.T) {
	// Arrange
	destination := template_rule.WinConditions{}
	source := template_rule.WinConditions{Classic: true, CityHold: true, CityHoldDays: 9}

	// Act
	destination.MergeWinConditionsIfDoesNotExist(source)

	// Assert
	assert.Equal(t, source, destination)
}

func TestWhenDestinationFieldIsSet_KeepsDestinationValue(t *testing.T) {
	// Arrange
	destination := template_rule.WinConditions{CityHoldDays: 6}
	source := template_rule.WinConditions{CityHoldDays: 9}

	// Act
	destination.MergeWinConditionsIfDoesNotExist(source)

	// Assert
	assert.Equal(t, 6, destination.CityHoldDays)
}

func TestWhenSourceIsEmpty_LeavesDestinationUnchanged(t *testing.T) {
	// Arrange
	destination := template_rule.WinConditions{Classic: true, DesertionDay: 3}
	expected := destination

	// Act
	destination.MergeWinConditionsIfDoesNotExist(template_rule.WinConditions{})

	// Assert
	assert.Equal(t, expected, destination)
}

func TestWhenFieldsOverlapPartially_FillsOnlyMissingFields(t *testing.T) {
	// Arrange
	destination := template_rule.WinConditions{CityHoldDays: 6}
	source := template_rule.WinConditions{CityHoldDays: 9, Desertion: true, DesertionDay: 3}
	expected := template_rule.WinConditions{CityHoldDays: 6, Desertion: true, DesertionDay: 3}

	// Act
	destination.MergeWinConditionsIfDoesNotExist(source)

	// Assert
	assert.Equal(t, expected, destination)
}

func TestWhenSourceHasTournamentDaySlices_CopiesSlices(t *testing.T) {
	// Arrange
	destination := template_rule.WinConditions{}
	source := template_rule.WinConditions{
		Tournament:             true,
		TournamentDays:         []int{7, 14},
		TournamentAnnounceDays: []int{5, 12},
	}

	// Act
	destination.MergeWinConditionsIfDoesNotExist(source)

	// Assert
	assert.Equal(t, source, destination)
}

func TestWhenSourceHasOnlyFieldsPresentInDestination_LeavesDestinationUnchanged(t *testing.T) {
	// Arrange - every source key already exists in destination, so nothing changes.
	destination := template_rule.WinConditions{Classic: true, CityHoldDays: 6}
	expected := destination
	source := template_rule.WinConditions{Classic: true, CityHoldDays: 9}

	// Act
	destination.MergeWinConditionsIfDoesNotExist(source)

	// Assert
	assert.Equal(t, expected, destination)
}
