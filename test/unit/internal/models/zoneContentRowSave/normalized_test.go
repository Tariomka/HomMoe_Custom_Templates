package zoneContentRowSave_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestWhenCountIsUnset_NormalizesCountToOne(t *testing.T) {
	// Arrange
	row := models.ZoneContentRowSave{Sid: "x"}

	// Act
	normalized := row.Normalized()

	// Assert
	assert.Equal(t, 1, normalized.Count)
}

func TestWhenRulesArePresentAndCountIsUnset_StillNormalizesCountToOne(t *testing.T) {
	// Arrange
	row := models.ZoneContentRowSave{
		Sid:   "x",
		Rules: []models.ContentRuleRowSave{{Name: "Guarded", IsGuarded: new(true)}},
	}

	// Act
	normalized := row.Normalized()

	// Assert
	assert.Equal(t, 1, normalized.Count)
}

func TestWhenCountIsPositive_KeepsCountUnchanged(t *testing.T) {
	// Arrange
	row := models.ZoneContentRowSave{Sid: "x", Count: 3}

	// Act
	normalized := row.Normalized()

	// Assert
	assert.Equal(t, 3, normalized.Count)
}
