#!/usr/bin/env bash
set -euo pipefail

CLUSTER_NAME="${KIND_CLUSTER_NAME:-runforge-local}"

if ! command -v kind >/dev/null 2>&1; then
  echo "kind is required but not installed" >&2
  exit 1
fi

if ! kind get clusters | grep -Fxq "${CLUSTER_NAME}"; then
  echo "kind cluster '${CLUSTER_NAME}' does not exist"
  exit 0
fi

kind delete cluster --name "${CLUSTER_NAME}"
echo "kind cluster '${CLUSTER_NAME}' deleted"
