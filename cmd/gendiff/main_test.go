package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCLIMissingArgs(t *testing.T) {
	app := NewApp()
	err := app.Run([]string{"gendiff"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "two file paths are required")
}

func TestCLISuccess(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.json")
	f2 := filepath.Join(dir, "b.json")
	require.NoError(t, os.WriteFile(f1, []byte(`{"key":"val1"}`), 0644))
	require.NoError(t, os.WriteFile(f2, []byte(`{"key":"val2"}`), 0644))

	// Перехватываем stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	app := NewApp()
	err := app.Run([]string{"gendiff", f1, f2})

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)

	require.NoError(t, err)
	assert.Contains(t, buf.String(), "key")
}
