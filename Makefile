CAVORITE_BIN := $(shell bazel cquery :cavorite --output=files 2>/dev/null)
LOCALSTORE_PLUGIN_BIN := $(shell bazel cquery plugins/localstore:localstore --output=files 2>/dev/null)

build: bazel_build

go_build: 
	go build

gazelle:
	bazel run :gazelle

bazel_build: gazelle
	bazel build :cavorite
	@echo Copy, paste, and execute this in your shell for convenience:
	@echo
	@echo CAVORITE_BIN=$(PWD)/$(CAVORITE_BIN)

bazel_build_docker:
	docker build --tag cavoritebazelbuild -f _ci/bazel_build/Dockerfile .
	docker run cavoritebazelbuild

localstore_plugin:
	bazel build plugins/localstore

with_localstore_plugin: bazel_build localstore_plugin
	@echo CAVORITE_BIN=$(PWD)/$(CAVORITE_BIN)
	@echo CAVORITE_PLUGIN=$(PWD)/$(LOCALSTORE_PLUGIN_BIN)

lint:
	docker build --tag cavoritelint -f _ci/lint/Dockerfile .
	docker run cavoritelint

test: gazelle
	bazel test //internal/...

minio:
	docker run -p 9000:9000 -p 9001:9001 quay.io/minio/minio server /data --console-address ":9001"

write_pluginproto:
	bazel run //internal/stores/pluginproto:write_generated_protos
