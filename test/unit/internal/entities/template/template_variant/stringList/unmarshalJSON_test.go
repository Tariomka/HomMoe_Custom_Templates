package stringList_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_variant"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDataIsSingleString_WrapsIntoSingleElementList(t *testing.T) {
	// Arrange
	value := gofakeit.Word()
	data := []byte(`"` + value + `"`)
	var list template_variant.StringList

	// Act
	err := json.Unmarshal(data, &list)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, template_variant.StringList{value}, list)
}

func TestWhenDataIsArray_DecodesEveryString(t *testing.T) {
	// Arrange
	data := []byte(`["first","second"]`)
	var list template_variant.StringList

	// Act
	err := json.Unmarshal(data, &list)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, template_variant.StringList{"first", "second"}, list)
}

func TestWhenDataIsNull_SetsListToNil(t *testing.T) {
	// Arrange
	list := template_variant.StringList{gofakeit.Word()}

	// Act
	err := list.UnmarshalJSON([]byte("null"))

	// Assert
	require.NoError(t, err)
	assert.Nil(t, list)
}

func TestWhenDataIsEmpty_SetsListToNil(t *testing.T) {
	// Arrange
	list := template_variant.StringList{gofakeit.Word()}

	// Act
	err := list.UnmarshalJSON([]byte{})

	// Assert
	require.NoError(t, err)
	assert.Nil(t, list)
}

func TestWhenSingleStringIsUnterminated_ReturnsError(t *testing.T) {
	// Arrange
	list := template_variant.StringList{}

	// Act
	err := list.UnmarshalJSON([]byte(`"unterminated`))

	// Assert
	assert.Error(t, err)
}

func TestWhenArrayHoldsNonStringElements_ReturnsError(t *testing.T) {
	// Arrange
	data := []byte(`[1,2]`)
	var list template_variant.StringList

	// Act
	err := json.Unmarshal(data, &list)

	// Assert
	assert.Error(t, err)
}
