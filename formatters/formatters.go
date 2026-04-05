package formatters

import "fmt"

// Format выбирает и применяет нужный форматер к дереву различий.
// Если format пустой, используется "stylish" по умолчанию.
// Возвращает отформатированную строку или ошибку, если формат неизвестен.
func Format(nodes []DiffNode, format string) (string, error) {
	// Форматер по умолчанию для библиотеки
	if format == "" {
		format = "stylish"
	}

	switch format {
	case "stylish":
		return FormatStylish(nodes), nil
	case "plain":
		return FormatPlain(nodes), nil
	default:
		return "", fmt.Errorf("unknown format: %s", format)
	}
}
