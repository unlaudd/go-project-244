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
