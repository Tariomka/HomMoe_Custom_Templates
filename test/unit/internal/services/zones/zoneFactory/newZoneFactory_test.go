package zoneFactory_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreOmitted_ReturnsUsableFactory(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	factory := zones.NewZoneFactory(nil, nil)

	// Assert
	assert.NotNil(t, factory)
}
