package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ParseFile читает файл и возвращает данные в виде map[string]interface{}
// Формат определяется по расширению файла
func ParseFile(filePath string) (map[string]interface{}, error) {
	// Приводим путь к абсолютному и разрешаем символы . и ..
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path %s: %w", filePath, err)
	}

	// Читаем содержимое файла
	content, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", absPath, err)
	}

	// Определяем формат по расширению
	ext := filepath.Ext(absPath)
	switch ext {
	case ".json":
		return parseJSON(content)
	case ".yaml", ".yml":
		// Пока возвращаем заглушку — реализуем на следующем шаге
		return nil, fmt.Errorf("format %s is not supported yet", ext)
	default:
		return nil, fmt.Errorf("unknown format: %s", ext)
	}
}

// parseJSON парсит JSON-данные в map[string]interface{}
func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}
