// Package formatters предоставляет реализацию stylish-форматера для вывода
// различий в виде иерархической структуры с отступами и маркерами состояний.
package formatters

import (
	"fmt"
	"sort"
	"strings"
)

// FormatStylish форматирует дерево различий в стильный вывод с отступами и маркерами.
// Вывод начинается с '{' и заканчивается '}', ключи сортируются по алфавиту.
func FormatStylish(nodes []DiffNode) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	// Начинаем с глубины 1, так как корневые скобки уже добавлены
	sb.WriteString(formatStylishNodes(nodes, 1))
	sb.WriteString("}")
	return sb.String()
}

// formatStylishNodes рекурсивно обходит узлы дерева и форматирует их.
// depth: уровень вложенности (1 = внутри корневых {}, 2 = первый уровень ключей и т.д.)
// Формула отступов: ключи — depth*4 пробелов, маркеры (+/-) — на 2 пробела левее.
func formatStylishNodes(nodes []DiffNode, depth int) string {
	var sb strings.Builder
	keyIndent := strings.Repeat(" ", depth*4)
	markerIndent := strings.Repeat(" ", depth*4-2)

	for _, node := range nodes {
		switch node.State {
		case "unchanged":
			// Неизменённые ключи выводятся с отступом, без маркера
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", keyIndent, node.Key, formatStylishValue(node.Value)))

		case "added":
			if isMap(node.Value) {
				// Добавлен объект: открываем блок, рекурсивно форматируем содержимое без маркеров
				sb.WriteString(fmt.Sprintf("%s+ %s: {\n", markerIndent, node.Key))
				sb.WriteString(formatStylishMapContent(node.Value.(map[string]interface{}), depth+1))
				sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
			} else {
				sb.WriteString(fmt.Sprintf("%s+ %s: %s\n", markerIndent, node.Key, formatStylishValue(node.Value)))
			}

		case "removed":
			if isMap(node.Value) {
				sb.WriteString(fmt.Sprintf("%s- %s: {\n", markerIndent, node.Key))
				sb.WriteString(formatStylishMapContent(node.Value.(map[string]interface{}), depth+1))
				sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
			} else {
				sb.WriteString(fmt.Sprintf("%s- %s: %s\n", markerIndent, node.Key, formatStylishValue(node.Value)))
			}

		case "changed":
			if node.Children != nil {
				// Оба значения — map: рекурсивно сравниваем детей, ключ без маркера
				sb.WriteString(fmt.Sprintf("%s%s: {\n", keyIndent, node.Key))
				sb.WriteString(formatStylishNodes(node.Children, depth+1))
				sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
			} else {
				// Изменено примитивное значение: показываем старое и новое с маркерами
				if isMap(node.OldValue) {
					sb.WriteString(fmt.Sprintf("%s- %s: {\n", markerIndent, node.Key))
					sb.WriteString(formatStylishMapContent(node.OldValue.(map[string]interface{}), depth+1))
					sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
				} else {
					sb.WriteString(fmt.Sprintf("%s- %s: %s\n", markerIndent, node.Key, formatStylishValue(node.OldValue)))
				}

				if isMap(node.NewValue) {
					sb.WriteString(fmt.Sprintf("%s+ %s: {\n", markerIndent, node.Key))
					sb.WriteString(formatStylishMapContent(node.NewValue.(map[string]interface{}), depth+1))
					sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
				} else {
					sb.WriteString(fmt.Sprintf("%s+ %s: %s\n", markerIndent, node.Key, formatStylishValue(node.NewValue)))
				}
			}
		}
	}
	return sb.String()
}

// formatStylishMapContent форматирует содержимое добавленного или удалённого объекта.
// Выводит ключи с отступом, но без маркеров + / -, так как маркер уже применён к родительскому ключу.
func formatStylishMapContent(m map[string]interface{}, depth int) string {
	var sb strings.Builder
	keyIndent := strings.Repeat(" ", depth*4)

	for _, key := range getSortedKeys(m) {
		val := m[key]
		if isMap(val) {
			sb.WriteString(fmt.Sprintf("%s%s: {\n", keyIndent, key))
			sb.WriteString(formatStylishMapContent(val.(map[string]interface{}), depth+1))
			sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", keyIndent, key, formatStylishValue(val)))
		}
	}
	return sb.String()
}

// formatStylishValue приводит примитивное значение к строке для stylish-вывода.
// Целые числа выводятся без десятичной точки для читаемости (42 вместо 42.0).
func formatStylishValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "null"
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isMap проверяет, является ли значение map[string]interface{}.
func isMap(val interface{}) bool {
	_, ok := val.(map[string]interface{})
	return ok
}

// getSortedKeys собирает и сортирует ключи карты в алфавитном порядке.
// Необходимо, так как порядок ключей в map в Go не гарантирован.
func getSortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
