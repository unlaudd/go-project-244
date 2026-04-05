package formatters

import (
	"fmt"
	"sort"
	"strings"
)

// FormatStylish форматирует дерево различий в стильный вывод с отступами и маркерами.
func FormatStylish(nodes []DiffNode) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(formatStylishNodes(nodes, 1))
	sb.WriteString("}")
	return sb.String()
}

// formatStylishNodes рекурсивно обходит узлы дерева
// depth: 1 = внутри корневых {}, 2 = первый уровень вложенности и т.д.
func formatStylishNodes(nodes []DiffNode, depth int) string {
	var sb strings.Builder
	keyIndent := strings.Repeat(" ", depth*4)
	markerIndent := strings.Repeat(" ", depth*4-2)

	for _, node := range nodes {
		switch node.State {
		case "unchanged":
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", keyIndent, node.Key, formatStylishValue(node.Value)))

		case "added":
			if isMap(node.Value) {
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
				// Оба значения map → рекурсивный диф
				sb.WriteString(fmt.Sprintf("%s%s: {\n", keyIndent, node.Key))
				sb.WriteString(formatStylishNodes(node.Children, depth+1))
				sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
			} else {
				// Значения изменились, но не map → показываем старое и новое
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

// formatStylishMapContent форматирует содержимое добавленного/удалённого объекта
// БЕЗ маркеров + / - (только с отступом)
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

// formatStylishValue приводит примитивное значение к строке для stylish-вывода
func formatStylishValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "null"
	case string:
		return v
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isMap проверяет, является ли значение map
func isMap(val interface{}) bool {
	_, ok := val.(map[string]interface{})
	return ok
}

// getSortedKeys собирает и сортирует ключи карты
func getSortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
