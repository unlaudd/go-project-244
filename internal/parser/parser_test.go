package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseFileJSON(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.json")
	require.NoError(t, os.WriteFile(file, []byte(`{"key":"val","num":42}`), 0644))

	res, err := ParseFile(file)
	require.NoError(t, err)
	assert.Equal(t, "val", res["key"])
	assert.Equal(t, float64(42), res["num"])
}

func TestParseFileYAML(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.yaml")
	require.NoError(t, os.WriteFile(file, []byte("key: val\nnum: 42"), 0644))

	res, err := ParseFile(file)
	require.NoError(t, err)
	assert.Equal(t, "val", res["key"])
	assert.Equal(t, 42, res["num"])
}

func TestParseFileUnknownExt(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "test.txt")
	require.NoError(t, os.WriteFile(file, []byte("data"), 0644))

	_, err := ParseFile(file)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format")
}

func TestParseFileInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(file, []byte(`{invalid}`), 0644))

	_, err := ParseFile(file)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}

func TestParseFileNotFound(t *testing.T) {
	_, err := ParseFile("nonexistent.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read file")
}
