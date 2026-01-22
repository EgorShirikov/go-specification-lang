# go-specification

Для тех, кто будет работать с go, советую установить gvm (Go Version Manager)
Внутри проекта `go applymod`

# generate overlay

`go run framework/cli/main.go project --plugin dynamic_checks ./samples/projects/simple/ ./overlays`

# run without checks

`go run samples/projects/simple/main.go`

# run with enabled checks

`go run --overlay overlays/gpma_overlay.json samples/projects/simple/main.go`

# different check verbosity modes

CHECK_MODE="silent|console|panic|file"

`CHECK_MODE=file CHECK_LOG_PATH=violations.txt go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go`
`CHECK_MODE=silent go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go`
`CHECK_MODE=panic go run --overlay overlays/gpma_overlay.json samples/projects/interactive/main.go`