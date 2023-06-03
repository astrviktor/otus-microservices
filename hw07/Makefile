# docker compose up and down
compose_up:
	@docker-compose -f ./deployments/compose/docker-compose.yml up --remove-orphans -d

compose_down:
	@docker-compose -f ./deployments/compose/docker-compose.yml down -v


# docker build image and push to dockerhub
DOCKER_IMG="astrviktor/orderservice:1.0.0"
BINARY_NAME=orderservice

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	@go build -o ${BINARY_NAME} cmd/main.go

docker_build:
	@docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG) \
		-f docker/Dockerfile .

docker_push:
	@docker push $(DOCKER_IMG)


# helm install and uninstall
helm_install:
	@helm install storage ./deployments/kubernetes/helm-charts/postgres/ \
       --set auth.postgresPassword=password \
       --set auth.database=orders
	@sleep 60
	@helm install orderservice ./deployments/kubernetes/helm-charts/orderservice/

helm_uninstall:
	@helm uninstall orderservice
	@sleep 10
	@helm uninstall storage


# newman test postman collection
newman_run:
	@newman run ./postman/hw07-orders.postman_collection.json --verbose