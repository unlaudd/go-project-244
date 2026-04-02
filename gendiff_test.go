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

// Helper: сравнение результата GenDiff с expected-файлом
func assertGenDiffEquals(t *testing.T, file1, file2, format, expectedFile, description string) {
	t.Helper()
	result, err := GenDiff(file1, file2, format)
	require.NoError(t, err, "GenDiff should not return error")

	expected := loadExpectedFixture(t, expectedFile)
	resultNormalized := normalizeLineEndings(result)

	assert.Equal(t, expected, resultNormalized, description)
}

func TestGenDiffJsonFlat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	assertGenDiffEquals(
		t,
		filepath.Join(fixtureDir, "file1.json"),
		filepath.Join(fixtureDir, "file2.json"),
		"stylish",
		filepath.Join(fixtureDir, "expected_stylish.txt"),
		"JSON output should match expected diff",
	)
}

func TestGenDiffYamlFlat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	assertGenDiffEquals(
		t,
		filepath.Join(fixtureDir, "file1.yml"),
		filepath.Join(fixtureDir, "file2.yml"),
		"stylish",
		filepath.Join(fixtureDir, "expected_stylish_yaml.txt"),
		"YAML output should match expected diff",
	)
}

func TestGenDiffIdenticalFiles(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	file1 := filepath.Join(fixtureDir, "file1.json")

	result, err := GenDiff(file1, file1, "stylish")
	require.NoError(t, err)

	assert.Contains(t, result, "host: hexlet.io")
	assert.Contains(t, result, "timeout: 50")
	assert.NotContains(t, result, "  - ")
	assert.NotContains(t, result, "  + ")
}

func TestGenDiffInvalidPath(t *testing.T) {
	_, err := GenDiff("nonexistent.json", "testdata/fixture/file2.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse")
}

func TestFormatValueFloatDecimal(t *testing.T) {
	// Создаём временные файлы с дробным числом
	dir := t.TempDir()
	file1 := filepath.Join(dir, "f1.json")
	file2 := filepath.Join(dir, "f2.json")

	// Пишем данные с дробным значением
	require.NoError(t, os.WriteFile(file1, []byte(`{"rate": 3.14}`), 0644))
	require.NoError(t, os.WriteFile(file2, []byte(`{"rate": 2.71}`), 0644))

	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err)
	// Проверяем, что дробные числа отображаются корректно
	assert.Contains(t, result, "3.14")
	assert.Contains(t, result, "2.71")
}

func TestGenDiffUnknownFormat(t *testing.T) {
	dir := t.TempDir()
	txtFile := filepath.Join(dir, "config.txt")
	require.NoError(t, os.WriteFile(txtFile, []byte("key: value"), 0644))

	_, err := GenDiff(txtFile, "testdata/fixture/file2.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format")
}

func TestGenDiffInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	badJSON := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(badJSON, []byte(`{"invalid": json}`), 0644)) // невалидный JSON

	_, err := GenDiff(badJSON, "testdata/fixture/file2.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestFormatValueSlice(t *testing.T) {
	// Этот тест проверяет, что слайсы обрабатываются через default-кейс
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.json")
	f2 := filepath.Join(dir, "b.json")

	// JSON с массивом значений
	require.NoError(t, os.WriteFile(f1, []byte(`{"items": [1, 2, 3]}`), 0644))
	require.NoError(t, os.WriteFile(f2, []byte(`{"items": [4, 5]}`), 0644))

	result, err := GenDiff(f1, f2, "stylish")
	require.NoError(t, err)
	// Проверяем, что слайс отформатирован (через %v)
	assert.Contains(t, result, "items:")
}
