// Package formatters содержит тесты для stylish-форматера.
// Проверяет отступы, маркеры + / - / пробел, и форматирование значений.
package formatters

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatStylish проверяет базовую структуру stylish-вывода:
// корректные отступы, маркеры состояний, обработку вложенных объектов
// и форматирование разных типов значений (включая null).
func TestFormatStylish(t *testing.T) {
	// Фикстура покрывает основные состояния: added, removed, unchanged, changed,
	// а также вложенный объект (group2) для проверки рекурсивного форматирования.
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

	// Проверяем структуру: вывод должен быть обернут в фигурные скобки
	assert.True(t, strings.HasPrefix(result, "{\n"), "Should start with {")
	assert.True(t, strings.HasSuffix(result, "}"), "Should end with }")

	// Проверяем маркеры состояний и отступы
	// Формула: ключи — 4*глубина пробелов, маркеры — на 2 пробела левее
	assert.Contains(t, result, "+ follow: false")
	assert.Contains(t, result, "  setting1: Value 1")
	assert.Contains(t, result, "- setting2: 200")

	// Проверяем обработку changed-состояния: оба значения выводятся рядом
	assert.Contains(t, result, "- setting3: true")
	assert.Contains(t, result, "+ setting3: null")

	// Проверяем форматирование удалённого вложенного объекта
	assert.Contains(t, result, "- group2: {")
}

// TestFormatStylishValue проверяет форматирование примитивных значений.
// Тест покрывает все ветки switch в formatStylishValue для предотвращения регрессий.
func TestFormatStylishValue(t *testing.T) {
	assert.Equal(t, "null", formatStylishValue(nil))
	assert.Equal(t, "hello", formatStylishValue("hello"))
	assert.Equal(t, "true", formatStylishValue(true))
	assert.Equal(t, "false", formatStylishValue(false))

	// Целые числа выводим без десятичной точки для читаемости
	assert.Equal(t, "42", formatStylishValue(42.0))
	assert.Equal(t, "3.14", formatStylishValue(3.14))

	// Неизвестные типы форматируются через %v — проверяем, что не паникует
	assert.Contains(t, formatStylishValue([]int{1, 2}), "1")
}
