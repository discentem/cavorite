FROM gcr.io/bazel-public/bazel:latest@sha256:25790da0aaf00d7aa0af59c6a4b5be38778cefedf3254945f79c3cbe329ffefe as bazel

COPY . .

RUN mkdir bazel-cavorite-output

RUN bazel run :cavorite