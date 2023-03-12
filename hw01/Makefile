DOCKER_IMG="astrviktor/health:1.0.0"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

docker-build:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG) \
		-f docker/Dockerfile .

docker-push:
	docker push $(DOCKER_IMG)

compose-up:
	docker-compose -f deployments/docker-compose.yml up --remove-orphans -d

compose-down:
	docker-compose -f deployments/docker-compose.yml down -v
