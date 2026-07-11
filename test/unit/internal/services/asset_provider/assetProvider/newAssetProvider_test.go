package assetProvider_test

import (
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/services/asset_provider"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenEmbeddedAssetsAreValid_ReturnsNoError(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	_, err := asset_provider.NewAssetProvider()

	// Assert
	assert.NoError(t, err)
}

func TestWhenEmbeddedAssetsAreValid_ReturnsProvider(t *testing.T) {
	t.Parallel()
	// Arrange

	// Act
	provider, err := asset_provider.NewAssetProvider()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, provider)
}

func TestWhenCalledTwice_ReturnsSameSingletonInstance(t *testing.T) {
	t.Parallel()
	// Arrange
	firstProvider, err := asset_provider.NewAssetProvider()
	require.NoError(t, err)

	// Act
	secondProvider, err := asset_provider.NewAssetProvider()

	// Assert
	require.NoError(t, err)
	assert.Same(t, firstProvider, secondProvider)
}
