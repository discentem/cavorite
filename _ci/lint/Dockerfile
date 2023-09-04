FROM golangci/golangci-lint:v1.54.1
WORKDIR /app

COPY . .

ENV GOFLAGS=-buildvcs=false

ENTRYPOINT golangci-lint run -v --timeout=10m