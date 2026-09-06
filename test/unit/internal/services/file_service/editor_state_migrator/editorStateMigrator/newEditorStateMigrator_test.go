package editorStateMigrator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenTheMigratorIsCreated_ReturnsUsableInstance(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	migrator := newMigrator()

	// Assert
	assert.NotNil(t, migrator)
}
