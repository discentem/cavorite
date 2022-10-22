FROM golangci/golangci-lint

WORKDIR /app

COPY . .

ENTRYPOINT cd $GITHUB_WORKSPACE && golangci-lint run 