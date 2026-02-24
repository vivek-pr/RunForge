# AIJob API (`runforge.runforge.io/v1alpha1`)

## Spec

- `image` (string, required): container image to run.
- `command` ([]string): optional entrypoint override.
- `args` ([]string): optional command arguments.
- `env` ([]object): explicit environment variables.
  - `name` (string, required)
  - `value` (string, optional)
- `envFrom` ([]object): environment imports.
  - `configMapRef.name` (string, optional)
  - `secretRef.name` (string, optional)
- `resources` (object): optional requests/limits.
  - `requests.cpu`, `requests.memory`, `requests.gpu`
  - `limits.cpu`, `limits.memory`, `limits.gpu`
- `restartPolicy` (string): `Never` or `OnFailure` (default `Never`).
- `backoffLimit` (int32): retry budget before failed.
- `activeDeadlineSeconds` (int64): max runtime.
- `ttlSecondsAfterFinished` (int32): TTL after completion.
- `nodeSelector` (map[string]string): scheduling selector.
- `tolerations` ([]Toleration): taint tolerations.
- `affinity` (Affinity): pod affinity rules.
- `serviceAccountName` (string): optional pod identity.

## Status

- `observedGeneration` (int64): last processed generation.
- `phase` (string): `Pending`, `Running`, `Succeeded`, `Failed`.
- `jobName` (string): backing Kubernetes Job name.
- `conditions` ([]object): condition timeline.
  - `type`
  - `reason`
  - `message`
  - `lastTransitionTime`
- `startTime` (Time)
- `completionTime` (Time)
- `lastError` (string)

## Example

```yaml
apiVersion: runforge.runforge.io/v1alpha1
kind: AIJob
metadata:
  name: summarize-docs
spec:
  image: ghcr.io/acme/ai-worker:1.2.3
  command: ["/app/worker"]
  args: ["run", "--task=summarize"]
  env:
    - name: MODEL
      value: gpt-4.1-mini
  envFrom:
    - secretRef:
        name: ai-api-keys
  resources:
    requests:
      cpu: "500m"
      memory: "512Mi"
    limits:
      cpu: "1"
      memory: "1Gi"
      gpu: "1"
  restartPolicy: OnFailure
  backoffLimit: 2
  activeDeadlineSeconds: 1800
  ttlSecondsAfterFinished: 600
  serviceAccountName: runforge-ai
```
