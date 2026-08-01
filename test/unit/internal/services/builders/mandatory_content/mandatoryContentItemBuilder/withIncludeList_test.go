package mandatoryContentItemBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/mandatory_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenIncludeListIsProvided_AppendsIncludeList(t *testing.T) {
	t.Parallel()
	// Arrange
	includeList := gofakeit.Word()
	builder := mandatory_content.NewContentItemBuilder("")

	// Act
	item := builder.WithIncludeList(includeList).Build()

	// Assert
	assert.Equal(t, entities.MandatoryContentItem{IncludeLists: []string{includeList}}, item)
}
