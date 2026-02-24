#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${KIND_CLUSTER_NAME:-runforge-local}"
CONFIG_PATH="${KIND_CONFIG_PATH:-hack/kind/kind-config.yaml}"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required but not installed" >&2
  exit 1
fi

if kind get clusters | grep -Fxq "${CLUSTER_NAME}"; then
  echo "kind cluster '${CLUSTER_NAME}' already exists"
  kubectl config use-context "kind-${CLUSTER_NAME}" >/dev/null
  echo "context set to kind-${CLUSTER_NAME}"
  exit 0
fi

kind create cluster --name "${CLUSTER_NAME}" --config "${CONFIG_PATH}" --wait 120s
kubectl config use-context "kind-${CLUSTER_NAME}" >/dev/null
echo "kind cluster '${CLUSTER_NAME}' is ready"
