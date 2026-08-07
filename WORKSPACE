load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "86d3dc8f59d253524f933aaf2f3c05896cb0b605fc35b460c0b4b039996124c6",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.60.0/rules_go-v0.60.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_resourcemanager_storage_armstorage",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storage/armstorage",
    sum = "h1:AifHbc4mg0x9zW52WOpKbsHaDKuRhlI7TVl47thgQ70=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_storage_azblob_testdata_perf",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/testdata/perf",
    sum = "h1:45Ajiuhu6AeJTFdwxn2OWXZTQOHdXT1U/aezrVu6HIM=",
    version = "v0.0.0-20240208231215-981108a6de20",
)

go_repository(
    name = "com_github_bufbuild_protocompile",
    importpath = "github.com/bufbuild/protocompile",
    sum = "h1:LbFKd2XowZvQ/kajzguUp2DC9UEIQhIq77fZZlaQsNA=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_dnaeon_go_vcr",
    importpath = "github.com/dnaeon/go-vcr",
    sum = "h1:zHCHvJYTMh1N7xnV7zf1m1GPBF9Ad0Jk/whtQ1663qI=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    importpath = "github.com/dustin/go-humanize",
    sum = "h1:GzkhY7T5VNhEkwH0PVJgjz+fX1rhBrR7pRT3mDkpeCY=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_fsouza_slognil",
    importpath = "github.com/fsouza/slognil",
    sum = "h1:O17M/Ro93X5JVjMYVgObkOoWLYz616ZMuzjAy0FLQag=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_go_logr_logr",
    importpath = "github.com/go-logr/logr",
    sum = "h1:pKouT5E8xu9zeFC39JXRDukb6JFQPXM5p5I91188VAQ=",
    version = "v1.4.1",
)

go_repository(
    name = "com_github_go_logr_stdr",
    importpath = "github.com/go-logr/stdr",
    sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_golang_jwt_jwt_v5",
    importpath = "github.com/golang-jwt/jwt/v5",
    sum = "h1:d/ix8ftRUorsN+5eMIlF4T6J8CAt9rch3My2winC1Jw=",
    version = "v5.2.0",
)

go_repository(
    name = "com_github_golang_snappy",
    importpath = "github.com/golang/snappy",
    sum = "h1:yAGX7huGHXlcLOEtBnF4w7FQwA26wojNCwOYAEhLjQM=",
    version = "v0.0.4",
)

go_repository(
    name = "com_github_google_go_pkcs11",
    importpath = "github.com/google/go-pkcs11",
    sum = "h1:OF1IPgv+F4NmqmJ98KTjdN97Vs1JxDPB3vbmYzV2dpk=",
    version = "v0.2.1-0.20230907215043-c6f79328ddf9",
)

go_repository(
    name = "com_github_google_s2a_go",
    importpath = "github.com/google/s2a-go",
    sum = "h1:60BLSyTrOV4/haCDW4zb1guZItoSq8foHCXrAnjBo/o=",
    version = "v0.1.7",
)

go_repository(
    name = "com_github_jhump_protoreflect",
    importpath = "github.com/jhump/protoreflect",
    sum = "h1:HUMERORf3I3ZdX05WaQ6MIpd/NJ434hTp5YiKgfCL6c=",
    version = "v1.15.1",
)

go_repository(
    name = "com_github_klauspost_compress",
    importpath = "github.com/klauspost/compress",
    sum = "h1:Ej5ixsIri7BrIjBkRZLTo6ghwrEtHFk7ijlczPW4fZ4=",
    version = "v1.17.4",
)

go_repository(
    name = "com_github_klauspost_cpuid_v2",
    importpath = "github.com/klauspost/cpuid/v2",
    sum = "h1:ndNyv040zDGIDh8thGkXYjnFtiN02M1PVVF+JE/48xc=",
    version = "v2.2.6",
)

go_repository(
    name = "com_github_minio_md5_simd",
    importpath = "github.com/minio/md5-simd",
    sum = "h1:Gdi1DZK69+ZVMoNHRXJyNcxrMA4dSxoYHZSQbirFg34=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_minio_minio_go_v7",
    importpath = "github.com/minio/minio-go/v7",
    sum = "h1:bnTOXOHjOqv/gcMuiVbN9o2ngRItvqE774dG9nq0Dzw=",
    version = "v7.0.66",
)

go_repository(
    name = "com_github_minio_sha256_simd",
    importpath = "github.com/minio/sha256-simd",
    sum = "h1:6kaan5IFmwTNynnKKpDHe6FWHohJOHhCPchzK49dzMM=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_montanaflynn_stats",
    importpath = "github.com/montanaflynn/stats",
    sum = "h1:r3y12KyNxj/Sb/iOE46ws+3mS1+MZca1wlHQFPsY/JU=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_nats_io_nats_go",
    importpath = "github.com/nats-io/nats.go",
    sum = "h1:/WFBHEc/dOKBF6qf1TZhrdEfTmOZ5JzdJ+Y3m6Y/p7E=",
    version = "v1.31.0",
)

