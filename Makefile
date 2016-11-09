SHORT_NAME := steward-cf

include versioning.mk

REPO_PATH := github.com/deis/${SHORT_NAME}
DEV_ENV_IMAGE := quay.io/deis/go-dev:0.19.0
DEV_ENV_WORK_DIR := /go/src/${REPO_PATH}
K8S_CLAIMER_SUFFIX := -e K8S_CLAIMER_AUTH_TOKEN=${K8S_CLAIMER_AUTH_TOKEN}
DEV_ENV_PREFIX := docker run -it --rm -v ${CURDIR}:${DEV_ENV_WORK_DIR} -w ${DEV_ENV_WORK_DIR}
DEV_ENV_CMD := ${DEV_ENV_PREFIX} ${DEV_ENV_IMAGE}

BINARY_DEST_DIR := rootfs/bin

LDFLAGS := -ldflags '-X main.version=${VERSION}'

all:
	@echo "Use a Makefile to control top-level building of the project."

# Allow developers to step into the containerized development environment
dev:
	${DEV_ENV_CMD} bash

bootstrap:
	${DEV_ENV_CMD} glide install

glideup:
	${DEV_ENV_CMD} glide up

build:
	${DEV_ENV_CMD} sh -c \
	"GOOS=linux GOARCH=amd64 \
	go build \
	${LDFLAGS} \
	-o ${BINARY_DEST_DIR}/${SHORT_NAME} ."

test: test-unit

test-unit:
	${DEV_ENV_CMD} sh -c 'go test $$(glide nv)'

test-all:
	@${DEV_ENV_PREFIX} \
		${K8S_CLAIMER_SUFFIX} \
		${DEV_ENV_IMAGE} \
		sh -c 'go run testing/test_driver.go go test -tags integration $$(glide nv)'

test-cover:
	@${DEV_ENV_PREFIX} \
		${K8S_CLAIMER_SUFFIX} \
		${DEV_ENV_IMAGE} \
		sh -c 'go run testing/test_driver.go _scripts/test-cover.sh'

docker-build: build
	${DEV_ENV_CMD} upx -9 ${BINARY_DEST_DIR}/${SHORT_NAME}
	docker build --rm -t ${IMAGE} rootfs
	docker tag ${IMAGE} ${MUTABLE_IMAGE}

install-namespace:
	kubectl get ns steward || kubectl create -f manifests/steward-namespace.yaml

DEPLOY_IMAGE ?= quay.io/deisci/${SHORT_NAME}:devel

install:
ifndef BROKER_NAME
	$(error BROKER_NAME is undefined)
endif
ifndef BROKER_ACCESS_SCHEME
	$(error BROKER_ACCESS_SCHEME is undefined)
endif
ifndef BROKER_HOST
	$(error BROKER_HOST is undefined)
endif
ifndef BROKER_PORT
	$(error BROKER_PORT is undefined)
endif
ifndef BROKER_USERNAME
	$(error BROKER_USERNAME is undefined)
endif
ifndef BROKER_PASSWORD
	$(error BROKER_PASSWORD is undefined)
endif
	sed "s/#broker_name#/${BROKER_NAME}/g" manifests/${SHORT_NAME}-template.yaml > manifests/${BROKER_NAME}-${SHORT_NAME}.yaml
	sed -i.bak "s/#broker_access_scheme#/${BROKER_ACCESS_SCHEME}/g" manifests/${BROKER_NAME}-${SHORT_NAME}.yaml
	sed -i.bak "s/#broker_host#/${BROKER_HOST}/g" manifests/${BROKER_NAME}-${SHORT_NAME}.yaml
	sed -i.bak "s/#broker_port#/${BROKER_PORT}/g" manifests/${BROKER_NAME}-${SHORT_NAME}.yaml
	sed -i.bak "s/#broker_username#/${BROKER_USERNAME}/g" manifests/${BROKER_NAME}-${SHORT_NAME}.yaml
	sed -i.bak "s/#broker_password#/${BROKER_PASSWORD}/g" manifests/${BROKER_NAME}-${SHORT_NAME}.yaml
	sed -i.bak "s#\#deploy_image\##${DEPLOY_IMAGE}#g" manifests/${BROKER_NAME}-${SHORT_NAME}.yaml
	rm manifests/${BROKER_NAME}-${SHORT_NAME}.yaml.bak
	kubectl get deployment ${BROKER_NAME}-${SHORT_NAME} --namespace=steward && \
	kubectl apply -f manifests/${BROKER_NAME}-${SHORT_NAME}.yaml || \
	kubectl create -f manifests/${BROKER_NAME}-${SHORT_NAME}.yaml

deploy: install-namespace install

dev-deploy: docker-build docker-push
	DEPLOY_IMAGE=${IMAGE} $(MAKE) deploy
