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

	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	app := NewApp()
	runErr := app.Run([]string{"gendiff", f1, f2})

	// ✅ errcheck: явно проверяем ошибки закрытия и копирования
	require.NoError(t, w.Close(), "failed to close write end of pipe")
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, copyErr := io.Copy(&buf, r)
	require.NoError(t, copyErr, "failed to read from pipe")

	require.NoError(t, runErr, "CLI execution should not return error")
	assert.Contains(t, buf.String(), "key")
}

func TestCLIErrorPath(t *testing.T) {
	// Проверяем, что при отсутствии аргументов возвращается ожидаемая ошибка
	// Примечание: печать "Error: ..." в stderr происходит в main(),
	// который мы не тестируем напрямую (из-за os.Exit)

	app := NewApp()
	err := app.Run([]string{"gendiff"}) // нет аргументов → ошибка

	// Проверяем только возврат ошибки, а не вывод в stderr
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "two file paths are required")
}
