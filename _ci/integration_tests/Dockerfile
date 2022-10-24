# https://stackoverflow.com/questions/72904999/create-minio-docker-image-with-content
FROM golang:1.18-alpine as GOBUILDER

WORKDIR /app

COPY . .

RUN echo $(which go) >> pwd.txt

RUN cat pwd.txt

RUN cd ./_ci/integration_tests/scripts && CGO_ENABLED=0 go build -o /usr/bin/create_bucket

FROM docker.io/minio/minio:latest

COPY --from=GOBUILDER /usr/bin/create_bucket /usr/bin/create_bucket

COPY . .

RUN echo $(ls) >> pwd.txt

RUN cat pwd.txt

RUN echo "$GITHUB_WORKSPACE" >> pwd.txt

RUN cat pwd.txt

ENTRYPOINT ["sh", "./_ci/integration_tests/run_minio_test.sh"]