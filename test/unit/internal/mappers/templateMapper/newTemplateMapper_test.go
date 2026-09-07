package templateMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheTemplateMapperIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	mapper := mappers.NewTemplateMapper()

	// Assert
	assert.NotNil(t, mapper)
}
