FROM golang:1.26.5-alpine

# Install golangci-lint
RUN apk add --no-cache git && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.60.3

WORKDIR /app

COPY . .

# golangci-lint compiles cavorite code but chokes on buildvcs info missing so let's disable that
ENV GOFLAGS=-buildvcs=false

ENTRYPOINT golangci-lint run -v --timeout=3m