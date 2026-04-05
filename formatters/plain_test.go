package formatters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPlain(t *testing.T) {
	nodes := []DiffNode{
		{Key: "common", State: "changed", Children: []DiffNode{
			{Key: "follow", State: "added", Value: false},
			{Key: "setting2", State: "removed", Value: 200.0},
			{Key: "setting3", State: "changed", OldValue: true, NewValue: nil},
			{Key: "setting5", State: "added", Value: map[string]interface{}{"key5": "value5"}},
			{Key: "setting6", State: "changed", Children: []DiffNode{
				{Key: "doge", State: "changed", Children: []DiffNode{
					{Key: "wow", State: "changed", OldValue: "", NewValue: "so much"},
				}},
			}},
		}},
		{Key: "group1", State: "changed", Children: []DiffNode{
			{Key: "nest", State: "changed", OldValue: map[string]interface{}{"key": "value"}, NewValue: "str"},
		}},
		{Key: "group2", State: "removed", Value: map[string]interface{}{"abc": 12345.0}},
	}

	result := FormatPlain(nodes)
	lines := strings.Split(result, "\n")

	assert.True(t, len(lines) > 0, "Output should not be empty")
	assert.Contains(t, result, "Property 'common.follow' was added with value: false")
	assert.Contains(t, result, "Property 'common.setting2' was removed")
	assert.Contains(t, result, "Property 'common.setting3' was updated. From true to null")
	assert.Contains(t, result, "Property 'common.setting5' was added with value: [complex value]")
	assert.Contains(t, result, "Property 'common.setting6.doge.wow' was updated. From '' to 'so much'")
	assert.Contains(t, result, "Property 'group1.nest' was updated. From [complex value] to 'str'")
	assert.Contains(t, result, "Property 'group2' was removed")
}

func TestFormatPlainValue(t *testing.T) {
	assert.Equal(t, "null", formatPlainValue(nil))
	assert.Equal(t, "'test'", formatPlainValue("test"))
	assert.Equal(t, "true", formatPlainValue(true))
	assert.Equal(t, "100", formatPlainValue(100.0))
	assert.Equal(t, "99.9", formatPlainValue(99.9))
	assert.Equal(t, "[complex value]", formatPlainValue(map[string]interface{}{"a": 1}))
	assert.Equal(t, "[complex value]", formatPlainValue([]interface{}{1, 2})) // default ветка
}
