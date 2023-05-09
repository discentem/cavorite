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
        "//internal/cli",
        "@com_github_google_logger//:logger",
    ],
)

go_binary(
    name = "pantri_but_go",
    embed = [":pantri_but_go_lib"],
    visibility = ["//visibility:public"],
)
