# auto generate server by openapi file

# go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4
# oapi-codegen --version

generate_openapi_server:
	@rm -rf ./internal/server/*
	@oapi-codegen -package user -generate chi-server -o ./internal/server/server.go ./api/user-openapi.yaml


generate_openapi_types:
	@rm -rf ./internal/types/*
	@oapi-codegen -package user -generate types -o ./internal/types/types.go ./api/user-openapi.yaml


# convert openapi to swagger

# https://github.com/LucyBot-Inc/api-spec-converter
# sudo npm install -g api-spec-converter
# api-spec-converter -h

convert_openapi_to_swagger:
	@rm -rf ./api/user-swagger.yaml
	@api-spec-converter --from=openapi_3 --to=swagger_2 --syntax=yaml ./api/user-openapi.yaml > ./api/user-swagger.yaml


# auto generate server by swagger file

# https://goswagger.io/install.html

swagger_server:
	@rm -rf ./internal/server/*
	@swagger generate server -f ./api/user-swagger.yaml -s ./internal/server/api -m ./internal/server/models


# docker compose up and down
compose_up:
	@docker-compose -f ./deployments/compose/docker-compose.yml up --remove-orphans -d

compose_down:
	@docker-compose -f ./deployments/compose/docker-compose.yml down -v


# docker build image and push to dockerhub
DOCKER_IMG="astrviktor/crudservice:1.0.2"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

docker_build:
	@docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		--platform linux/amd64 \
		-t $(DOCKER_IMG) \
		-f docker/Dockerfile .

docker_push:
	@docker push $(DOCKER_IMG)

