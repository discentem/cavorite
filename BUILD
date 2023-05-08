load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix github.com/discentem/cavorite
# gazelle:exclude .sl
gazelle(name = "gazelle")

go_library(
    name = "cavorite_lib",
    srcs = ["main.go"],
    importpath = "github.com/discentem/cavorite",
    visibility = ["//visibility:private"],
    deps = [
        "//internal/cmd/root",
        "@com_github_google_logger//:logger",
    ],
)

go_binary(
    name = "cavorite",
    embed = [":cavorite_lib"],
    visibility = ["//visibility:public"],
)
