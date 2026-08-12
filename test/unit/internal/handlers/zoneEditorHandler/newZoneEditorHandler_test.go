package zoneEditorHandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenHandlerIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	fixture := newZoneEditorHandlerFixture()

	// Assert
	assert.NotNil(t, fixture.handler)
}
