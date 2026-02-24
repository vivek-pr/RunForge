# Developer Decisions

## 2026-02-24

- Reset repository to a clean baseline before re-bootstrap.
- Chose a tool-agnostic Makefile so `make test` succeeds even before application code exists.
- Made developer tracking files mandatory and validated them in `make lint`.
- Pinned `golangci-lint` to `v1.64.8` in both local Makefile and CI workflow to keep lint rules consistent across environments.
- Pinned `kubebuilder` to `v4.12.0` as the scaffold baseline for this repo because it supports Kubernetes `v1.35` APIs and matches the locally installed toolchain.
- AIJob spec immutability policy: core execution fields are immutable after create (`image`, `command`, `args`, `env`, `envFrom`, `resources`, `restartPolicy`, scheduling fields). Mutable fields are retained for lifecycle management (`ttlSecondsAfterFinished` and metadata labels/annotations).
