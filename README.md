# RunForge

Bootstrap repository for a Go-based Kubernetes operator project.

## Stack decision

- Language: Go
- Framework direction: Kubebuilder (recommended)

`kubebuilder` is not required for this initial scaffold, but this repo is set up so you can run Kubebuilder initialization next once installed.

## Repository layout

- `docs/` project documentation
- `hack/` helper scripts
- `pkg/version/` small starter package so lint/test run immediately

## Quick start

```bash
make lint
go test ./...
```

## Next step (Kubebuilder)

After installing Kubebuilder:

```bash
kubebuilder init --domain example.com --repo github.com/vivekpradhan/runforge
```
