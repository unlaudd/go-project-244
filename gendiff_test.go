package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func loadExpectedFixture(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read fixture: %s", path)
	return normalizeLineEndings(string(content))
}

// Тест stylish через GenDiff (интеграционный)
func TestGenDiffNestedStylish(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	result, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.json"),
		filepath.Join(fixtureDir, "file2_nested.json"),
		"stylish",
	)
	require.NoError(t, err)

	expected := loadExpectedFixture(t, filepath.Join(fixtureDir, "expected_stylish_nested.txt"))
	assert.Equal(t, expected, normalizeLineEndings(result))
}

// Тест plain через GenDiff
func TestGenDiffNestedPlain(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	result, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.json"),
		filepath.Join(fixtureDir, "file2_nested.json"),
		"plain",
	)
	require.NoError(t, err)

	assert.Contains(t, result, "Property 'common.follow' was added with value: false")
	assert.Contains(t, result, "Property 'common.setting5' was added with value: [complex value]")
	assert.Contains(t, result, "Property 'common.setting6.doge.wow' was updated. From '' to 'so much'")
	assert.Contains(t, result, "Property 'group2' was removed")
	assert.NotContains(t, result, "unchanged") // plain пропускает неизменённые
}

// Тест формата по умолчанию (пустая строка → stylish)
func TestGenDiffDefaultFormat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	result, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.json"),
		filepath.Join(fixtureDir, "file2_nested.json"),
		"", // пустой формат
	)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result, "{\n"), "Default format should be stylish")
}

// Ошибки
func TestGenDiffUnknownFormat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")

	// Передаём валидные файлы, но неизвестный формат
	_, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.json"),
		filepath.Join(fixtureDir, "file2_nested.json"),
		"markdown", // ← неизвестный формат
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format: markdown")
}

func TestGenDiffInvalidPath(t *testing.T) {
	_, err := GenDiff("nonexistent.json", "testdata/fixture/file1_nested.json", "stylish")
	assert.Error(t, err)
	// Ошибка может быть "failed to read" или "failed to parse" в зависимости от реализации парсера
	assert.True(t, strings.Contains(err.Error(), "failed to read") || strings.Contains(err.Error(), "failed to parse"))
}

func TestGenDiffInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	badJSON := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(badJSON, []byte(`{invalid}`), 0644))

	_, err := GenDiff(badJSON, "testdata/fixture/file1_nested.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}
