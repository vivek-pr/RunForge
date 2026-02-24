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

### Command
`make manifests`

### Output
"/Users/vivekpradhan/vscode/InterViewPreparation/RunForge/bin/controller-gen" rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases

### Command
`make kind-up`

### Output
No kind clusters found.
Creating cluster "runforge-local" ...
...
kind cluster 'runforge-local' is ready

### Command
`make install`

### Output
"/Users/vivekpradhan/vscode/InterViewPreparation/RunForge/bin/controller-gen" rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
customresourcedefinition.apiextensions.k8s.io/aijobs.runforge.runforge.io created

### Command
`kubectl get crd | grep aijobs`

### Output
aijobs.runforge.runforge.io   2026-02-24T07:59:55Z

### Command
`kubectl api-resources | grep -i aijob`

### Output
aijobs                                           runforge.runforge.io/v1alpha1     true         AIJob

### Command
`make uninstall`

### Output
"/Users/vivekpradhan/vscode/InterViewPreparation/RunForge/bin/controller-gen" rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
customresourcedefinition.apiextensions.k8s.io "aijobs.runforge.runforge.io" deleted

### Command
`make generate && make manifests && make install`

### Output
"/Users/vivekpradhan/vscode/InterViewPreparation/RunForge/bin/controller-gen" object:headerFile="hack/boilerplate.go.txt" paths="./..."
"/Users/vivekpradhan/vscode/InterViewPreparation/RunForge/bin/controller-gen" rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases
customresourcedefinition.apiextensions.k8s.io/aijobs.runforge.runforge.io configured

### Command
`kubectl explain aijob.spec`

### Output
GROUP:      runforge.runforge.io
KIND:       AIJob
VERSION:    v1alpha1

FIELDS include: `image` (required), `command`, `args`, `env`, `envFrom`, `resources`, `restartPolicy`, `backoffLimit`, `activeDeadlineSeconds`, `ttlSecondsAfterFinished`, `nodeSelector`, `tolerations`, `affinity`, `serviceAccountName`.

### Command
`kubectl apply -f examples/aijob-success.yaml`

### Output
aijob.runforge.runforge.io/aijob-success created

### Command
`kubectl get aijob -o yaml`

### Output
List contains `aijob-success` with persisted `spec` fields from `examples/aijob-success.yaml` (image, command/args, env/envFrom, resources, restartPolicy, backoffLimit, deadlines/ttl, scheduling, serviceAccountName).
