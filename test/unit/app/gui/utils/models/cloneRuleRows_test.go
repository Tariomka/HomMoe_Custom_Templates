package models_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenRuleListIsNil_NilIsReturned(t *testing.T) {
	// Arrange & Act
	cloned := utils.CloneRuleRows(nil)

	// Assert
	assert.Nil(t, cloned)
}

func TestWhenRuleListIsEmpty_NilIsReturned(t *testing.T) {
	// Arrange
	rules := []models.ContentRuleRowSave{}

	// Act
	cloned := utils.CloneRuleRows(rules)

	// Assert
	assert.Nil(t, cloned)
}

func TestWhenRulesArePresent_EqualCopyIsReturned(t *testing.T) {
	// Arrange
	isGuarded := gofakeit.Bool()
	rules := []models.ContentRuleRowSave{
		{Name: "Guarded", IsGuarded: &isGuarded},
		{Name: "Distance to town", DistanceName: "Near"},
	}

	// Act
	cloned := utils.CloneRuleRows(rules)

	// Assert
	assert.Equal(t, rules, cloned)
}

func TestWhenCloneIsMutated_OriginalRulesStayUnchanged(t *testing.T) {
	// Arrange
	rules := []models.ContentRuleRowSave{{Name: "Guarded"}}
	cloned := utils.CloneRuleRows(rules)

	// Act
	cloned[0].Name = gofakeit.Name()

	// Assert
	assert.Equal(t, "Guarded", rules[0].Name)
}
