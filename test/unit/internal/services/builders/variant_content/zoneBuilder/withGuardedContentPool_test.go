package zoneBuilder_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/template_model"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/builders/variant_content"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestWhenGuardedContentPoolIsProvided_SetsGuardedContentPoolOnBuiltZone(t *testing.T) {
	t.Parallel()
	// Arrange
	expectedPool := []string{gofakeit.Word(), gofakeit.Word()}
	builder := variant_content.NewZoneBuilder()

	// Act
	zone := builder.WithGuardedContentPool(expectedPool).Build()

	// Assert
	assert.Equal(t, template_model.Zone{GuardedContentPool: expectedPool}, zone)
}
