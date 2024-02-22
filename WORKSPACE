load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "80a98277ad1311dacd837f9b16db62887702e9f1d1c4c9f796d0121a46c8e184",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.46.0/rules_go-v0.46.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.46.0/rules_go-v0.46.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

http_archive(
    name = "bazel_gazelle",
    sha256 = "32938bda16e6700063035479063d9d24c60eda8d79fd4739563f50d331cb3209",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.35.0/bazel-gazelle-v0.35.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

############################################################
# Define your own dependencies here using go_repository.
# Else, dependencies declared by rules_go/gazelle will be used.
# The first declaration of an external repository "wins".
############################################################

http_archive(
    name = "bazel_skylib",
    sha256 = "cd55a062e763b9349921f0f5db8c3933288dc8ba4f76dd9416aac68acee3cb94",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.5.0/bazel-skylib-1.5.0.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.5.0/bazel-skylib-1.5.0.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

go_repository(
    name = "com_github_fsouza_fake_gcs_server",
    importpath = "github.com/fsouza/fake-gcs-server",
    sum = "h1:UXPNcCDrS3j7T3CBjZm+x9zBD9Cc5jqQZlm0mcbZAIE=",
    version = "v1.45.1",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:oxx1eChJGI6Uks2ZC4W1zpLlVgqB8ner4EuQwV4Ik1Y=",
    version = "v1.9.2",
)

go_repository(
    name = "com_github_gorilla_handlers",
    importpath = "github.com/gorilla/handlers",
    sum = "h1:9lRY6j8DEeeBT10CvO9hGW0gmky0BprnvDI5vfhUHH4=",
    version = "v1.5.1",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:i40aqfkR1h2SlN9hojwV5ZA91wcXFOvkdNIeFDP5koI=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_google_renameio_v2",
    importpath = "github.com/google/renameio/v2",
    sum = "h1:UifI23ZTGY8Tt29JbYFiuyIU3eX+RNFtUwefq9qAhxg=",
    version = "v2.0.0",
)

go_repository(
    name = "com_github_felixge_httpsnoop",
    importpath = "github.com/felixge/httpsnoop",
    sum = "h1:s/nj+GCswXYzN5v2DpNMuMQYe+0DDwt5WVCU6CWBdXk=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_pkg_xattr",
    importpath = "github.com/pkg/xattr",
    sum = "h1:5883YPCtkSd8LFbs13nXplj9g9tlrwoJRjgpgMu1/fE=",
    version = "v0.4.9",
)

go_repository(
    name = "com_github_sourcegraph_conc",
    importpath = "github.com/sourcegraph/conc",
    sum = "h1:OQTbbt6P72L20UqAkXXuLOj79LfEanQ+YQFNpLA9ySo=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_hashicorp_go_multierror",
    importpath = "github.com/hashicorp/go-multierror",
    sum = "h1:H5DkEtf6CXdFp0N0Em5UCwQpXMWke8IA0+lD48awMYo=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_hashicorp_errwrap",
    importpath = "github.com/hashicorp/errwrap",
    sum = "h1:OxrOeh75EUXMY8TBjag2fzXGZ40LB6IKw45YeGUDY2I=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go",
    importpath = "github.com/Azure/azure-sdk-for-go",
    sum = "h1:fcYLmCpyNYRnvJbPerq7U0hS+6+I79yEDJBqVNcqUzU=",
    version = "v68.0.0+incompatible",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_storage_azblob",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob",
    sum = "h1:u/LLAOFgsMv7HmNL4Qufg58y+qElGOt5qv0z1mURkRY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_azidentity",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/azidentity",
    sum = "h1:vcYCAze6p19qBW7MhZybIsqD8sMV8js0NyQM8JDnVtg=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_azcore",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/azcore",
    sum = "h1:8kDqDngH+DmVBiCtIjCFTGa7MBnsIOkF9IccInFEbjk=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_azuread_microsoft_authentication_library_for_go",
    importpath = "github.com/AzureAD/microsoft-authentication-library-for-go",
    sum = "h1:OBhqkivkhkMqLPymWEppkm7vgPQY2XsHoEkaMQ0AdZY=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_internal",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/internal",
    sum = "h1:sXr+ck84g/ZlZUOZiNELInmMgOsuGwdjjVkEIde0OtY=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_kylelemons_godebug",
    importpath = "github.com/kylelemons/godebug",
    sum = "h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_pkg_browser",
    importpath = "github.com/pkg/browser",
    sum = "h1:KoWmjvw+nsYOo29YJK9vDA65RGE3NrOnUtO7a+RF9HU=",
    version = "v0.0.0-20210911075715-681adbf594b8",
)

go_repository(
    name = "com_github_golang_jwt_jwt_v4",
    importpath = "github.com/golang-jwt/jwt/v4",
    sum = "h1:7cYmW1XlMY7h7ii7UhUyChSgS5wUJEnm9uZVTGqOWzg=",
    version = "v4.5.0",
)

go_repository(
    name = "com_github_hashicorp_go_plugin",
    importpath = "github.com/hashicorp/go-plugin",
    sum = "h1:g6Lj3USwF5LaB8HlvCxPjN2X4nFE08ko2BJNVpl7TIE=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_mitchellh_go_testing_interface",
    importpath = "github.com/mitchellh/go-testing-interface",
    sum = "h1:jrgshOhYAUVNMAJiKbEu7EqAwgJJ2JqpQmpLJOu07cU=",
    version = "v1.14.1",
)

go_repository(
    name = "com_github_oklog_run",
    importpath = "github.com/oklog/run",
    sum = "h1:GEenZ1cK0+q0+wsJew9qUg/DyD8k3JzYsZAi5gYi2mA=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_hashicorp_yamux",
    importpath = "github.com/hashicorp/yamux",
    sum = "h1:yrQxtgseBDrq9Y652vSRDvsKCJKOUD+GzTS4Y0Y8pvE=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_googleapis_googleapis",
    importpath = "github.com/googleapis/googleapis",
    sum = "h1:H/JKoEVPUZjvazqVLKO2fdSFTudVmozSgOpeMJHe5n8=",
    version = "v0.0.0-20230907000724-75346295df28",
)

go_repository(
    name = "com_github_carolynvs_aferox",
    importpath = "github.com/carolynvs/aferox",
    sum = "h1:CMT50zX88amTMbFfFIWSTKRVRaOw6sejUMbbKiCD4zE=",
    version = "v0.3.0",
)

bazel_skylib_workspace()

load("//:repositories.bzl", "go_repositories")

# gazelle:proto disable_global

# Ref: https://github.com/bazelbuild/bazel/issues/10270
http_archive(
    name = "zlib",
    build_file = "@com_google_protobuf//:third_party/zlib.BUILD",
    sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
    strip_prefix = "zlib-1.2.11",
    urls = [
        "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
        "https://zlib.net/zlib-1.2.11.tar.gz",
    ],
)

go_repository(
    name = "com_github_spf13_cobra",
    importpath = "github.com/spf13/cobra",
    sum = "h1:hyqWnYt1ZQShIddO5kBpj3vu05/++x6tJ6dg8EC572I=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_spf13_viper",
    importpath = "github.com/spf13/viper",
    sum = "h1:js3yy885G8xwJa6iOISGFwd+qlUo5AvyXb7CiihdtiU=",
    version = "v1.15.0",
)

go_repository(
    name = "com_github_spf13_jwalterweatherman",
    importpath = "github.com/spf13/jwalterweatherman",
    sum = "h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_spf13_cast",
    importpath = "github.com/spf13/cast",
    sum = "h1:rj3WzYc11XZaIZMPKmwP96zkFEnnAmV8s6XbB2aY32w=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_spf13_pflag",
    importpath = "github.com/spf13/pflag",
    sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
    version = "v1.0.5",
)

go_repository(
    name = "com_github_subosito_gotenv",
    importpath = "github.com/subosito/gotenv",
    sum = "h1:X1TuBLAMDFbaTAChgCBLu3DU3UPyELpnF2jjJ2cz/S8=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_magiconair_properties",
    importpath = "github.com/magiconair/properties",
    sum = "h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=",
    version = "v1.8.7",
)

go_repository(
    name = "com_github_pelletier_go_toml_v2",
    importpath = "github.com/pelletier/go-toml/v2",
    sum = "h1:nrzqCb7j9cDFj2coyLNLaZuJTLjWjlaz6nvTvIwycIU=",
    version = "v2.0.6",
)

go_repository(
    name = "in_gopkg_ini_v1",
    importpath = "gopkg.in/ini.v1",
    sum = "h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=",
    version = "v1.67.0",
)

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

go_rules_dependencies()

go_register_toolchains(version = "1.20.5")

gazelle_dependencies()

# Ref: https://github.com/rules-proto-grpc/rules_proto_grpc/blob/master/docs/overriding_deps.rst
http_archive(
    name = "com_google_protobuf",
    sha256 = "da82be8acc5347c7918ef806ebbb621b24988f7e1a19b32cd7fc73bc29b59186",
    strip_prefix = "protobuf-3.25.3",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.25.3.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.25.3.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

http_archive(
    name = "rules_proto_grpc",
    sha256 = "2a0860a336ae836b54671cbbe0710eec17c64ef70c4c5a88ccfd47ea6e3739bd",
    strip_prefix = "rules_proto_grpc-4.6.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/releases/download/4.6.0/rules_proto_grpc-4.6.0.tar.gz"],
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_repos", "rules_proto_grpc_toolchains")

rules_proto_grpc_toolchains()

rules_proto_grpc_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

load("@rules_proto_grpc//:repositories.bzl", "bazel_gazelle", "io_bazel_rules_go")  # buildifier: disable=same-origin-load

io_bazel_rules_go()

bazel_gazelle()

load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")

rules_proto_grpc_go_repos()

protobuf_deps()

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
    name = "com_github_ash2k_bazel_tools",
    commit = "2add5bb84c2837a82a44b57e83c7414247aed43a",
    remote = "https://github.com/ash2k/bazel-tools.git",
)

load("@com_github_ash2k_bazel_tools//golangcilint:deps.bzl", "golangcilint_dependencies")

golangcilint_dependencies()

http_archive(
    name = "aspect_bazel_lib",
    sha256 = "04feedcd06f71d0497a81fdd3220140a373ff9d2bff94620fbd50b774f96d8e0",
    strip_prefix = "bazel-lib-1.40.2",
    url = "https://github.com/aspect-build/bazel-lib/releases/download/v1.40.2/bazel-lib-v1.40.2.tar.gz",
)

load("@aspect_bazel_lib//lib:repositories.bzl", "aspect_bazel_lib_dependencies")

aspect_bazel_lib_dependencies()

http_archive(
    name = "googleapis",
    sha256 = "88f13880b96a542e4a3695c02468bb3abd0ac6c665e9eee5931544779d495fc4",
    strip_prefix = "googleapis-694219fe5d75f264a7af30b9e3cff1a22dc2756f",
    urls = [
        "https://github.com/googleapis/googleapis/archive/694219fe5d75f264a7af30b9e3cff1a22dc2756f.zip",
    ],
)

load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
)
