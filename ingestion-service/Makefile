IMG=grow-ingestion-service

.PHONY: docker-build
docker-build:
	docker build -t ${IMG} .

.PHONY: docker-run
docker-run:
	docker run -t ${IMG} ${ARGS}

.PHONY: local-run
local-run:
	go run . --log-level=debug