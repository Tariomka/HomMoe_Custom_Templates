package mandatoryContentProvider_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhenDependenciesAreProvided_ReturnsUsableProvider(t *testing.T) {
	t.Parallel()
	// Arrange & Act
	provider := newMandatoryContentProvider()

	// Assert
	assert.NotNil(t, provider)
}
