// Package main содержит точку входа CLI-утилиты gendiff и её тесты.
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

// TestCLIMissingArgs проверяет, что утилита возвращает ошибку при отсутствии аргументов.
func TestCLIMissingArgs(t *testing.T) {
	app := NewApp()
	err := app.Run([]string{"gendiff"})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "two file paths are required")
}

// TestCLISuccess проверяет успешное выполнение утилиты с валидными аргументами.
// Тест перехватывает stdout, чтобы проверить вывод без загрязнения консоли.
func TestCLISuccess(t *testing.T) {
	dir := t.TempDir()
	f1 := filepath.Join(dir, "a.json")
	f2 := filepath.Join(dir, "b.json")
	require.NoError(t, os.WriteFile(f1, []byte(`{"key":"val1"}`), 0644))
	require.NoError(t, os.WriteFile(f2, []byte(`{"key":"val2"}`), 0644))

	// Перехватываем stdout для проверки вывода
	// Сохраняем оригинал, чтобы восстановить после теста
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	app := NewApp()
	runErr := app.Run([]string{"gendiff", f1, f2})

	// Восстанавливаем stdout и читаем перехваченный вывод
	// Ошибки закрытия и копирования проверяем явно (требование errcheck)
	require.NoError(t, w.Close(), "failed to close write end of pipe")
	os.Stdout = oldStdout

	var buf bytes.Buffer
	_, copyErr := io.Copy(&buf, r)
	require.NoError(t, copyErr, "failed to read from pipe")

	require.NoError(t, runErr, "CLI execution should not return error")
	assert.Contains(t, buf.String(), "key")
}

// TestCLIErrorPath проверяет обработку ошибки валидации аргументов.
// Примечание: вывод "Error: ..." в stderr происходит в main(),
// который не тестируется напрямую из-за os.Exit(1).
// Здесь проверяем только возврат ошибки из Action-функции.
func TestCLIErrorPath(t *testing.T) {
	app := NewApp()
	err := app.Run([]string{"gendiff"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "two file paths are required")
}
