FROM golangci/golangci-lint:v1.54-alpine

WORKDIR /app

COPY . .

CMD golangci-lint run --timeout=10m