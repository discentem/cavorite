FROM gcr.io/bazel-public/bazel:latest@sha256:3453374b7d76157c876e4e4134f657d8acf6ce63e65842c2bf6275ba13e19e85 as bazel

COPY . .

RUN mkdir bazel-cavorite-output

RUN bazel run :cavorite