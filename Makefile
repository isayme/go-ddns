.PHONY: dev
dev:
	CONF_FILE_PATH=${PWD}/config/dev.json go run main.go

APP_NAME := ddns
APP_VERSION := $(shell git describe --tags --always)
APP_PKG := $(shell echo ${PWD} | sed -e "s\#${GOPATH}/src/\#\#g")
APP_BUILDTIME := $(shell date -u +"%FT%TZ")
APP_GITSHA1 := $(shell git rev-parse HEAD)

.PHONY: image
image:
	docker build \
	--build-arg APP_NAME=${APP_NAME} \
	--build-arg APP_VERSION=${APP_VERSION} \
	--build-arg APP_BUILDTIME=${APP_BUILDTIME} \
	--build-arg APP_GITSHA1=${APP_GITSHA1} \
	--build-arg APP_PKG=${APP_PKG} \
	-t ${APP_NAME}:${APP_VERSION} .

.PHONY: publish
publish: image
	docker tag ${APP_NAME}:${APP_VERSION} isayme/${APP_NAME}:${APP_VERSION}
	docker push isayme/${APP_NAME}:${APP_VERSION}
	docker tag ${APP_NAME}:${APP_VERSION} isayme/${APP_NAME}:latest
	docker push isayme/${APP_NAME}:latest
