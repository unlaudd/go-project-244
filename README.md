# gendiff

CLI-утилита для сравнения конфигурационных файлов (JSON, YAML) и вывода различий.

Поддерживает три формата вывода:
- **stylish** (по умолчанию) — иерархический вывод с отступами и маркерами `+` / `-`
- **plain** — плоский список изменений с полными путями: `Property 'common.setting' was updated...`
- **json** — структурированный JSON для интеграции с другими системами

### Статус тестов и линтера Hexlet:
[![Actions Status](https://github.com/unlaudd/go-project-244/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/unlaudd/go-project-244/actions)

## Демо
[![asciicast](https://asciinema.org/a/IlUcGgVUFNjwZoad.svg)](https://asciinema.org/a/IlUcGgVUFNjwZoad)

## Статус CI и качества кода
[![CI](https://github.com/unlaudd/go-project-244/actions/workflows/ci.yml/badge.svg)](https://github.com/unlaudd/go-project-244/actions/workflows/ci.yml)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=unlaudd_go-project-244&metric=alert_status)](https://sonarcloud.io/dashboard?id=unlaudd_go-project-244)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=unlaudd_go-project-244&metric=coverage)](https://sonarcloud.io/dashboard?id=unlaudd_go-project-244)

## Установка

```bash
# Сборка из исходников
make build

# Или вручную
go build -o bin/gendiff ./cmd/gendiff
```
