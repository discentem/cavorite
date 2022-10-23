FROM golang:1.18-alpine as GOBUILD
# FROM ubuntu:latest as GOBUILD

# This avoids an error about gopath
# WORKDIR /app

RUN echo "$GITHUB_WORKSPACE" >> pwd.txt

RUN cat pwd.txt

ENTRYPOINT [ "sh", "-c", "CGO_ENABLED=0 go build $GITHUB_WORKSPACE -o /usr/bin/pantri_but_go" ]

FROM golang:1.18-alpine as RUNPANTRI

COPY --from=GOBUILD /usr/bin/pantri_but_go /usr/bin/pantri_but_go

ENTRYPOINT ["/usr/bin/pantri_but_go", "-h"]