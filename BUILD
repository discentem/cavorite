load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/example/project
gazelle(name = "gazelle")

go_library(
    name = "project_lib",
    srcs = ["main.go"],
    importpath = "github.com/example/project",
    visibility = ["//visibility:private"],
    deps = [
        "@com_github_google_logger//:go_default_library",
        "@com_github_spf13_afero//:go_default_library",
        "@com_github_urfave_cli_v2//:go_default_library",
    ],
)

go_binary(
    name = "project",
    embed = [":project_lib"],
    visibility = ["//visibility:public"],
)
