# go-specification-lang

Инструмент для добавления динамических проверок контрактов (предусловий, постусловий, инвариантов) в Go-программы на основе аннотаций в стиле ACSL.

## Требования

- Go 1.21+
- (опционально) [gvm](https://github.com/moovweb/gvm) для управления версиями Go
- (опционально) Java — для перегенерации ANTLR-парсера

## Сборка

### Через Makefile (рекомендуется)

`make` автоматически запускает `go generate` для сборки ANTLR-парсера, а затем компилирует бинарник:

```bash
cd framework
make build
# бинарник появится в framework/build/framework
```

### Напрямую через `go build`

Перед сборкой необходимо вручную сгенерировать ANTLR-парсер (требует Java):

```bash
cd framework
go generate ./...
go build -o build/framework ./cli/main.go
```

## CLI

CLI-приложение (`framework/cli/main.go`) принимает Go-файлы или целые проекты, прогоняет их через один или несколько **плагинов** и записывает результат. Плагины трансформируют AST исходного кода — например, вставляют динамические проверки контрактов. Можно подключать несколько плагинов одновременно через повторный флаг `--plugin`.

Полная справка по любой команде:

```bash
go run framework/cli/main.go --help
go run framework/cli/main.go file --help
go run framework/cli/main.go project --help
```

Список зарегистрированных плагинов:

```bash
go run framework/cli/main.go plugins
```

### Обработка одного файла

```bash
go run framework/cli/main.go file --plugin <plugin> <source.go> <output.go>
```

### Обработка целого проекта

Обходит все `.go`-файлы в директории проекта, применяет плагины и записывает результат в `<output-dir/>`. Дополнительно создаёт файл `gpma_overlay.json`, который можно передать в `go run --overlay` для запуска изменённого кода без перезаписи оригиналов:

```bash
go run framework/cli/main.go project --plugin <plugin> <project-dir/> <output-dir/>
```

## Плагин `dynamic_checks`

Плагин `dynamic_checks` читает ACSL-аннотации в комментариях (`/*@ ... */`) и вставляет в код вызовы `checks.Check(...)` для проверки предусловий (`requires`), постусловий (`ensures`) и утверждений (`assert`) во время выполнения.

### Генерация overlay и запуск

Сначала сгенерируйте overlay — директорию с изменёнными файлами и файлом `gpma_overlay.json`:

```bash
go run framework/cli/main.go project --plugin dynamic_checks \
  ./samples/projects/simple/ ./overlays/
```

После этого можно запустить программу без проверок (оригинальный код):

```bash
go run samples/projects/simple/main.go
```

Или с включёнными проверками через overlay (оригинальные файлы не изменяются):

```bash
go run --overlay overlays/gpma_overlay.json samples/projects/simple/main.go
```

### Режим строгости (`strict_mode`)

**Строгий** (по умолчанию) — предусловия и постусловия проверяются безусловно:

```bash
go run framework/cli/main.go project -p dynamic_checks \
  --plugin-config 'dynamic_checks:strict_mode:true' \
  samples/projects/interactive/ overlays/
```

**Нестрогий** — постусловия проверяются только если предусловия выполнены; нарушение предусловий само по себе не является ошибкой:

```bash
go run framework/cli/main.go project -p dynamic_checks \
  --plugin-config 'dynamic_checks:strict_mode:false' \
  samples/projects/interactive/ overlays/
```

### Режим реакции на нарушения (`CHECK_MODE`)

Поведение при нарушении контракта управляется переменной окружения `CHECK_MODE`:

| Значение  | Поведение                               |
|-----------|-----------------------------------------|
| `panic`   | паника (по умолчанию)                   |
| `console` | вывод в stderr                          |
| `file`    | запись в файл (требует `CHECK_LOG_PATH`) |
| `silent`  | игнорировать нарушения                  |

```bash
# Запись нарушений в файл
CHECK_MODE=file CHECK_LOG_PATH=violations.txt \
  go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go

# Тихий режим — нарушения игнорируются
CHECK_MODE=silent go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go

# Паника при нарушении
CHECK_MODE=panic go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go
```

## Тесты

Запуск всех тестов фреймворка:

```bash
cd framework
make test
# или напрямую:
go test ./...
```

### Golden-тесты плагина `dynamic_checks`

Golden-тесты проверяют, что плагин генерирует ожидаемый код для каждого проекта из `samples/projects/`.

**Запуск проверки:**

```bash
cd framework
go test ./plugin/dynamic_checks/... -v
```

**Обновление golden-файлов** (после изменения логики генерации):

```bash
cd framework
go test ./plugin/dynamic_checks/... -update
```

Флаг `-update` перезаписывает файлы в `framework/plugin/dynamic_checks/testdata/golden/` актуальным выводом плагина. После обновления зафиксируйте изменения в git, чтобы зафиксировать новое ожидаемое поведение.
