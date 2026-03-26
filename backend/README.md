# Backend

Go API for the payment dashboard.

## Prerequisites

- Go 1.21+
- CGO enabled
- A local C toolchain for `sqlite3`

## Setup

```bash
cp env.sample .env
make dep
make tool-openapi
make openapi-gen
make gen-secret
```

## Run

```bash
CGO_ENABLED=1 go run main.go
```

Or via the backend makefile:

```bash
make run
```

The API starts on `http://localhost:8080`.

## OpenAPI

Generate server/types from the root spec:

```bash
make tool-openapi
make openapi-gen
```

## Tests

Run all backend tests:

```bash
CGO_ENABLED=1 go test ./... -v
```

Run one package:

```bash
go test ./internal/module/payment/usecase -v
```

## Coverage Scope

- Repository tests for merchant, payment, and user modules
- Usecase tests for auth, merchant, payment, and user modules
