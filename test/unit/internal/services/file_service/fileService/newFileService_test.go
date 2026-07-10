package fileService_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/file_service"
	"github.com/stretchr/testify/assert"
)

func TestWhenServiceIsConstructed_ReturnsNonNilInstance(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	service := file_service.NewFileService()

	// Assert
	assert.NotNil(t, service)
}
