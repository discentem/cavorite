CAVORITE_BIN := $(shell bazel cquery :cavorite --output=files 2>/dev/null)

build: bazel_build

bazel_build_docker:
	docker build --tag cavoritebazelbuild -f _ci/bazel_build/Dockerfile .
	docker run cavoritebazelbuild

bazel_build: gazelle
	bazel build :cavorite
	@echo Copy, past, and execute this in your shell for convenience:
	@echo
	@echo CAVORITE_BIN=$(CAVORITE_BIN)

ci: compile_and_run lint

compile_and_run:
	docker build --tag cnr -f _ci/compile_and_run/Dockerfile .
	docker run cnr

lint:
	docker build --tag cavoritebutgolint -f _ci/lint/Dockerfile .
	docker run cavoritebutgolint

gazelle:
	bazel run :gazelle

minio:
	docker run -p 9000:9000 -p 9001:9001 quay.io/minio/minio server /data --console-address ":9001"
