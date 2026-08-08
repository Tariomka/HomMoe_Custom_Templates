package bonusEntryService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/bonuses"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheServiceIsConstructed_ReturnsAUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	service := bonuses.NewBonusEntryService()

	// Assert
	assert.NotNil(t, service)
}
