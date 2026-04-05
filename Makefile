# Стандартные цели для сборки, тестирования и поддержки качества кода.
# Использование: make <target>
# Доступные цели: см. `make help`

.PHONY: build lint test fmt tidy clean help

# Сборка исполняемого файла в bin/gendiff
build:
	go build -o bin/gendiff ./cmd/gendiff

# Форматирование кода через go fmt (изменяет файлы на месте)
fmt:
	go fmt ./...

# Запуск линтера golangci-lint (только проверка, без изменений)
lint:
	golangci-lint run

# Запуск тестов с выводом покрытия
test:
	go test -v -cover ./...

# Генерация отчёта о покрытии в coverage.out (для CI/SonarCloud)
cover:
	go test -coverprofile=coverage.out ./...

# Очистка зависимостей и приведение go.mod в порядок
tidy:
	go mod tidy

# Полная очистка артефактов сборки
clean:
	rm -rf bin/* coverage.out

# Комбинированная проверка: форматирование + линтер + тесты
check: fmt lint test

# Справка по доступным целям (парсит комментарии с ##)
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'
