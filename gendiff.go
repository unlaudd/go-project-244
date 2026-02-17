package code

import (
	"fmt"
	"sort"
	"strings"

	"code/internal/parser"
)

// GenDiff сравнивает два конфигурационных файла и возвращает разницу в виде строки
func GenDiff(filepath1, filepath2, format string) (string, error) {
	data1, err := parser.ParseFile(filepath1)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath1, err)
	}

	data2, err := parser.ParseFile(filepath2)
	if err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", filepath2, err)
	}

	return buildDiff(data1, data2), nil
}

func buildDiff(data1, data2 map[string]interface{}) string {
	var result strings.Builder
	result.WriteString("{\n")

	// Собираем все уникальные ключи
	keySet := make(map[string]bool)
	for k := range data1 {
		keySet[k] = true
	}
	for k := range data2 {
		keySet[k] = true
	}

	// Сортируем ключи в алфавитном порядке
	keys := make([]string, 0, len(keySet))
	for k := range keySet {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Формируем вывод диффа
	for _, key := range keys {
		val1, exists1 := data1[key]
		val2, exists2 := data2[key]

		if exists1 && !exists2 {
			// Ключ только в первом файле
			result.WriteString(fmt.Sprintf("  - %s: %s\n", key, formatValue(val1)))
		} else if !exists1 && exists2 {
			// Ключ только во втором файле
			result.WriteString(fmt.Sprintf("  + %s: %s\n", key, formatValue(val2)))
		} else if val1 == val2 {
			// Значения совпадают
			result.WriteString(fmt.Sprintf("    %s: %s\n", key, formatValue(val1)))
		} else {
			// Значения различаются: сначала строка из первого файла, затем из второго
			result.WriteString(fmt.Sprintf("  - %s: %s\n", key, formatValue(val1)))
			result.WriteString(fmt.Sprintf("  + %s: %s\n", key, formatValue(val2)))
		}
	}

	result.WriteString("}")
	return result.String()
}

// formatValue форматирует значение для вывода (без кавычек для строк, lowercase для bool)
func formatValue(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case float64:
		// JSON числа приходят как float64
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
