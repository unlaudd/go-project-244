// Package formatters содержит тесты для plain-форматера.
// Проверяет форматирование путей, обработку разных типов значений
// и маркер [complex value] для вложенных структур.
package formatters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatPlain проверяет вывод plain-форматера на комплексном наборе данных.
// Тестирует: полные пути до вложенных ключей, кавычки для строк, [complex value]
// для объектов, обработку null и changed-состояний с oldValue/newValue.
func TestFormatPlain(t *testing.T) {
	// Фикстура имитирует реальное дерево различий с разной глубиной вложенности
	// и комбинацией состояний: added, removed, unchanged, changed.
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
			// Проверяем смену типа: объект → строка
			{Key: "nest", State: "changed", OldValue: map[string]interface{}{"key": "value"}, NewValue: "str"},
		}},
		{Key: "group2", State: "removed", Value: map[string]interface{}{"abc": 12345.0}},
	}

	result := FormatPlain(nodes)
	lines := strings.Split(result, "\n")

	assert.True(t, len(lines) > 0, "Output should not be empty")

	// Проверяем форматирование путей: точка как разделитель уровней вложенности
	assert.Contains(t, result, "Property 'common.follow' was added with value: false")
	assert.Contains(t, result, "Property 'common.setting2' was removed")

	// Проверяем обработку null и формат changed с двумя значениями
	assert.Contains(t, result, "Property 'common.setting3' was updated. From true to null")

	// Проверяем маркер [complex value] для вложенных объектов
	assert.Contains(t, result, "Property 'common.setting5' was added with value: [complex value]")

	// Проверяем глубокую вложенность: путь должен включать все уровни до ключа
	assert.Contains(t, result, "Property 'common.setting6.doge.wow' was updated. From '' to 'so much'")

	// Проверяем смену типа значения: объект в строку
	assert.Contains(t, result, "Property 'group1.nest' was updated. From [complex value] to 'str'")
	assert.Contains(t, result, "Property 'group2' was removed")
}

// TestFormatPlainValue проверяет форматирование примитивных значений.
// Убеждаемся, что строки в кавычках, числа без лишних знаков,
// а сложные типы возвращают [complex value].
func TestFormatPlainValue(t *testing.T) {
	assert.Equal(t, "null", formatPlainValue(nil))
	assert.Equal(t, "'test'", formatPlainValue("test"))
	assert.Equal(t, "true", formatPlainValue(true))
	assert.Equal(t, "100", formatPlainValue(100.0))
	assert.Equal(t, "99.9", formatPlainValue(99.9))

	// map и slice считаются сложными значениями и заменяются маркером
	assert.Equal(t, "[complex value]", formatPlainValue(map[string]interface{}{"a": 1}))
	assert.Equal(t, "[complex value]", formatPlainValue([]interface{}{1, 2}))
}
