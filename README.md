# go-specification

Для тех, кто будет работать с go, советую установить gvm (Go Version Manager)
Внутри проекта `go applymod`

# generate overlay

`go run framework/cli/main.go project --plugin dynamic_checks ./samples/projects/simple/ ./overlays`

# run without checks

```bash
go run samples/projects/simple/main.go
```

# run with enabled checks

```bash
go run --overlay overlays/gpma_overlay.json samples/projects/simple/main.go
```

# different check verbosity modes

CHECK_MODE="silent|console|panic|file"

```bash
CHECK_MODE=file CHECK_LOG_PATH=violations.txt go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go
```

```bash
CHECK_MODE=silent go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go
```

```bash
CHECK_MODE=panic go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go
```