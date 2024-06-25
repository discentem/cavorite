FROM gcr.io/bazel-public/bazel:latest@sha256:65e2c7b5814b000163817ddcfedb16ce3abe4ddd745d683ad399cab14398e4e2 as bazel

COPY . .

RUN mkdir bazel-cavorite-output

RUN bazel run :cavorite