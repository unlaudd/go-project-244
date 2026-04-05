// Package parser содержит тесты для парсера конфигурационных файлов.
// Проверяет поддержку форматов JSON/YAML, обработку ошибок и работу с путями.
package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseFileJSON проверяет парсинг валидного JSON-файла.
// Убеждаемся, что строки и числа корректно преобразуются в map[string]interface{}.
func TestParseFileJSON(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.json")
	require.NoError(t, os.WriteFile(file, []byte(`{"key":"val","num":42}`), 0644))

	res, err := ParseFile(file)
	require.NoError(t, err)
	assert.Equal(t, "val", res["key"])
	// JSON-парсер в Go всегда возвращает числа как float64
	assert.Equal(t, float64(42), res["num"])
}

// TestParseFileYAML проверяет парсинг валидного YAML-файла.
// Отличается от JSON тем, что числа могут возвращаться как int, а не float64.
func TestParseFileYAML(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.yaml")
	require.NoError(t, os.WriteFile(file, []byte("key: val\nnum: 42"), 0644))

	res, err := ParseFile(file)
	require.NoError(t, err)
	assert.Equal(t, "val", res["key"])
	// YAML-парсер сохраняет тип числа, поэтому 42 приходит как int
	assert.Equal(t, 42, res["num"])
}

// TestParseFileUnknownExt проверяет обработку файла с неподдерживаемым расширением.
// Ошибка должна содержать понятное сообщение о неизвестном формате.
func TestParseFileUnknownExt(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")
	require.NoError(t, os.WriteFile(file, []byte("data"), 0644))

	_, err := ParseFile(file)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format")
}

// TestParseFileInvalidJSON проверяет обработку синтаксически невалидного JSON.
// Ошибка должна указывать на проблему парсинга, а не чтения файла.
func TestParseFileInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(file, []byte(`{invalid}`), 0644))

	_, err := ParseFile(file)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

// TestParseFileNotFound проверяет обработку несуществующего файла.
// Ошибка должна указывать на проблему чтения, а не парсинга.
func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("nonexistent.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}
