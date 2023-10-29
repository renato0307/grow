IMG=grow-service-ingestion

.PHONY: docker-build
docker-build:
	docker build -t ${IMG} .

.PHONY: docker-run
docker-run:
	docker run -t ${IMG} ${ARGS}