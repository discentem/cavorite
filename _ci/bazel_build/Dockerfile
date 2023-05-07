FROM gcr.io/bazel-public/bazel:latest@sha256:91d0c1fcbb9e46c562dab1b99eba360188fe5d1e2890f502b7af509219a3ff5d as bazel

COPY . .

RUN mkdir bazel-pantri-output

RUN bazel --output_user_root /home/ubuntu/bazel-pantri-output run :pantri_but_go

RUN chown ubuntu /github/home/.cache/bazel/