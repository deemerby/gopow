PROJECT_NAME          := go-pow
DOCKERFILE_PATH       := $(CURDIR)/docker

# configuration for image names
IMAGE_VERSION         ?= latest
IMAGE_REGISTRY        ?= dlahuta

# configuration for server binary and image
SERVER_IMAGE          := $(IMAGE_REGISTRY)/$(PROJECT_NAME)
SERVER_DOCKERFILE     := $(DOCKERFILE_PATH)/Dockerfile

.PHONY clean:
clean:
	@docker rmi -f $(shell docker images -q $(SERVER_IMAGE)) || true
	rm .docker-*


.docker-$(IMAGE_NAME)-$(IMAGE_VERSION):
	@docker build --platform linux/amd64 -f $(SERVER_DOCKERFILE) -t $(SERVER_IMAGE):$(IMAGE_VERSION) .
	@docker image prune -f --filter label=stage=server-intermediate
	touch $@

.PHONY: docker
docker: .docker-$(IMAGE_NAME)-$(IMAGE_VERSION)

.push-$(IMAGE_NAME)-$(IMAGE_VERSION):
ifndef IMAGE_REGISTRY
	@(echo "Please set IMAGE_REGISTRY variable in Makefile to use push command"; exit 1)
else
	@docker push $(SERVER_IMAGE):$(IMAGE_VERSION)
endif
	touch $@

.PHONY: push
push: .push-$(IMAGE_NAME)-$(IMAGE_VERSION)

.PHONY vendor: 
	@go mod tidy
	@go mod download
	@go mod vendor

client: 
	@go run ./cmd/ --type.server=false

server: 
	@go run ./cmd/
