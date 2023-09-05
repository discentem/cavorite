#!/usr/bin/env bash
GOPACKAGESDRIVER_BAZEL_BUILD_FLAGS=--strategy=GoStdlibList=local
exec bazel run --tool_tag=gopackagesdriver -- @io_bazel_rules_go//go/tools/gopackagesdriver "${@}"