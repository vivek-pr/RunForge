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

### Command
`make kind-up`

### Output
No kind clusters found.
Creating cluster "runforge-local" ...
...
kind cluster 'runforge-local' is ready

### Command
`kubectl get nodes -o wide`

### Output
NAME                           STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                         KERNEL-VERSION     CONTAINER-RUNTIME
runforge-local-control-plane   Ready    control-plane   29s   v1.33.1   172.18.0.2    <none>        Debian GNU/Linux 12 (bookworm)   6.12.65-linuxkit   containerd://2.1.1

### Command
`make kind-down`

### Output
Deleting cluster "runforge-local" ...
Deleted nodes: ["runforge-local-control-plane"]
kind cluster 'runforge-local' deleted
