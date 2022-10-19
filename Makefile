ci: compile_and_run lint

compile_and_run:
	docker build --tag cnr -f _ci/compile_and_run/Dockerfile .
	docker run cnr

lint:
	docker build --tag pantributgolint -f _ci/lint/Dockerfile .
	docker run pantributgolint