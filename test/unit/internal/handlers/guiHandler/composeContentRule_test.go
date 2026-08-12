package guiHandler_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos"
	"github.com/stretchr/testify/assert"
)

func TestWhenTheEditorSelectionIsUnknown_ComposeContentRuleReportsItInvalid(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newProductionGuiHandler()

	// Act
	result := handler.ComposeContentRule(dtos.ContentRuleCompositionRequestDto{})

	// Assert
	assert.False(t, result.Valid)
}
