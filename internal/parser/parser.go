package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ParseFile читает файл и возвращает данные в виде map[string]interface{}
// Формат определяется по расширению файла
func ParseFile(filePath string) (map[string]interface{}, error) {
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

// parseJSON парсит JSON-данные в map[string]interface{}
func parseJSON(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}
	return result, nil
}

// parseYAML парсит YAML-данные в map[string]interface{}
func parseYAML(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return result, nil
}
