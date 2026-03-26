# Payment Dashboard

A full-stack payment monitoring dashboard built with Go (backend) and React + TypeScript (frontend).

## Prerequisites

- Go 1.21+ (with CGO support — requires a C compiler like `gcc` or `clang`)
- Node.js 20+
- Docker & Docker Compose (optional, for containerized setup)

## Quick Start (Docker Compose)

The easiest way to run both services:

```bash
make run
# or
docker compose up --build
```

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## Manual Setup

### Backend

First-time setup:

```bash
cd backend
cp env.sample .env
make dep
make tool-openapi
make openapi-gen
make gen-secret
```

Start the API:

```bash
cd backend
CGO_ENABLED=1 go run main.go
```

Or via the backend makefile:

```bash
make run
```

The backend starts on http://localhost:8080. On first run it automatically:

- Creates a SQLite database (`dashboard.db`)
- Seeds 2 test users and 50 sample payments

### Frontend

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

The frontend starts on http://localhost:5173 (Vite dev server).

### Frontend Production Build

```bash
cd frontend
npm run build
npm run preview
```

## Test Accounts

| Email              | Password | Role      |
| ------------------ | -------- | --------- |
| cs@test.com        | password | cs        |
| operation@test.com | password | operation |

## Data Seeding

Data is automatically seeded on first run:

- **Users**: 2 accounts (cs and operation roles)
- **Payments**: 50 sample payments with randomized merchants, amounts, statuses, and dates

To reset data, delete `backend/dashboard.db` and restart the backend.

## Running Tests

### Backend Tests

```bash
cd backend
CGO_ENABLED=1 go test ./... -v
```

### Frontend Unit Tests

```bash
cd frontend
npm test
```

### Frontend Type Check

```bash
cd frontend
npx tsc --noEmit
```

### All Tests

```bash
make test
```

### Make Targets

```bash
make test-backend
make test-frontend
make test-frontend-typecheck
```

## Testing Strategy

### Backend

- **Repository tests** cover merchant, payment, and user persistence behavior
- **Usecase tests** cover auth, merchant, payment, and user business rules
- Tests use an **in-memory SQLite** database where persistence behavior matters

### Frontend

- **Vitest + Testing Library** cover shared components, auth context, and page flows
- **Shared test setup** is defined in `frontend/src/test/setup.ts`
- **Type safety** is checked separately with `npx tsc --noEmit`

## API Documentation

The API follows the OpenAPI specification defined in `openapi.yaml`.

## Make Commands

```bash
make help           # Show all available commands
make install        # Install all dependencies
make run-backend    # Run backend dev server
make run-frontend   # Run frontend dev server
make run            # Run both with docker-compose
make build          # Build docker images
make test-backend   # Run backend Go tests
make test-frontend  # Run frontend Vitest suite
make test-frontend-typecheck # Run frontend TypeScript check
make test           # Run backend + frontend tests
make clean          # Clean build artifacts
```
