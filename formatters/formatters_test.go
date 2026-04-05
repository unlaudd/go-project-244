// Package formatters содержит тесты для фабрики форматеров.
// Тестирует маршрутизацию форматов и базовую корректность вывода.
package formatters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testNodes — минимальная фикстура для тестирования фабрики.
// Содержит один узел, достаточный для проверки всех веток switch по форматам.
var testNodes = []DiffNode{
	{Key: "key", State: "added", Value: "val"},
}

// TestFormat проверяет фабрику форматеров: выбор по умолчанию,
// явные форматы и обработку неизвестных значений.
func TestFormat(t *testing.T) {
	t.Run("empty string defaults to stylish", func(t *testing.T) {
		// Пустая строка должна подставлять "stylish" — дефолт для библиотеки
		res, err := Format(testNodes, "")
		require.NoError(t, err)
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

	t.Run("json format", func(t *testing.T) {
		// Проверяем, что JSON-вывод содержит ожидаемые поля с форматированием MarshalIndent
		res, err := Format(testNodes, "json")
		require.NoError(t, err)
		assert.Contains(t, res, `"status": "added"`)
		assert.Contains(t, res, `"value": "val"`)
	})

	t.Run("unknown format returns error", func(t *testing.T) {
		// Проверяем обработку неподдерживаемого формата.
		// Используем "xml", так как он заведомо отсутствует в switch.
		_, err := Format(testNodes, "xml")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unknown format: xml")
	})
}
