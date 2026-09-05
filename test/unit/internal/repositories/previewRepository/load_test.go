package previewRepository_test

import (
	"image"
	"path/filepath"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/common/common_errors"
	"github.com/Tariomka/hommoe_custom_templates/internal/repositories"
	"github.com/stretchr/testify/assert"
)

func TestWhenPreviewLoadIsRequested_ReportsNotImplemented(t *testing.T) {
	t.Parallel()
	// Arrange
	previewPath := filepath.Join(t.TempDir(), "T.png")

	// Act
	err := repositories.NewPreviewRepository().Load(previewPath, &image.RGBA{})

	// Assert
	assert.ErrorIs(t, err, common_errors.ErrNotImplemented)
}
