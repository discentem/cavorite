ci: compile_and_run lint

compile_and_run:
	docker build --tag pantri_compile_run -f _ci/compile_and_run/Dockerfile .
	docker run pantri_compile_run

lint:
	docker build --tag pantri_lint -f _ci/lint/Dockerfile .
	docker run pantri_lint

integration_tests:
	docker build --tag pantri_integration -f _ci/integration_tests/Dockerfile .
	docker run pantri_integration