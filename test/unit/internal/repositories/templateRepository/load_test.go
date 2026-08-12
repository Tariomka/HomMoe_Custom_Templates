package templateRepository_test

import (
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestWhenTemplateLoadIsRequested_ReportsNotImplemented(t *testing.T) {
	t.Parallel()
	// Arrange
	templatePath := filepath.Join(t.TempDir(), "T.rmg.json")

	// Act
	_, err := repositories.NewTemplateRepository().Load(templatePath)

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNotImplemented)
}
