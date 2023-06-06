load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "56d8c5a5c91e1af73eca71a6fab2ced959b67c86d12ba37feedb0a2dfea441a6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.37.0/rules_go-v0.37.0.zip",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "ecba0f04f96b4960a5b250c8e8eeec42281035970aa8852dda73098274d14a1d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.29.0/bazel-gazelle-v0.29.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

############################################################
# Define your own dependencies here using go_repository.
# Else, dependencies declared by rules_go/gazelle will be used.
# The first declaration of an external repository "wins".
############################################################

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "bazel_skylib",
    sha256 = "b8a1527901774180afc798aeb28c4634bdccf19c4d98e7bdd1ce79d1fe9aaad7",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.4.1/bazel-skylib-1.4.1.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.4.1/bazel-skylib-1.4.1.tar.gz",
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
    name = "com_github_schollz_progressbar_v3",
    importpath = "github.com/schollz/progressbar/v3",
    sum = "h1:o8rySDYiQ59Mwzy2FELeHY5ZARXZTVJC7iHD6PEFUiE=",
    version = "v3.13.1",
)

go_repository(
    name = "com_github_mattn_go_runewidth",
    importpath = "github.com/mattn/go-runewidth",
    sum = "h1:+xnbZSEeDbOIg5/mE6JF0w6n9duR1l3/WmbinWVwUuU=",
    version = "v0.0.14",
)

go_repository(
    name = "com_github_mitchellh_colorstring",
    importpath = "github.com/mitchellh/colorstring",
    sum = "h1:62I3jR2EmQ4l5rM/4FEfDWcRD+abF5XlKShorW5LRoQ=",
    version = "v0.0.0-20190213212951-d06e56a500db",
)

go_repository(
    name = "com_github_rivo_uniseg",
    importpath = "github.com/rivo/uniseg",
    sum = "h1:8TfxU8dW6PdqD27gjM8MVNuicgxIjxpm4K7x4jp8sis=",
    version = "v0.4.4",
)

bazel_skylib_workspace()

load("//:repositories.bzl", "go_repositories")

# gazelle:proto disable_global

# Ref: https://github.com/rules-proto-grpc/rules_proto_grpc/blob/master/docs/overriding_deps.rst
http_archive(
    name = "com_google_protobuf",
    sha256 = "8b28fdd45bab62d15db232ec404248901842e5340299a57765e48abe8a80d930",
    strip_prefix = "protobuf-3.20.1",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.20.1.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.20.1.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

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

go_register_toolchains(version = "1.19.5")

gazelle_dependencies()
