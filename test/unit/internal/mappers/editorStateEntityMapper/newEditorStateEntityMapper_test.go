package editorStateEntityMapper_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/mappers"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheEditorStateEntityMapperIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	mapper := mappers.NewEditorStateEntityMapper()

	// Assert
	assert.NotNil(t, mapper)
}
