# docker compose up and down
compose_up:
	@docker-compose -f ./deployments/compose/docker-compose.yml up --remove-orphans -d

compose_down:
	@docker-compose -f ./deployments/compose/docker-compose.yml down -v


# docker build orderservice image and push to dockerhub
DOCKER_IMG_ORDER="astrviktor/orderservice:1.0.2"
BINARY_NAME_ORDER=orderservice

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_order:
	cd ./apps/orderservice && \
	go build -o ${BINARY_NAME_ORDER} main.go && \
	cd ../..

docker_build_order:
	cd ./apps/orderservice && \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG_ORDER) \
		-f Dockerfile . && \
	cd ../..

docker_push_order:
	@docker push $(DOCKER_IMG_ORDER)


# docker build billing image and push to dockerhub
DOCKER_IMG_BILLING="astrviktor/billing:1.0.4"
BINARY_NAME_BILLING=billing

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_billing:
	cd ./apps/billing && \
	go build -o ${BINARY_NAME_BILLING} main.go && \
    cd ../..

docker_build_billing:
	cd ./apps/billing && \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG_BILLING) \
		-f Dockerfile . && \
	cd ../..

docker_push_billing:
	@docker push $(DOCKER_IMG_BILLING)


# docker build notification image and push to dockerhub
DOCKER_IMG_NOTIFICATION="astrviktor/notification:1.0.2"
BINARY_NAME_NOTIFICATION=notification

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_notification:
	cd ./apps/notification && \
	go build -o ${BINARY_NAME_NOTIFICATION} main.go && \
    cd ../..

docker_build_notification:
	cd ./apps/notification && \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG_NOTIFICATION) \
		-f Dockerfile . && \
	cd ../..

docker_push_notification:
	@docker push $(DOCKER_IMG_NOTIFICATION)


# helm install and uninstall
helm_install:
	@helm install storage ./deployments/kubernetes/helm-charts/postgres/ \
       --set auth.postgresPassword=password \
       --set auth.database=orders
	@sleep 60
	@helm install rabbitmq ./deployments/kubernetes/helm-charts/rabbitmq/
	@sleep 30
	@helm install orderservice ./deployments/kubernetes/helm-charts/orderservice/
	@sleep 10
	@helm install billing ./deployments/kubernetes/helm-charts/billing/
	@sleep 10
	@helm install notification ./deployments/kubernetes/helm-charts/notification/
	@sleep 120

helm_uninstall:
	@helm uninstall orderservice
	@helm uninstall billing
	@helm uninstall notification
	@helm uninstall rabbitmq
	@helm uninstall storage
	@sleep 10

# newman test postman collection
newman_run:
	@newman run ./postman/hw10.postman_collection.json --verbose