package assetProvider_test

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// TestWhenAssetsFolderIsScanned_EveryEmbeddedSpriteIsDecoded guards against
// artwork rotting unused in the embed folder: the arena sprites sat there
// unreferenced while the feature was only half implemented.
func TestWhenAssetsFolderIsScanned_EveryEmbeddedSpriteIsDecoded(t *testing.T) {
	t.Parallel()
	// Arrange
	assetsDirectory, err := filepath.Abs(filepath.Join(
		"..", "..", "..", "..", "..", "..", "internal", "services", "asset_provider", "assets"))
	require.NoError(t, err)
	entries, err := os.ReadDir(assetsDirectory)
	require.NoError(t, err)

	spriteNames := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".png" {
			spriteNames = append(spriteNames, entry.Name())
		}
	}
	sort.Strings(spriteNames)

	// Act
	decodedNames := expectedDecodedSpriteNames()

	// Assert
	assert.Equal(t, decodedNames, spriteNames)
}

// expectedDecodedSpriteNames mirrors every sprite buildAssetProvider decodes.
func expectedDecodedSpriteNames() []string {
	names := []string{"background.png", "gladiator_arena.png"}
	for index := 1; index <= 8; index++ {
		names = append(names, fmt.Sprintf("player_%d.png", index))
	}
	for _, quality := range []string{"none", "low", "medium", "high", "highest"} {
		names = append(names,
			"neutral_"+quality+".png",
			"neutral_"+quality+"_castle.png",
			"neutral_"+quality+"_arena.png")
	}
	sort.Strings(names)
	return names
}
