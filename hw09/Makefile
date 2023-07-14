# docker compose up and down
compose_up:
	@docker-compose -f ./deployments/compose/docker-compose.yml up --remove-orphans -d

compose_down:
	@docker-compose -f ./deployments/compose/docker-compose.yml down -v


# docker build orderservice image and push to dockerhub
DOCKER_IMG_ORDER="astrviktor/orderservice:1.0.1"
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


# docker build payment image and push to dockerhub
DOCKER_IMG_PAYMENT="astrviktor/payment:1.0.1"
BINARY_NAME_PAYMENT=payment

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_payment:
	cd ./apps/payment && \
	go build -o ${BINARY_NAME_PAYMENT} main.go && \
    cd ../..

docker_build_payment:
	cd ./apps/payment && \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG_PAYMENT) \
		-f Dockerfile . && \
	cd ../..

docker_push_payment:
	@docker push $(DOCKER_IMG_PAYMENT)


# docker build warehouse image and push to dockerhub
DOCKER_IMG_WAREHOUSE="astrviktor/warehouse:1.0.1"
BINARY_NAME_WAREHOUSE=warehouse

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_warehouse:
	cd ./apps/warehouse && \
	go build -o ${BINARY_NAME_WAREHOUSE} main.go && \
    cd ../..

docker_build_warehouse:
	cd ./apps/warehouse && \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG_WAREHOUSE) \
		-f Dockerfile . && \
	cd ../..

docker_push_warehouse:
	@docker push $(DOCKER_IMG_WAREHOUSE)


# docker build delivery image and push to dockerhub
DOCKER_IMG_DELIVERY="astrviktor/delivery:1.0.1"
BINARY_NAME_DELIVERY=delivery

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_delivery:
	cd ./apps/delivery && \
	go build -o ${BINARY_NAME_DELIVERY} main.go && \
    cd ../..

docker_build_delivery:
	cd ./apps/delivery && \
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG_DELIVERY) \
		-f Dockerfile . && \
	cd ../..

docker_push_delivery:
	@docker push $(DOCKER_IMG_DELIVERY)


# helm install and uninstall
helm_install:
	@helm install storage ./deployments/kubernetes/helm-charts/postgres/ \
       --set auth.postgresPassword=password \
       --set auth.database=orders
	@sleep 60
	@helm install orderservice ./deployments/kubernetes/helm-charts/orderservice/
	@helm install payment ./deployments/kubernetes/helm-charts/payment/
	@helm install warehouse ./deployments/kubernetes/helm-charts/warehouse/
	@helm install delivery ./deployments/kubernetes/helm-charts/delivery/

helm_uninstall:
	@helm uninstall orderservice
	@helm uninstall payment
	@helm uninstall warehouse
	@helm uninstall delivery
	@helm uninstall storage
	@sleep 10

# newman test postman collection
newman_run:
	@newman run ./postman/saga.postman_collection.json --verbose