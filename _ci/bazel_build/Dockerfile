FROM gcr.io/bazel-public/bazel:latest@sha256:ad5916a7e27eeddbddb003d7671d5a8202c7346e3c682c059b27467b47c0f6d4 as bazel

COPY . .

RUN mkdir bazel-cavorite-output

RUN bazel run :cavorite