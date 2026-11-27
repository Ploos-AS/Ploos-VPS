# Local development

Requires Go 1.25 or newer. Run `go run ./cmd/api`, then request `http://127.0.0.1:8080/healthz` or `/readyz`. Override the bind address with `PLOOS_VPS_LISTEN`. Run `go test -race ./...` and `go vet ./...`. Readiness currently checks only the process; database readiness will be added with persistence. No authentication, provisioning or production deployment yet.
