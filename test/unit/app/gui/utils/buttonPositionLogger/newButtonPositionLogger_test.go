package buttonPositionLogger_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/app/gui/utils"
	"github.com/stretchr/testify/assert"
)

func TestWhenCreated_ReturnsNonNilLogger(t *testing.T) {
	t.Parallel()
	// Arrange
	handler := newRecordingHandler(levelDebug)

	// Act
	logger := utils.NewButtonPositionLogger(newSlogLogger(handler))

	// Assert
	assert.NotNil(t, logger)
}
