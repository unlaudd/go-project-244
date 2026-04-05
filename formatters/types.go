// Package formatters предоставляет типы и интерфейсы для форматеров вывода.
// Центральный тип — DiffNode, используемый как промежуточное представление (AST)
// между логикой сравнения и конкретными реализациями форматирования.
package formatters

// DiffNode представляет узел различий в конфигурации.
// Используется как промежуточное представление (AST) между
// логикой сравнения и форматерами вывода.
type DiffNode struct {
	Key      string      // Имя ключа в конфигурации
	State    string      // Состояние: "added", "removed", "unchanged", "changed"
	Value    interface{} // Значение для added/removed/unchanged
	OldValue interface{} // Предыдущее значение (только для changed)
	NewValue interface{} // Новое значение (только для changed)
	Children []DiffNode  // Вложенные узлы (заполняется, если значение — map)
}
