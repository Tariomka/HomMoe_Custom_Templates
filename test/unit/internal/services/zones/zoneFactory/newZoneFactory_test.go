package zoneFactory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreProvided_ReturnsInstance(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	factory := newZoneFactory()

	// Assert
	assert.NotNil(t, factory)
}
