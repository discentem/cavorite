# https://stackoverflow.com/questions/72904999/create-minio-docker-image-with-content
FROM docker.io/minio/minio:latest as MINIOSOURCE

RUN echo $(which minio) >> pwd.txt

RUN cat pwd.txt

FROM golang:1.18-alpine

WORKDIR /app

COPY . .

# COPY --from=MINIOSOURCE /usr/bin/pantri_not_go /usr/bin/pantri_not_go

# RUN CGO_ENABLED=0 go build -o /usr/bin/pantri_not_go

# RUN echo $(ls) >> pwd.txt

# RUN cat pwd.txt

# RUN echo "$GITHUB_WORKSPACE" >> pwd.txt

# RUN cat pwd.txt

ENTRYPOINT ["sh", "echo", "'hello'"]