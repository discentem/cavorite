# Adapted from https://github.com/bazelbuild/rules_go/issues/2111#issuecomment-1355927231
# Updated for rules_go 0.51.0+ compatibility
load(
    "@rules_go//go:def.bzl",
    "GoLibrary",
    "go_context",
)

def _output_go_library_srcs_impl(ctx):
    go = go_context(ctx)

    srcs_of_library = []
    importpath = ""
    for src in ctx.attr.deps:
        lib = src[GoLibrary]
        
        # In rules_go 0.51.0+, directly access source files from the provider
        if hasattr(lib, "srcs"):
            srcs_of_library.extend(lib.srcs)
        else:
            # Fallback: try to get srcs from the provider's direct outputs
            if hasattr(lib, "_source"):
                srcs_of_library.extend(lib._source.srcs)
        
        if importpath and lib.importpath != importpath:
            fail(
                "importpath of all deps must match, got {} and {}",
                importpath,
                lib.importpath,
            )
        importpath = lib.importpath

    if len(srcs_of_library) != 1:
        fail("expected exactly one src for library, got {}", len(srcs_of_library))

    if not ctx.attr.out:
        fail("must specify out for now")

    # Run a command to copy the src file to the out location.
    _copy(ctx, srcs_of_library[0], ctx.outputs.out)

def _copy(ctx, in_file, out_file):
    # based on https://github.com/bazelbuild/examples/blob/main/rules/shell_command/rules.bzl
    ctx.actions.run_shell(
        # Input files visible to the action.
        inputs = [in_file],
        # Output files that must be created by the action.
        outputs = [out_file],
        progress_message = "Copying {} to {}".format(in_file.path, out_file.path),
        arguments = [in_file.path, out_file.path],
        command = """cp "$1" "$2" """,
    )

output_go_library_srcs = rule(
    implementation = _output_go_library_srcs_impl,
    attrs = {
        "deps": attr.label_list(
            providers = [GoLibrary],
            aspects = [],
        ),
        "out": attr.output(
            doc = ("Name of output .go file. If not specified, the file name " +
                   "of the generated source file will be used."),
            mandatory = False,
        ),
        "_go_context_data": attr.label(
            default = "@rules_go//:go_context_data",
        ),
    },
    toolchains = ["@rules_go//go:toolchain"],
)

def write_go_proto_srcs(name, go_proto_library, src, visibility = None):
    generated_src = "__generated_" + src
    output_go_library_srcs(
        name = name + "_generated",
        deps = [go_proto_library],
        out = generated_src,
        visibility = ["//visibility:private"],
    )
