FROM gcr.io/bazel-public/bazel:latest@sha256:33ffeb123eb231271309835c6348103ae15520856e74aa4d35b0c49c33f9619a as bazel

COPY . .

RUN mkdir bazel-cavorite-output

RUN bazel run :cavorite