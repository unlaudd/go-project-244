package formatters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Минимальный набор узлов для тестирования фабрики
var testNodes = []DiffNode{
	{Key: "key", State: "added", Value: "val"},
}

func TestFormat(t *testing.T) {
	t.Run("empty string defaults to stylish", func(t *testing.T) {
		res, err := Format(testNodes, "")
		require.NoError(t, err)
		// У stylish вывод начинается с "{\n", у plain — с "Property"
		assert.True(t, strings.HasPrefix(res, "{\n"), "Default format should be stylish")
	})

	t.Run("explicit stylish format", func(t *testing.T) {
		res, err := Format(testNodes, "stylish")
		require.NoError(t, err)
		assert.True(t, strings.HasPrefix(res, "{\n"))
	})

	t.Run("plain format", func(t *testing.T) {
		res, err := Format(testNodes, "plain")
		require.NoError(t, err)
		assert.Contains(t, res, "Property 'key' was added with value: 'val'")
	})

	t.Run("unknown format returns error", func(t *testing.T) {
		_, err := Format(testNodes, "json")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown format: json")
	})
}
