# Verification

## 2026-02-24

### Command
`make help`

### Output
help       Print available targets
tools      Verify required local tools are installed
fmt        Run formatting step
lint       Run lint checks
test       Run test suite

### Command
`make fmt && make lint && make test`

### Output
format check complete (no source formatters configured yet)
lint passed
no tests yet

### Command
`make tidy && make ci`

### Output
ok  	github.com/vivekpradhan/runforge/pkg/version	(cached)

### Command
`make help`

### Output
help       Print available targets
tools      Verify required local tools are installed
fmt        Format Go code
lint       Run golangci-lint
test       Run unit tests
tidy       Tidy Go module files
ci         Run local CI-equivalent checks
