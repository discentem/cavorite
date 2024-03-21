FROM gcr.io/bazel-public/bazel:latest@sha256:7430f06fff16c8860ec486e1fe0e8c7ce0209a78605b97f41ab6e17f1330ab85 as bazel

COPY . .

RUN mkdir bazel-cavorite-output

RUN bazel run :cavorite