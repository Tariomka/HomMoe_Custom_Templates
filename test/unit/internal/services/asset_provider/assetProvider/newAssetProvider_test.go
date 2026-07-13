package assetProvider_test

import (
	"sync"
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

func TestWhenCalledConcurrently_ReturnsSameSingletonInstance(t *testing.T) {
	t.Parallel()
	// Arrange
	const goroutineCount = 16
	providers := make([]*asset_provider.AssetProvider, goroutineCount)
	var waitGroup sync.WaitGroup
	waitGroup.Add(goroutineCount)

	// Act
	for index := range goroutineCount {
		go func() {
			defer waitGroup.Done()
			provider, err := asset_provider.NewAssetProvider()
			assert.NoError(t, err)
			providers[index] = provider
		}()
	}
	waitGroup.Wait()

	// Assert
	for index := 1; index < goroutineCount; index++ {
		assert.Same(t, providers[0], providers[index])
	}
}
