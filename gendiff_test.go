package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper: нормализация окончаний строк
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// Helper: загрузка и нормализация expected-файла
func loadExpectedFixture(t *testing.T, fixturePath string) string {
	t.Helper()
	content, err := os.ReadFile(fixturePath)
	require.NoError(t, err, "Failed to read fixture: %s", fixturePath)
	return normalizeLineEndings(string(content))
}

// ============================================================================
// Тесты на вложенные структуры (покрывают и плоские случаи)
// ============================================================================

func TestGenDiffNestedJSON(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	result, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.json"),
		filepath.Join(fixtureDir, "file2_nested.json"),
		"stylish",
	)
	require.NoError(t, err)

	expected := loadExpectedFixture(t, filepath.Join(fixtureDir, "expected_stylish_nested.txt"))
	assert.Equal(t, expected, normalizeLineEndings(result), "Nested JSON diff should match expected")
}

func TestGenDiffNestedYAML(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	result, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.yml"),
		filepath.Join(fixtureDir, "file2_nested.yml"),
		"stylish",
	)
	require.NoError(t, err)

	expected := loadExpectedFixture(t, filepath.Join(fixtureDir, "expected_stylish_nested.txt"))
	assert.Equal(t, expected, normalizeLineEndings(result), "Nested YAML diff should match expected")
}

func TestGenDiffIdenticalNested(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	file := filepath.Join(fixtureDir, "file1_nested.json")

	result, err := GenDiff(file, file, "stylish")
	require.NoError(t, err)

	assert.NotContains(t, result, "  + ")
	assert.NotContains(t, result, "  - ")
	assert.Contains(t, result, "setting1: Value 1")
}

// ============================================================================
// Тесты на обработку ошибок (граничные случаи — оставляем!)
// ============================================================================

func TestGenDiffInvalidPath(t *testing.T) {
	_, err := GenDiff("nonexistent.json", "testdata/fixture/file1_nested.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestGenDiffUnknownFormat(t *testing.T) {
	dir := t.TempDir()
	txtFile := filepath.Join(dir, "config.txt")
	require.NoError(t, os.WriteFile(txtFile, []byte("key: value"), 0644))

	_, err := GenDiff(txtFile, "testdata/fixture/file1_nested.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format")
}

func TestGenDiffInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	badJSON := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(badJSON, []byte(`{"invalid": json}`), 0644))

	_, err := GenDiff(badJSON, "testdata/fixture/file1_nested.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

// ============================================================================
// Тесты на форматирование значений (утилитарные — оставляем)
// ============================================================================

func TestFormatValueFloatDecimal(t *testing.T) {
	dir := t.TempDir()
	file1 := filepath.Join(dir, "f1.json")
	file2 := filepath.Join(dir, "f2.json")

	require.NoError(t, os.WriteFile(file1, []byte(`{"rate": 3.14}`), 0644))
	require.NoError(t, os.WriteFile(file2, []byte(`{"rate": 2.71}`), 0644))

	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)
	assert.Contains(t, result, "3.14")
	assert.Contains(t, result, "2.71")
}

func TestFormatValueSlice(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.json")
	f2 := filepath.Join(dir, "b.json")

	require.NoError(t, os.WriteFile(f1, []byte(`{"items": [1, 2, 3]}`), 0644))
	require.NoError(t, os.WriteFile(f2, []byte(`{"items": [4, 5]}`), 0644))

	result, err := GenDiff(f1, f2, "stylish")
	require.NoError(t, err)
	assert.Contains(t, result, "items:")
}

func TestFormatValueNil(t *testing.T) {
	result := formatValue(nil)
	assert.Equal(t, "null", result)
}
