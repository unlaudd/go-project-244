package formatters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatStylish(t *testing.T) {
	nodes := []DiffNode{
		{Key: "common", State: "changed", Children: []DiffNode{
			{Key: "follow", State: "added", Value: false},
			{Key: "setting1", State: "unchanged", Value: "Value 1"},
			{Key: "setting2", State: "removed", Value: 200.0},
			{Key: "setting3", State: "changed", OldValue: true, NewValue: nil},
		}},
		{Key: "group2", State: "removed", Value: map[string]interface{}{"abc": 12345.0}},
	}

	result := FormatStylish(nodes)

	assert.True(t, strings.HasPrefix(result, "{\n"), "Should start with {")
	assert.True(t, strings.HasSuffix(result, "}"), "Should end with }")
	assert.Contains(t, result, "+ follow: false")
	assert.Contains(t, result, "  setting1: Value 1")
	assert.Contains(t, result, "- setting2: 200")
	assert.Contains(t, result, "- setting3: true")
	assert.Contains(t, result, "+ setting3: null")
	assert.Contains(t, result, "- group2: {")
}

func TestFormatStylishValue(t *testing.T) {
	// Покрываем все ветки switch
	assert.Equal(t, "null", formatStylishValue(nil))
	assert.Equal(t, "hello", formatStylishValue("hello"))
	assert.Equal(t, "true", formatStylishValue(true))
	assert.Equal(t, "false", formatStylishValue(false))
	assert.Equal(t, "42", formatStylishValue(42.0))
	assert.Equal(t, "3.14", formatStylishValue(3.14))
	assert.Contains(t, formatStylishValue([]int{1, 2}), "1") // default ветка (%v)
}
