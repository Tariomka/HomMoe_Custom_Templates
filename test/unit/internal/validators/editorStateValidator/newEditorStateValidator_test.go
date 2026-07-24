package editorStateValidator_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/validators"
	"github.com/stretchr/testify/assert"
)

func TestWhenValidatorIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	validator := validators.NewEditorStateValidator()

	// Assert
	assert.NotNil(t, validator)
}