go_repository(
    name = "com_github_nats_io_nkeys",
    importpath = "github.com/nats-io/nkeys",
    sum = "h1:IzVe95ru2CT6ta874rt9saQRkWfe2nFj1NtvYSLqMzY=",
    version = "v0.4.6",
)

go_repository(
    name = "com_github_nats_io_nuid",
    importpath = "github.com/nats-io/nuid",
    sum = "h1:5iA8DT8V7q8WK2EScv2padNa/rTESc1KdnPw4TC2paw=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_rs_xid",
    importpath = "github.com/rs/xid",
    sum = "h1:mKX4bl4iPYJtEIxp6CYiUuLQ/8DYMoz0PUdtGgMFRVc=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_sagikazarmark_locafero",
    importpath = "github.com/sagikazarmark/locafero",
    sum = "h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_sagikazarmark_slog_shim",
    importpath = "github.com/sagikazarmark/slog-shim",
    sum = "h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_google_cloud_go_apigeeregistry",
    importpath = "cloud.google.com/go/apigeeregistry",
    sum = "h1:C+QU2K+DzDjk4g074ouwHQGkoff1h5OMQp6sblCVreQ=",
    version = "v0.8.3",
)

go_repository(
    name = "com_google_cloud_go_dataproc_v2",
    importpath = "cloud.google.com/go/dataproc/v2",
    sum = "h1:/u81Fd+BvCLp+xjctI1DiWVJn6cn9/s3Akc8xPH02yk=",
    version = "v2.4.0",
)

go_repository(
    name = "io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc",
    importpath = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
    sum = "h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=",
    version = "v0.49.0",
)

go_repository(
    name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
    importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
    sum = "h1:jq9TW8u3so/bN+JPT166wjOI6/vQPF6Xe7nMNIltagk=",
    version = "v0.49.0",
)

go_repository(
    name = "io_opentelemetry_go_otel",
    importpath = "go.opentelemetry.io/otel",
    sum = "h1:0LAOdjNmQeSTzGBzduGe/rU4tZhMwL5rWgtp9Ku5Jfo=",
    version = "v1.24.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_metric",
    importpath = "go.opentelemetry.io/otel/metric",
    sum = "h1:6EhoGWWK28x1fbpA4tYTOWBkPefTDQnb8WSGXlc88kI=",
    version = "v1.24.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_sdk",
    importpath = "go.opentelemetry.io/otel/sdk",
    sum = "h1:FTt8qirL1EysG6sTQRZ5TokkU8d0ugCj8htOgThZXQ8=",
    version = "v1.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_trace",
    importpath = "go.opentelemetry.io/otel/trace",
    sum = "h1:CsKnnL4dUAr/0llH9FKuc698G04IrpWV0MQA/Y1YELI=",
    version = "v1.24.0",
)

go_repository(
    name = "org_golang_google_genproto_googleapis_api",
    importpath = "google.golang.org/genproto/googleapis/api",
    sum = "h1:9g7erC9qu44ks7UK4gDNlnk4kOxZG707xKm4jVniy6o=",
    version = "v0.0.0-20240221002015-b0ce06bbee7c",
)

go_repository(
    name = "org_golang_google_genproto_googleapis_bytestream",
    importpath = "google.golang.org/genproto/googleapis/bytestream",
    sum = "h1:YqFWYZXim8bG9v68xU8WjTZmYKb5M5dMeSOWIp6jogI=",
    version = "v0.0.0-20240304161311-37d4d3c04a78",
)

go_repository(
    name = "org_golang_google_genproto_googleapis_rpc",
    importpath = "google.golang.org/genproto/googleapis/rpc",
    sum = "h1:Xs9lu+tLXxLIfuci70nG4cpwaRC+mRQPUL7LoIeDJC4=",
    version = "v0.0.0-20240304161311-37d4d3c04a78",
)

go_repository(
    name = "tech_einride_go_aip",
    importpath = "go.einride.tech/aip",
    sum = "h1:XfV+NQX6L7EOYK11yoHHFtndeaWh3KbD9/cN/6iWEt8=",
    version = "v0.66.0",
)

go_rules_dependencies()

http_archive(
    name = "bazel_gazelle",
    sha256 = "de86bb09ddac514d9831cd0af5addb49f196eff66cf7c779a595b410cbe40253",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.51.0/bazel-gazelle-v0.51.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.51.0/bazel-gazelle-v0.51.0.tar.gz",
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
    sha256 = "3b5b49006181f5f8ff626ef8ddceaa95e9bb8ad294f7b5d7b11ea9f7ddaf8c59",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.9.0/bazel-skylib-1.9.0.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.9.0/bazel-skylib-1.9.0.tar.gz",
    ],
)

