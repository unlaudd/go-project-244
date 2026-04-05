package formatters

// DiffNode представляет узел различий в конфигурации.
// Используется как промежуточное представление (AST) между
// логикой сравнения и форматерами вывода.
type DiffNode struct {
	Key      string      // Имя ключа
	State    string      // "added", "removed", "unchanged", "changed"
	Value    interface{} // Текущее значение (для added/removed/unchanged)
	OldValue interface{} // Старое значение (для changed)
	NewValue interface{} // Новое значение (для changed)
	Children []DiffNode  // Вложенные ключи (только если значение — map)
}
