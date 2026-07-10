package topologyProvider_test

import (
	"slices"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/models/config"
	"github.com/Tariomka/hommoe_custom_templates/internal/services/template_generator/providers"
	"github.com/stretchr/testify/assert"
)

func TestWhenCalled_ReturnsSameProviderForChaining(t *testing.T) {
	t.Parallel()
	// Arrange
	provider := providers.NewTopologyProvider()

	// Act
	returned := provider.ShufflePlayerZones(true)

	// Assert
	assert.Same(t, provider, returned)
}

func TestWhenShuffleEnabled_DoesNotMutateInputLabels(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D", "E", "F", "G", "H"}
	expectedLabels := slices.Clone(playerLabels)
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := providers.NewTopologyProvider().ShufflePlayerZones(true)

	// Act
	provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.Equal(t, expectedLabels, playerLabels)
}

func TestWhenShuffleEnabled_PreservesSpawnZoneNameSet(t *testing.T) {
	t.Parallel()
	// Arrange
	playerLabels := []string{"A", "B", "C", "D"}
	configuration := config.NewGeneratorConfig()
	configuration.Topology = config.TopologyRing
	tuning := buildVariantInputs(configuration, playerLabels, nil)
	provider := providers.NewTopologyProvider().ShufflePlayerZones(true)

	// Act
	variant := provider.CreateTopologyVariant(*configuration, playerLabels, nil, tuning, "")

	// Assert
	assert.Equal(t, []string{"Spawn-A", "Spawn-B", "Spawn-C", "Spawn-D"}, spawnZoneNames(variant))
}
