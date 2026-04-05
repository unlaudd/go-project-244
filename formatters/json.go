package formatters

import (
	"encoding/json"
)

// jsonNode — внутреннее представление для JSON-сериализации.
// Поля с тегами `omitempty` не попадут в вывод, если они пустые.
type jsonNode struct {
	Status   string              `json:"status"` // "added", "removed", "unchanged", "changed"
	Value    interface{}         `json:"value,omitempty"`
	OldValue interface{}         `json:"oldValue,omitempty"`
	NewValue interface{}         `json:"newValue,omitempty"`
	Children map[string]jsonNode `json:"children,omitempty"`
}

// FormatJSON форматирует дерево различий в структурированный JSON.
func FormatJSON(nodes []DiffNode) string {
	// Корневой объект: ключи верхнего уровня - их статусы
	root := make(map[string]jsonNode)

	for _, node := range nodes {
		root[node.Key] = buildJSONNode(node)
	}

	// MarshalIndent для красивого вывода с отступами
	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		// Здесь возвращаем пустой объект, так как ошибка маловероятна
		return "{}"
	}

	return string(output)
}

// buildJSONNode рекурсивно преобразует DiffNode в jsonNode
func buildJSONNode(node DiffNode) jsonNode {
	jn := jsonNode{
		Status: node.State,
	}

	switch node.State {
	case "added":
		if node.Children != nil {
			// Добавлен объект - рекурсивно строим детей
			jn.Children = buildChildrenMap(node.Children)
		} else {
			jn.Value = node.Value
		}
	case "removed":
		if node.Children != nil {
			jn.Children = buildChildrenMap(node.Children)
		} else {
			jn.Value = node.Value
		}
	case "unchanged":
		if node.Children != nil {
			jn.Children = buildChildrenMap(node.Children)
		} else {
			jn.Value = node.Value
		}
	case "changed":
		if node.Children != nil {
			// Изменён объект - рекурсивно сравниваем детей
			jn.Children = buildChildrenMap(node.Children)
		} else {
			// Изменено примитивное значение - показываем старое и новое
			jn.OldValue = node.OldValue
			jn.NewValue = node.NewValue
		}
	}

	return jn
}

// buildChildrenMap преобразует слайс детей в map[string]jsonNode для JSON
func buildChildrenMap(children []DiffNode) map[string]jsonNode {
	result := make(map[string]jsonNode, len(children))
	for _, child := range children {
		result[child.Key] = buildJSONNode(child)
	}
	return result
}
