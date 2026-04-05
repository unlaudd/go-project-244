package formatters

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatJSON(t *testing.T) {
	nodes := []DiffNode{
		{Key: "common", State: "changed", Children: []DiffNode{
			{Key: "follow", State: "added", Value: false},
			{Key: "setting1", State: "unchanged", Value: "Value 1"},
			{Key: "setting2", State: "removed", Value: 200.0},
			{Key: "setting3", State: "changed", OldValue: true, NewValue: nil},
			{Key: "setting5", State: "added", Value: map[string]interface{}{"key5": "value5"}},
		}},
		{Key: "group2", State: "removed", Value: map[string]interface{}{"abc": 12345.0}},
	}

	result := FormatJSON(nodes)

	// 1. Проверяем, что это валидный JSON
	var parsed map[string]interface{}
	err := json.Unmarshal([]byte(result), &parsed)
	require.NoError(t, err, "Output should be valid JSON")

	// 2. Проверяем структуру корневого ключа "common"
	common, ok := parsed["common"].(map[string]interface{})
	require.True(t, ok, "'common' should be an object")
	assert.Equal(t, "changed", common["status"])

	// 3. Дети находятся внутри "children"
	children, ok := common["children"].(map[string]interface{})
	require.True(t, ok, "'common' should have 'children' object")

	// 4. Проверяем "follow" внутри children
	follow, ok := children["follow"].(map[string]interface{})
	require.True(t, ok, "'follow' should be an object")
	assert.Equal(t, "added", follow["status"])
	assert.Equal(t, false, follow["value"])

	// 5. Проверяем "setting3" (changed с oldValue/newValue)
	setting3, ok := children["setting3"].(map[string]interface{})
	require.True(t, ok)
	assert.Equal(t, "changed", setting3["status"])
	assert.Equal(t, true, setting3["oldValue"])
	assert.Nil(t, setting3["newValue"])

	// 6. Проверяем, что в выводе есть ключевые поля
	assert.Contains(t, result, `"status": "added"`)
	assert.Contains(t, result, `"value": false`)
	assert.Contains(t, result, `"children"`)
}

func TestFormatJSONNestedObject(t *testing.T) {
	// Проверяем рекурсивную обработку вложенных объектов
	nodes := []DiffNode{
		{Key: "root", State: "changed", Children: []DiffNode{
			{Key: "nested", State: "changed", Children: []DiffNode{
				{Key: "deep", State: "added", Value: "value"},
			}},
		}},
	}

	result := FormatJSON(nodes)

	var parsed map[string]interface{}
	require.NoError(t, json.Unmarshal([]byte(result), &parsed))

	root := parsed["root"].(map[string]interface{})
	children := root["children"].(map[string]interface{})
	nested := children["nested"].(map[string]interface{})
	nestedChildren := nested["children"].(map[string]interface{})
	deep := nestedChildren["deep"].(map[string]interface{})

	assert.Equal(t, "added", deep["status"])
	assert.Equal(t, "value", deep["value"])
}
