package testLayoutChecker_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/cmd/testlayoutcheck/checker"
	"github.com/stretchr/testify/assert"
)

func TestWhenCheckerIsCreated_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	instance := checker.NewTestLayoutChecker()

	// Assert
	assert.NotNil(t, instance)
}
