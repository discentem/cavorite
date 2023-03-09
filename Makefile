PANTRI_BIN := $(shell bazel cquery :pantri_but_go --output=files 2>/dev/null)

bazel_build: gazelle
	bazel build :pantri_but_go
	@echo Set pantri_bin for convenience:
	@echo
	@echo PANTRI_BIN=$(PANTRI_BIN)

ci: compile_and_run lint

compile_and_run:
	docker build --tag cnr -f _ci/compile_and_run/Dockerfile .
	docker run cnr

lint:
	docker build --tag pantributgolint -f _ci/lint/Dockerfile .
	docker run pantributgolint

gazelle:
	bazel run :gazelle

minio:
	docker run -p 9000:9000 -p 9001:9001 quay.io/minio/minio server /data --console-address ":9001"
