// Package code содержит интеграционные тесты для библиотеки gendiff.
// Проверяет сквозную работу: парсинг → построение дерева → форматирование.
package code

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// normalizeLineEndings приводит окончания строк к Unix-стилю (\n).
// Необходим для кроссплатформенных тестов: Windows использует \r\n,
// а фикстуры хранятся с \n, что вызывает ложные падения тестов.
func normalizeLineEndings(s string) string {
	return strings.ReplaceAll(s, "\r\n", "\n")
}

// loadExpectedFixture загружает ожидаемый вывод из файла фикстуры.
// Автоматически нормализует окончания строк для кроссплатформенности.
// Помечен как t.Helper(), чтобы в логах ошибок показывалась строка вызова теста, а не этой функции.
func loadExpectedFixture(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read fixture: %s", path)
	return normalizeLineEndings(string(content))
}

// TestGenDiffNestedStylish проверяет интеграционный сценарий для формата stylish.
// Сравнивает результат работы GenDiff с эталонным файлом-фикстурой.
// Тестирует рекурсивное сравнение вложенных структур с правильными отступами.
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

// TestGenDiffNestedPlain проверяет интеграционный сценарий для формата plain.
// Убеждаемся, что вывод содержит полные пути до ключей, кавычки для строк
// и маркер [complex value] для вложенных объектов.
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
	// plain-форматер намеренно пропускает неизменённые ключи, чтобы вывод был компактным
	assert.NotContains(t, result, "unchanged")
}

// TestGenDiffDefaultFormat проверяет, что пустая строка в параметре format
// подставляет дефолтное значение "stylish". Это важно для обратной совместимости
// и корректной работы библиотеки при программном вызове без указания формата.
func TestGenDiffDefaultFormat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")
	result, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.json"),
		filepath.Join(fixtureDir, "file2_nested.json"),
		"", // пустой формат → должен сработать дефолт
	)
	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result, "{\n"), "Default format should be stylish")
}

// TestGenDiffUnknownFormat проверяет обработку неподдерживаемого формата вывода.
// Ошибка должна возвращаться на уровне форматера, а не парсера, поэтому передаём валидные файлы.
func TestGenDiffUnknownFormat(t *testing.T) {
	fixtureDir := filepath.Join("testdata", "fixture")

	_, err := GenDiff(
		filepath.Join(fixtureDir, "file1_nested.json"),
		filepath.Join(fixtureDir, "file2_nested.json"),
		"markdown", // неподдерживаемый формат
	)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown format: markdown")
}

// TestGenDiffInvalidPath проверяет обработку ошибки при чтении несуществующего файла.
// Допускаем два варианта сообщения об ошибке, так как парсер может упасть
// либо на этапе разрешения пути, либо на этапе чтения файла.
func TestGenDiffInvalidPath(t *testing.T) {
	_, err := GenDiff("nonexistent.json", "testdata/fixture/file1_nested.json", "stylish")
	assert.Error(t, err)
	// Гибкая проверка: ошибка может возникнуть на разных этапах (path resolution или file read)
	assert.True(t, strings.Contains(err.Error(), "failed to read") || strings.Contains(err.Error(), "failed to parse"))
}

// TestGenDiffInvalidJSON проверяет обработку синтаксически невалидного JSON.
// Убеждаемся, что ошибка парсинга корректно прокидывается наверх с понятным сообщением.
func TestGenDiffInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	badJSON := filepath.Join(dir, "bad.json")
	require.NoError(t, os.WriteFile(badJSON, []byte(`{invalid}`), 0644))

	_, err := GenDiff(badJSON, "testdata/fixture/file1_nested.json", "stylish")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse JSON")
}
