PATH_SNOWY = github.com/trussle/snowy

UNAME_S := $(shell uname -s)
SED ?= sed -i
ifeq ($(UNAME_S),Darwin)
	SED += '' --
endif

.PHONY: all
all: install
	$(MAKE) clean build

.PHONY: install
install:
	go get github.com/Masterminds/glide
	go get github.com/mjibson/esc
	go get github.com/mattn/goveralls
	glide install

.PHONY: build
build: dist/documents \
	dist/prompt

dist/documents:
	go build -o dist/documents ${PATH_SNOWY}/cmd/documents

dist/prompt:
	go build -o dist/prompt ${PATH_SNOWY}/cmd/prompt

pkg/store/mocks/store.go:
	mockgen -package=mocks -destination=pkg/store/mocks/store.go ${PATH_SNOWY}/pkg/store Store
	@ $(SED) 's/github.com\/trussle\/snowy\/vendor\///g' ./pkg/store/mocks/store.go

pkg/repository/mocks/repository.go:
	mockgen -package=mocks -destination=pkg/repository/mocks/repository.go ${PATH_SNOWY}/pkg/repository Repository
	@ $(SED) 's/github.com\/trussle\/snowy\/vendor\///g' ./pkg/repository/mocks/repository.go

pkg/metrics/mocks/metrics.go:
	mockgen -package=mocks -destination=pkg/metrics/mocks/metrics.go ${PATH_SNOWY}/pkg/metrics Gauge,HistogramVec,Counter
	@ $(SED) 's/github.com\/trussle\/snowy\/vendor\///g' ./pkg/metrics/mocks/metrics.go

pkg/metrics/mocks/observer.go:
	mockgen -package=mocks -destination=pkg/metrics/mocks/observer.go github.com/prometheus/client_golang/prometheus Observer
	@ $(SED) 's/github.com\/trussle\/snowy\/vendor\///g' ./pkg/metrics/mocks/observer.go

.PHONY: build-mocks
build-mocks: pkg/store/mocks/store.go \
	pkg/repository/mocks/repository.go \
	pkg/metrics/mocks/metrics.go \
	pkg/metrics/mocks/observer.go

.PHONY: clean
clean: FORCE
	rm -f dist/documents
	rm -f dist/prompt

.PHONY: clean-mocks
clean-mocks: FORCE
	rm -f pkg/store/mocks/store.go
	rm -f pkg/repository/mocks/repository.go
	rm -f pkg/metrics/mocks/metrics.go
	rm -f pkg/metrics/mocks/observer.go

FORCE:

.PHONY: integration-tests
integration-tests:
	docker-compose run documents go test -v -tags=integration ./cmd/... ./pkg/...

.PHONY: unit-tests
unit-tests:
	docker-compose run documents go test -v ./cmd/... ./pkg/...

.PHONY: documentation
documentation:
	go test -v -tags=documentation ./pkg/... -run=TestDocumentation_

.PHONY: coverage-integration-tests
coverage-integration-tests:
	docker-compose run documents go test -covermode=count -coverprofile=bin/integration/coverage.out -v -tags=integration ${COVER_PKG}

.PHONY: coverage-unit-tests
coverage-unit-tests:
	docker-compose run documents go test -covermode=count -coverprofile=bin/unit/coverage.out -v ${COVER_PKG}

.PHONY: coverage-view-integration
coverage-view-integration:
	go tool cover -html=bin/integration/coverage.out

.PHONY: coverage-view-unit
coverage-view-unit:
	go tool cover -html=bin/unit/coverage.out

.PHONY: build-ui
build-ui: ui/scripts/snowy.js pkg/ui/static.go

.PHONY: clean-ui
clean-ui: FORCE
	@ rm -f pkg/ui/static.go
	$(MAKE) -C ./ui clean

ui/scripts/snowy.js:
	$(MAKE) -C ./ui

pkg/ui/static.go:
	esc -o="pkg/ui/static.go" -ignore="elm-stuff|Makefile|src|elm-package.json" -pkg="ui" ui

PWD ?= ${GOPATH}/src/${PATH_SNOWY}
TAG ?= dev
BRANCH ?= dev
ifeq ($(BRANCH),master)
	TAG=latest
endif

.PHONY: build-docker
build-docker:
	@echo "Building '${TAG}' for '${BRANCH}'"
	docker run --rm -v ${PWD}:/go/src/${PATH_SNOWY} -w /go/src/${PATH_SNOWY} iron/go:dev go build -o documents ${PATH_SNOWY}/cmd/documents
	docker build -t teamtrussle/snowy:${TAG} .

.PHONY: push-docker-tag
push-docker-tag: FORCE
	@echo "Pushing '${TAG}' for '${BRANCH}'"
	docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
	docker push teamtrussle/snowy:${TAG}

.PHONY: push-docker
ifeq ($(TAG),latest)
push-docker: FORCE
	@echo "Pushing '${TAG}' for '${BRANCH}'"
	docker login -u ${DOCKER_HUB_USERNAME} -p ${DOCKER_HUB_PASSWORD}
	docker push teamtrussle/snowy:${TAG}
else
push-docker: FORCE
	@echo "Pushing requires branch '${BRANCH}' to be master"
endif
