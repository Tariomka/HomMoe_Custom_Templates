package editorStateDto_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/dtos/editor_state_dto"
	"github.com/stretchr/testify/assert"
)

func TestWhenPlayerCountChanges_ReportsChanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_dto.NewDefaultEditorStateDto()
	incoming := previous
	incoming.PlayerCount++

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.True(t, changed)
}

func TestWhenStatesAreIdentical_ReportsUnchanged(t *testing.T) {
	t.Parallel()
	// Arrange
	previous := editor_state_dto.NewDefaultEditorStateDto()
	incoming := previous

	// Act
	changed := previous.LayoutDefiningOptionsChanged(&incoming)

	// Assert
	assert.False(t, changed)
}
