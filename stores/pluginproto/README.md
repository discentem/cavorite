# Generate plugin.pb.go for non-bazel builds

For non-bazel builds and [linting](../../../_ci/lint/Dockerfile), this file needs to exist in the source tree.

To generate plugin.pb.go run `bazel run //internal/stores/pluginproto:write_generated_protos`