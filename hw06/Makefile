DOCKER_IMG="astrviktor/auth:1.1.5"

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



# docker compose

compose-up:
	docker-compose -f docker-compose/docker-compose.yml up --remove-orphans -d

compose-down:
	docker-compose -f docker-compose/docker-compose.yml down -v

# helm install

helm_install_auth:
	@helm install auth ./helm-charts/auth/

# helm uninstall

helm_uninstall_auth:
	@helm uninstall auth


# newman test postman collection

newman_run:
	@newman run ./postman/hw06-auth.postman_collection.json --verbose