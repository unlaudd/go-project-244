package formatters

import (
	"fmt"
	"strings"
)

// FormatPlain форматирует дерево различий в плоский человеко-читаемый вывод.
func FormatPlain(nodes []DiffNode) string {
	var lines []string
	renderPlain(nodes, []string{}, &lines)
	return strings.Join(lines, "\n")
}

// renderPlain рекурсивно обходит дерево и собирает строки вывода.
func renderPlain(nodes []DiffNode, path []string, lines *[]string) {
	for _, node := range nodes {
		// Формируем полный путь до текущего ключа
		currentPath := append(path, node.Key)
		pathStr := strings.Join(currentPath, ".")

		switch node.State {
		case "added":
			val := formatPlainValue(node.Value)
			*lines = append(*lines, fmt.Sprintf("Property '%s' was added with value: %s", pathStr, val))
		case "removed":
			*lines = append(*lines, fmt.Sprintf("Property '%s' was removed", pathStr))
		case "unchanged":
			// Пропускаем неизменённые ключи
			continue
		case "changed":
			if node.Children != nil {
				// Вложенный объект изменился → рекурсивно обходим детей
				renderPlain(node.Children, currentPath, lines)
			} else {
				oldVal := formatPlainValue(node.OldValue)
				newVal := formatPlainValue(node.NewValue)
				*lines = append(*lines, fmt.Sprintf("Property '%s' was updated. From %s to %s", pathStr, oldVal, newVal))
			}
		}
	}
}

// formatPlainValue форматирует отдельное значение для plain-вывода.
func formatPlainValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case float64:
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case map[string]interface{}:
		return "[complex value]"
	default:
		// Слайсы, указатели и другие типы тоже считаем сложными
		return "[complex value]"
	}
}
