load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/discentem/pantri_but_go
# gazelle:exclude .sl
gazelle(name = "gazelle")

go_library(
    name = "pantri_but_go_lib",
    srcs = ["main.go"],
    importpath = "github.com/discentem/pantri_but_go",
    visibility = ["//visibility:private"],
    deps = [
        "//internal/config",
        "//internal/stores",
        "@com_github_google_logger//:logger",
        "@com_github_spf13_afero//:afero",
        "@com_github_urfave_cli_v2//:cli",
    ],
)

go_binary(
    name = "pantri_but_go",
    embed = [":pantri_but_go_lib"],
    visibility = ["//visibility:public"],
)