load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

go_repository(
    name = "com_github_fsouza_fake_gcs_server",
    importpath = "github.com/fsouza/fake-gcs-server",
    sum = "h1:i/rV62ZOh/3y2aBl4jaGxIf0sySpVnTaot54r4BpgaE=",
    version = "v1.47.8",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:dueUQJ1C2q9oE3F7wvmSGAaVtTmUizReu6fjN8uqzbQ=",
    version = "v1.9.3",
)

go_repository(
    name = "com_github_gorilla_handlers",
    importpath = "github.com/gorilla/handlers",
    sum = "h1:cLTUSsNkgcwhgRqvCNmdbRWG0A3N4F+M2nWKdScwyEE=",
    version = "v1.5.2",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:TuBL49tXwgrFYWhqrNgrUNEY92u81SPhu7sTdzQEiWY=",
    version = "v1.8.1",
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
    sum = "h1:NFTV2Zj1bL4mc9sqWACXbQFVBBg2W3GPvqp8/ESS2Wg=",
    version = "v1.0.4",
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
    sum = "h1:hLrqtEDnRye3+sgx6z4qVLNuviH3MR5aQ0ykNJa/UYA=",
    version = "v1.0.0",
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
    sum = "h1:fXPMAmuh0gDuRDey0atC8cXBuKIlqCzCkL8sm1n9Ov0=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_azidentity",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/azidentity",
    sum = "h1:sO0/P7g68FrryJzljemN+6GTssUXdANk6aJ7T1ZxnsQ=",
    version = "v1.5.1",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_azcore",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/azcore",
    sum = "h1:c4k2FIYIh4xtwqrQwV0Ct1v5+ehlNXj5NI/MWVsiTkQ=",
    version = "v1.9.2",
)

go_repository(
    name = "com_github_azuread_microsoft_authentication_library_for_go",
    importpath = "github.com/AzureAD/microsoft-authentication-library-for-go",
    sum = "h1:DzHpqpoJVaCgOUdVHxE8QB52S6NiVdDQvGlny1qvPqA=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_azure_azure_sdk_for_go_sdk_internal",
    importpath = "github.com/Azure/azure-sdk-for-go/sdk/internal",
    sum = "h1:LqbJ/WzJUwBf8UiaSzgX7aMclParm9/5Vgp+TY51uBQ=",
    version = "v1.5.2",
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
    sum = "h1:+mdjkGKdHQG3305AYmdv1U2eRNDiU2ErMBj1gwrq8eQ=",
    version = "v0.0.0-20240102092130-5ac0b6a4141c",
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
    sum = "h1:7aJaZx1B85qltLMc546zn58BxxfZdR/W22ej9CFoEf0=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_spf13_viper",
    importpath = "github.com/spf13/viper",
    sum = "h1:LUXCnvUvSM6FXAsj6nnfc8Q2tp1dIgUfY9Kc8GsSOiQ=",
    version = "v1.18.2",
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
    sum = "h1:GEiTHELF+vaR5dhz3VqZfFSzZjYbgeKDpBxQVS4GYJ0=",
    version = "v1.6.0",
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
    sum = "h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=",
    version = "v1.6.0",
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
    sum = "h1:FnwAJ4oYMvbT/34k9zzHuZNrhlz48GB3/s6at6/MHO4=",
    version = "v2.1.0",
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
    sha256 = "ccff8964efdc4052f0b3579ad503dba28729c28fb0cf4245c060ec17667666aa",
    strip_prefix = "protobuf-3.29.5",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.29.5.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.29.5.tar.gz",
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
    commit = "415483a9e13342a6603a710b0296f6d85b8d26bf",
    remote = "https://github.com/ash2k/bazel-tools.git",
)

load("@com_github_ash2k_bazel_tools//golangcilint:deps.bzl", "golangcilint_dependencies")

golangcilint_dependencies()

http_archive(
    name = "aspect_bazel_lib",
    sha256 = "94e192033ca8027f26de71c9000a67ef9c73695c2b88e2c559045170917ead0c",
    strip_prefix = "bazel-lib-2.22.5",
    url = "https://github.com/aspect-build/bazel-lib/releases/download/v2.22.5/bazel-lib-v2.22.5.tar.gz",
)

load("@aspect_bazel_lib//lib:repositories.bzl", "aspect_bazel_lib_dependencies")

aspect_bazel_lib_dependencies()

http_archive(
    name = "googleapis",
    sha256 = "f1b33bb68bd4e3362e5147b4ebe09bbf97de9c9b9afadf752b6b62c4b73406a4",
    strip_prefix = "googleapis-0054cdcbd8ea44298f329d916d8173dd736fbaaa",
    urls = [
        "https://github.com/googleapis/googleapis/archive/0054cdcbd8ea44298f329d916d8173dd736fbaaa.zip",
    ],
)

load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
)
