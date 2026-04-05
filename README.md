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

После сборки утилита будет доступна по пути `./bin/gendiff`.

## Использование

```bash
# Сравнение двух файлов (формат stylish по умолчанию)
./bin/gendiff config1.json config2.json
```

### Явное указание формата
```bash
./bin/gendiff -f plain config1.yml config2.yml
./bin/gendiff --format json config1.json config2.json
```

### Справка
```bash
./bin/gendiff --help
```

### Примеры вывода

#### stylish (по умолчанию):

```bash
{
    common: {
      + follow: false
        setting1: Value 1
      - setting2: 200
    }
}
```

#### plain:
```bash
Property 'common.follow' was added with value: false
Property 'common.setting2' was removed
```

#### json:
```json
{
  "common": {
    "children": {
      "follow": { "status": "added", "value": false },
      "setting2": { "status": "removed", "value": 200 }
    },
    "status": "changed"
  }
}
```

### Разработка

```bash
# Запуск тестов
make test

# Проверка покрытия
make cover

# Линтер
make lint

# Форматирование кода
make fmt

# Все проверки разом
make check
```