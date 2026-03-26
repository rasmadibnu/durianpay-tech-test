# Frontend

React + TypeScript + Vite client for the payment dashboard.

## Setup

```bash
npm install
```

## Run

```bash
npm run dev
```

The app starts on `http://localhost:5173`.

## Build

```bash
npm run build
npm run preview
```

## Tests

Run the frontend test suite:

```bash
npm test
```

Run tests in watch mode:

```bash
npm run test:watch
```

Run the TypeScript check:

```bash
npx tsc --noEmit
```

## Test Stack

- Vitest
- Testing Library
- JSDOM
- Shared setup in `src/test/setup.ts`

## Current Coverage

- Login page
- Dashboard page
- Merchants page
- Users page
