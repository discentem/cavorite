load("@bazel_gazelle//:def.bzl", "gazelle")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

# gazelle:prefix github.com/discentem/cavorite
# gazelle:exclude .sl
# gazelle:proto disable_global
# gazelle:proto package
# gazelle:proto_group go_package
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)

go_library(
    name = "cavorite_lib",
    srcs = ["main.go"],
    importpath = "github.com/discentem/cavorite",
    visibility = ["//visibility:private"],
    deps = [
        "//internal/cli",
        "//internal/program",
        "@com_github_google_logger//:logger",
    ],
)

go_binary(
    name = "cavorite",
    embed = [":cavorite_lib"],
    visibility = ["//visibility:public"],
)
