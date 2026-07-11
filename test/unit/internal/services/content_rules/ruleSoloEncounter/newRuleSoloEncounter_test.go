package ruleSoloEncounter_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/content_rules"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenStateIsSupplied_StoresIt(t *testing.T) {
	t.Parallel()
	// Arrange
	isSoloEncounter := gofakeit.Bool()

	// Act
	rule := content_rules.NewRuleSoloEncounter(isSoloEncounter)

	// Assert
	assert.Equal(t, isSoloEncounter, rule.IsSoloEncounter)
}
