package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// normalizeLineEndings заменяет \r\n на \n для кросс-платформенного сравнения
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

func TestGenDiffJsonFlat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	file1 := filepath.Join(fixtureDir, "file1.json")
	file2 := filepath.Join(fixtureDir, "file2.json")
	expectedFile := filepath.Join(fixtureDir, "expected_stylish.txt")

	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err, "GenDiff should not return error")

	expectedBytes, err := os.ReadFile(expectedFile)
	require.NoError(t, err, "Failed to read expected fixture")
	expected := normalizeLineEndings(string(expectedBytes))
	resultNormalized := normalizeLineEndings(result)

	assert.Equal(t, expected, resultNormalized, "Output should match expected diff")
}

func TestGenDiffYamlFlat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	file1 := filepath.Join(fixtureDir, "file1.yml")
	file2 := filepath.Join(fixtureDir, "file2.yml")
	expectedFile := filepath.Join(fixtureDir, "expected_stylish_yaml.txt")

	result, err := GenDiff(file1, file2, "stylish")
	require.NoError(t, err, "GenDiff should not return error for YAML files")

	expectedBytes, err := os.ReadFile(expectedFile)
	require.NoError(t, err, "Failed to read expected YAML fixture")
	expected := normalizeLineEndings(string(expectedBytes))
	resultNormalized := normalizeLineEndings(result)

	assert.Equal(t, expected, resultNormalized, "YAML output should match expected diff")
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
