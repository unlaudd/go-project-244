// Package formatters предоставляет реализацию plain-форматера для вывода
// различий в виде плоского списка изменений с полными путями.
package formatters

import (
	"fmt"
	"strings"
)

// FormatPlain форматирует дерево различий в плоский человеко-читаемый вывод.
// Каждый изменённый ключ выводится как отдельная строка с полным путём.
func FormatPlain(nodes []DiffNode) string {
	var lines []string
	renderPlain(nodes, []string{}, &lines)
	return strings.Join(lines, "\n")
}

// renderPlain рекурсивно обходит дерево и собирает строки вывода.
// path накапливает цепочку ключей от корня до текущего узла.
// lines — аккумулятор результатов, передаётся по указателю для избежания копирования.
func renderPlain(nodes []DiffNode, path []string, lines *[]string) {
	for _, node := range nodes {
		// Формируем полный путь: соединяем ключи через точку
		// Создаём новый слайс, чтобы избежать проблем с переиспользованием буфера при рекурсии
		currentPath := append(append([]string(nil), path...), node.Key)
		pathStr := strings.Join(currentPath, ".")

		switch node.State {
		case "added":
			val := formatPlainValue(node.Value)
			*lines = append(*lines, fmt.Sprintf("Property '%s' was added with value: %s", pathStr, val))
		case "removed":
			*lines = append(*lines, fmt.Sprintf("Property '%s' was removed", pathStr))
		case "unchanged":
			// Пропускаем неизменённые ключи — они не интересны в диффе
			continue
		case "changed":
			if node.Children != nil {
				// Вложенный объект изменился — рекурсивно обрабатываем детей
				// с тем же путём, чтобы сохранить полную иерархию
				renderPlain(node.Children, currentPath, lines)
			} else {
				// Изменено примитивное значение — показываем старое и новое
				oldVal := formatPlainValue(node.OldValue)
				newVal := formatPlainValue(node.NewValue)
				*lines = append(*lines, fmt.Sprintf("Property '%s' was updated. From %s to %s", pathStr, oldVal, newVal))
			}
		}
	}
}

// formatPlainValue форматирует отдельное значение для plain-вывода.
// Строки оборачиваются в одинарные кавычки, сложные типы заменяются маркером.
// Целые числа выводятся без десятичной точки для читаемости.
func formatPlainValue(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return "null"
	case string:
		return fmt.Sprintf("'%s'", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case int:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float64:
		// Целые числа выводим без десятичной точки для читаемости
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	case map[string]interface{}:
		return "[complex value]"
	default:
		// Слайсы, указатели и другие типы тоже считаем сложными
		// и заменяем маркером, чтобы не перегружать вывод
		return "[complex value]"
	}
}
