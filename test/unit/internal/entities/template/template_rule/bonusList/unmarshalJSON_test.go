package bonusList_test

import (
	"encoding/json"
	"testing"

	"github.com/Tariomka/hommoe_custom_templates/internal/entities/template/template_rule"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWhenDataIsSingleObject_WrapsBonusIntoSingleElementList(t *testing.T) {
	// Arrange
	bonusSid := gofakeit.Word()
	data := []byte(`{"sid":"` + bonusSid + `","receiverSide":-1,"parameters":["1"]}`)
	expected := template_rule.BonusList{
		{SID: bonusSid, ReceiverSide: -1, Parameters: []string{"1"}},
	}
	var list template_rule.BonusList

	// Act
	err := json.Unmarshal(data, &list)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, list)
}

func TestWhenDataIsArray_DecodesEveryBonus(t *testing.T) {
	// Arrange
	data := []byte(
		`[{"sid":"first","receiverSide":0,"parameters":[]},{"sid":"second","receiverSide":1,"parameters":["7"]}]`,
	)
	expected := template_rule.BonusList{
		{SID: "first", ReceiverSide: 0, Parameters: []string{}},
		{SID: "second", ReceiverSide: 1, Parameters: []string{"7"}},
	}
	var list template_rule.BonusList

	// Act
	err := json.Unmarshal(data, &list)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expected, list)
}

func TestWhenDataIsNull_SetsListToNil(t *testing.T) {
	// Arrange
	list := template_rule.BonusList{{SID: gofakeit.Word()}}

	// Act
	err := list.UnmarshalJSON([]byte("null"))

	// Assert
	require.NoError(t, err)
	assert.Nil(t, list)
}

func TestWhenDataIsEmpty_SetsListToNil(t *testing.T) {
	// Arrange
	list := template_rule.BonusList{{SID: gofakeit.Word()}}

	// Act
	err := list.UnmarshalJSON([]byte{})

	// Assert
	require.NoError(t, err)
	assert.Nil(t, list)
}

func TestWhenSingleObjectHasInvalidFieldType_ReturnsError(t *testing.T) {
	// Arrange
	data := []byte(`{"receiverSide":"notANumber"}`)
	var list template_rule.BonusList

	// Act
	err := json.Unmarshal(data, &list)

	// Assert
	assert.Error(t, err)
}

func TestWhenArrayElementHasInvalidFieldType_ReturnsError(t *testing.T) {
	// Arrange
	data := []byte(`[{"receiverSide":"notANumber"}]`)
	var list template_rule.BonusList

	// Act
	err := json.Unmarshal(data, &list)

	// Assert
	assert.Error(t, err)
}
