// Package code предоставляет основную библиотеку для сравнения
// конфигурационных файлов. Экспортирует функцию GenDiff, которая
// принимает пути к файлам и формат вывода, возвращая разницу как строку.
package code

import (
	"code/formatters"
	"code/internal/parser"
	"fmt"
	"reflect"
	"sort"
)

// GenDiff сравнивает два конфигурационных файла и возвращает разницу.
// Если format пустой, используется формат "stylish" по умолчанию.
// Ошибки парсинга оборачиваются с контекстом для удобной отладки.
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parser.ParseFile(filepath1)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath1, err)
	}

	data2, err := parser.ParseFile(filepath2)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath2, err)
	}

	// Строим промежуточное представление (AST) различий.
	// Это разделяет логику сравнения и форматирования: новые форматы
	// можно добавлять, не меняя алгоритм построения дерева.
	diffTree := buildDiffTree(data1, data2)

	return formatters.Format(diffTree, format)
}

// buildDiffTree рекурсивно строит дерево различий между двумя картами.
// Возвращает слайс узлов, где каждый узел описывает состояние ключа:
// добавлен, удалён, неизменён или изменён. Для вложенных объектов
// рекурсивно строит детей, для примитивов — сохраняет значения.
func buildDiffTree(data1, data2 map[string]interface{}) []formatters.DiffNode {
	keys := getSortedKeys(data1, data2)
	var diff []formatters.DiffNode

	for _, key := range keys {
		val1, ok1 := data1[key]
		val2, ok2 := data2[key]

		switch {
		case ok1 && !ok2:
			// Ключ присутствует только в первом файле → помечаем как удалённый
			diff = append(diff, formatters.DiffNode{Key: key, State: "removed", Value: val1})
		case !ok1 && ok2:
			// Ключ присутствует только во втором файле → помечаем как добавленный
			diff = append(diff, formatters.DiffNode{Key: key, State: "added", Value: val2})
		case ok1 && ok2:
			m1, isMap1 := val1.(map[string]interface{})
			m2, isMap2 := val2.(map[string]interface{})

			if isMap1 && isMap2 {
				// Оба значения — объекты: рекурсивно сравниваем их содержимое.
				// Это ключевая особенность: вложенные структуры раскрываются,
				// а не сравниваются как целое, что даёт детальный дифф.
				diff = append(diff, formatters.DiffNode{
					Key:      key,
					State:    "changed",
					OldValue: val1,
					NewValue: val2,
					Children: buildDiffTree(m1, m2),
				})
			} else if reflect.DeepEqual(val1, val2) {
				// Значения совпадают (включая сравнение слайсов и вложенных структур)
				diff = append(diff, formatters.DiffNode{Key: key, State: "unchanged", Value: val1})
			} else {
				// Значения различаются: примитивы или разные типы (map vs string и т.п.)
				// Сохраняем оба значения для отображения в форматере.
				diff = append(diff, formatters.DiffNode{
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

// getSortedKeys собирает и сортирует уникальные ключи из двух карт.
// Сортировка необходима для детерминированного вывода: порядок ключей
// в map в Go не гарантирован, а для диффа важен стабильный алфавитный порядок.
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
