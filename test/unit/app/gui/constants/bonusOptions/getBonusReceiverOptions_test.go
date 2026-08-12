package bonusOptions_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/constants"
	"github.com/Tariomka/hommoe_custom_templates/internal/registry"
	"github.com/stretchr/testify/assert"
)

func TestWhenBonusReceiversAreRequested_ReturnsTheStartingHeroAndAllHeroesFilters(t *testing.T) {
	t.Parallel()
	// Arrange
	receiversFilters := registry.GetReceiversFilterValues()

	// Act
	receivers := constants.GetBonusReceiverOptions()

	// Assert
	assert.Equal(t, []string{receiversFilters.StartingHero, receiversFilters.AllHeroes}, receivers)
}
