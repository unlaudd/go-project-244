// Package formatters предоставляет реализацию JSON-форматера для вывода
// структурированных различий между конфигурационными файлами.
package formatters

import (
	"encoding/json"
)

// jsonNode — внутреннее представление узла для JSON-сериализации.
// Использует теги omitempty, чтобы исключать пустые поля из вывода.
type jsonNode struct {
	Status   string              `json:"status"`
	Value    interface{}         `json:"value,omitempty"`
	OldValue interface{}         `json:"oldValue,omitempty"`
	NewValue interface{}         `json:"newValue,omitempty"`
	Children map[string]jsonNode `json:"children,omitempty"`
}

// FormatJSON форматирует дерево различий в структурированный JSON.
// Вывод использует отступы для читаемости (MarshalIndent).
func FormatJSON(nodes []DiffNode) string {
	root := make(map[string]jsonNode)
	for _, node := range nodes {
		root[node.Key] = buildJSONNode(node)
	}

	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		// Ошибка сериализации маловероятна при корректной структуре,
		// но на всякий случай возвращаем валидный пустой объект.
		return "{}"
	}

	return string(output)
}

// buildJSONNode преобразует DiffNode в jsonNode для сериализации.
// Рекурсивно обрабатывает вложенные объекты через поле Children.
func buildJSONNode(node DiffNode) jsonNode {
	jn := jsonNode{Status: node.State}

	switch node.State {
	case "added", "removed", "unchanged":
		if node.Children != nil {
			jn.Children = buildChildrenMap(node.Children)
		} else {
			jn.Value = node.Value
		}
	case "changed":
		if node.Children != nil {
			jn.Children = buildChildrenMap(node.Children)
		} else {
			// Для изменённых примитивов сохраняем оба значения,
			// чтобы потребитель мог увидеть, что именно изменилось.
			jn.OldValue = node.OldValue
			jn.NewValue = node.NewValue
		}
	}

	return jn
}

// buildChildrenMap преобразует слайс детей в map[string]jsonNode.
// JSON-формат требует объектной структуры для вложенных ключей,
// поэтому используем map вместо слайса.
func buildChildrenMap(children []DiffNode) map[string]jsonNode {
	result := make(map[string]jsonNode, len(children))
	for _, child := range children {
		result[child.Key] = buildJSONNode(child)
	}
	return result
}
