// Package parser предоставляет функции для чтения и парсинга
// конфигурационных файлов в форматах JSON и YAML.
// Возвращает данные в виде map[string]interface{} для дальнейшей обработки.
package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ParseFile читает файл и возвращает данные в виде map[string]interface{}.
// Формат определяется по расширению файла (.json, .yaml, .yml).
// Путь к файлу разрешается в абсолютный для корректной работы из любой рабочей директории.
func ParseFile(filePath string) (map[string]interface{}, error) {
	// Преобразуем путь в абсолютный, чтобы избежать проблем с относительными путями
	// при запуске утилиты из разных директорий
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path %s: %w", filePath, err)
	}

	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", absPath, err)
	}

	ext := filepath.Ext(absPath)
	switch ext {
	case ".json":
		return parseJSON(content)
	case ".yaml", ".yml":
		return parseYAML(content)
	default:
		return nil, fmt.Errorf("unknown format: %s", ext)
	}
}

// parseJSON парсит JSON-данные в map[string]interface{}.
// Использует стандартный encoding/json, который возвращает числа как float64.
func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		// Оборачиваем ошибку с %w, чтобы вызывающий код мог проверить тип ошибки через errors.Is/As
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

// parseYAML парсит YAML-данные в map[string]interface{}.
// Использует gopkg.in/yaml.v3, который сохраняет нативные типы (int, bool и т.д.).
func parseYAML(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return result, nil
}
