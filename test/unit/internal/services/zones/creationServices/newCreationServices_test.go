package creationServices_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/zones"
	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreOmitted_CreatesCompleteServiceSet(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	services := zones.NewCreationServices(nil, nil)

	// Assert
	assert.NotNil(t, services.ZoneFactory)
}

func TestWhenDependenciesAreOmitted_CreatesCastleFactory(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	services := zones.NewCreationServices(nil, nil)

	// Assert
	assert.NotNil(t, services.CastleFactory)
}

func TestWhenDependenciesAreOmitted_CreatesRoadFactory(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	services := zones.NewCreationServices(nil, nil)

	// Assert
	assert.NotNil(t, services.RoadFactory)
}

func TestWhenCastleFactoryIsProvided_PreservesInstance(t *testing.T) {
	t.Parallel()
	// Arrange
	castleFactory := zones.NewCastleFactory()

	// Act
	services := zones.NewCreationServices(castleFactory, nil)

	// Assert
	assert.Same(t, castleFactory, services.CastleFactory)
}

func TestWhenRoadFactoryIsProvided_PreservesInstance(t *testing.T) {
	t.Parallel()
	// Arrange
	roadFactory := zones.NewRoadFactory()

	// Act
	services := zones.NewCreationServices(nil, roadFactory)

	// Assert
	assert.Same(t, roadFactory, services.RoadFactory)
}
