FROM golangci/golangci-lint

WORKDIR /app

COPY . .

ENTRYPOINT golangci-lint run --timeout=10m