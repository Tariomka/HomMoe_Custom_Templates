package mandatoryContentItemMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/stretchr/testify/assert"
)

func TestWhenMapperIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	mapper := mappers.NewMandatoryContentItemMapper()

	// Assert
	assert.NotNil(t, mapper)
}
