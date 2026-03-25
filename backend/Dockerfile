FROM golang:1.23-alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .
COPY openapi.yaml /openapi.yaml

RUN CGO_ENABLED=1 go build -o bin/server main.go

FROM alpine:3.19

RUN apk add --no-cache sqlite-libs ca-certificates

WORKDIR /app

COPY --from=builder /app/bin/server .
COPY --from=builder /openapi.yaml ./openapi.yaml

ENV OPENAPIYAML_LOCATION=./openapi.yaml

EXPOSE 8080

CMD ["./server"]
