package generatorConfigMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapperIsCreated_ReturnsNonNilInstance(t *testing.T) {
	// Arrange - nothing to prepare.

	// Act
	mapper := mappers.NewConfigMapper()

	// Assert
	assert.NotNil(t, mapper)
}
