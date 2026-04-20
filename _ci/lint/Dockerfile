FROM golangci/golangci-lint:v1.64.8
WORKDIR /app

COPY . .

# golangci-lint compiles cavorite code but chokes on buildvcs info missing so let's disable that
ENV GOFLAGS=-buildvcs=false

ENTRYPOINT golangci-lint run -v --timeout=3m