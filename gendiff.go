package code

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"code/internal/parser"
)

// DiffNode представляет узел различий в конфигурации
type DiffNode struct {
	Key      string
	State    string // "added", "removed", "unchanged", "changed"
	Value    interface{}
	OldValue interface{}
	NewValue interface{}
	Children []DiffNode // Только если оба значения - map
}

// GenDiff сравнивает два конфигурационных файла и возвращает разницу
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parser.ParseFile(filepath1)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath1, err)
	}

	data2, err := parser.ParseFile(filepath2)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath2, err)
	}

	diffTree := buildDiffTree(data1, data2)

	// Форматер по умолчанию для библиотеки
	if format == "" {
		format = "stylish"
	}

	switch format {
	case "stylish":
		return formatStylish(diffTree), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}

// buildDiffTree строит дерево различий рекурсивно
func buildDiffTree(data1, data2 map[string]interface{}) []DiffNode {
	keys := getSortedKeys(data1, data2)
	var diff []DiffNode

	for _, key := range keys {
		val1, ok1 := data1[key]
		val2, ok2 := data2[key]

		switch {
		case ok1 && !ok2:
			diff = append(diff, DiffNode{Key: key, State: "removed", Value: val1})
		case !ok1 && ok2:
			diff = append(diff, DiffNode{Key: key, State: "added", Value: val2})
		case ok1 && ok2:
			m1, isMap1 := val1.(map[string]interface{})
			m2, isMap2 := val2.(map[string]interface{})

			if isMap1 && isMap2 {
				diff = append(diff, DiffNode{
					Key:      key,
					State:    "changed",
					OldValue: val1,
					NewValue: val2,
					Children: buildDiffTree(m1, m2),
				})
			} else if reflect.DeepEqual(val1, val2) {
				diff = append(diff, DiffNode{Key: key, State: "unchanged", Value: val1})
			} else {
				diff = append(diff, DiffNode{
					Key:      key,
					State:    "changed",
					OldValue: val1,
					NewValue: val2,
				})
			}
		}
	}
	return diff
}

// formatStylish форматирует дерево различий в строку
func formatStylish(diff []DiffNode) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString(formatNodes(diff, 1))
	sb.WriteString("}")
	return sb.String()
}

// formatNodes обходит узлы дерева и форматирует их
// depth: 1 = внутри корневых {}, 2 = внутри первого уровня вложенности и т.д.
func formatNodes(nodes []DiffNode, depth int) string {
	var sb strings.Builder

	// Формула отступов из задания
	keyIndent := strings.Repeat(" ", depth*4)
	markerIndent := strings.Repeat(" ", depth*4-2)

	for _, node := range nodes {
		switch node.State {
		case "unchanged":
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", keyIndent, node.Key, formatValue(node.Value)))

		case "added":
			if isMap(node.Value) {
				sb.WriteString(fmt.Sprintf("%s+ %s: {\n", markerIndent, node.Key))
				sb.WriteString(formatMapContent(node.Value.(map[string]interface{}), depth+1))
				sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
			} else {
				sb.WriteString(fmt.Sprintf("%s+ %s: %s\n", markerIndent, node.Key, formatValue(node.Value)))
			}

		case "removed":
			if isMap(node.Value) {
				sb.WriteString(fmt.Sprintf("%s- %s: {\n", markerIndent, node.Key))
				sb.WriteString(formatMapContent(node.Value.(map[string]interface{}), depth+1))
				sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
			} else {
				sb.WriteString(fmt.Sprintf("%s- %s: %s\n", markerIndent, node.Key, formatValue(node.Value)))
			}

		case "changed":
			if node.Children != nil {
				// Оба значения map → рекурсивный диф
				sb.WriteString(fmt.Sprintf("%s%s: {\n", keyIndent, node.Key))
				sb.WriteString(formatNodes(node.Children, depth+1))
				sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
			} else {
				// Типы различаются или оба примитивы → показываем старое и новое
				if isMap(node.OldValue) {
					sb.WriteString(fmt.Sprintf("%s- %s: {\n", markerIndent, node.Key))
					sb.WriteString(formatMapContent(node.OldValue.(map[string]interface{}), depth+1))
					sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
				} else {
					sb.WriteString(fmt.Sprintf("%s- %s: %s\n", markerIndent, node.Key, formatValue(node.OldValue)))
				}

				if isMap(node.NewValue) {
					sb.WriteString(fmt.Sprintf("%s+ %s: {\n", markerIndent, node.Key))
					sb.WriteString(formatMapContent(node.NewValue.(map[string]interface{}), depth+1))
					sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
				} else {
					sb.WriteString(fmt.Sprintf("%s+ %s: %s\n", markerIndent, node.Key, formatValue(node.NewValue)))
				}
			}
		}
	}
	return sb.String()
}

// formatMapContent форматирует содержимое добавленного/удалённого объекта БЕЗ маркеров
func formatMapContent(m map[string]interface{}, depth int) string {
	var sb strings.Builder
	keyIndent := strings.Repeat(" ", depth*4)
	keys := getSortedKeysFromMap(m)

	for _, key := range keys {
		val := m[key]
		if isMap(val) {
			sb.WriteString(fmt.Sprintf("%s%s: {\n", keyIndent, key))
			sb.WriteString(formatMapContent(val.(map[string]interface{}), depth+1))
			sb.WriteString(fmt.Sprintf("%s}\n", keyIndent))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s: %s\n", keyIndent, key, formatValue(val)))
		}
	}
	return sb.String()
}

// formatValue приводит значение к строковому представлению
func formatValue(val interface{}) string {
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

// getSortedKeys собирает и сортирует все уникальные ключи из двух карт
func getSortedKeys(m1, m2 map[string]interface{}) []string {
	keySet := make(map[string]bool)
	for k := range m1 {
		keySet[k] = true
	}
	for k := range m2 {
		keySet[k] = true
	}

	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// getSortedKeysFromMap собирает и сортирует ключи одной карты
func getSortedKeysFromMap(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
