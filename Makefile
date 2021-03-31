default: build

.PHONY: build

APP_NAME="grid650-array-serial"
git_revision=git rev-parse --short HEAD
BUILD_TIME=`date +%FT%T%z`
BUILD_DATE=`date +%F`
GIT_REVISION=`${git_revision}`
GIT_COMMIT=`git rev-parse HEAD`
GIT_BRANCH=`git rev-parse --symbolic-full-name --abbrev-ref HEAD`
GIT_DIRTY=`test -z "$$(git status --porcelain)" && echo "clean" || echo "dirty"`
VERSION=`git describe --tag --abbrev=0 --exact-match HEAD 2> /dev/null || (echo 'Git tag not found, fallback to commit id' >&2; ${git_revision})`
GO_VERSION?=$(shell echo '$$(go version |cut -d " " -f 3 | cut -b 3-)')
GO_PKGS=$(shell go list ./... | grep -v '/vendor/' | grep -v 'machinery' | grep -v 'test' | grep -v 'mocks')
BUILD_META_PKG_PATH=github.com/XSAM/go-hybrid
INJECT_VARIABLE=-X ${BUILD_META_PKG_PATH}/metadata.gitVersion=${VERSION} -X ${BUILD_META_PKG_PATH}/metadata.gitCommit=${GIT_COMMIT} -X ${BUILD_META_PKG_PATH}/metadata.gitBranch=${GIT_BRANCH} -X ${BUILD_META_PKG_PATH}/metadata.gitTreeState=${GIT_DIRTY} -X ${BUILD_META_PKG_PATH}/metadata.buildTime=${BUILD_TIME}

FLAGS=-trimpath -ldflags "-s ${INJECT_VARIABLE}"
DEBUG_FLAGS=-gcflags "all=-N -l" ${FLAGS}

export CGO_ENABLED=1

unit-test:
	@echo "> Running Unit Test..."
	go test --tags=unittest ./...

integration-test:
	@echo "> Running Integration Test..."
	go test --tags=integration ./...

build:
	@echo "> Building binaries for current os"
	go build -o bin/${APP_NAME} ${FLAGS} cmd/main.go

build-linux:
	@echo "> Building linux dist binaries..."
	GOARCH=amd64 GOOS=linux go build -o bin/${APP_NAME} ${FLAGS} cmd/main.go

build-windows:
	@echo "> Building windows dist binaries..."
	GOARCH=amd64 GOOS=windows go build -o bin/${APP_NAME} ${FLAGS} cmd/main.go

build-darwin:
	@echo "> Building darwin dist binaries..."
	GOARCH=amd64 GOOS=darwin go build -o bin/${APP_NAME} ${FLAGS} cmd/main.go

debug-build:
	@echo "> Building for debug"
	go build -o bin/${APP_NAME} ${DEBUG_FLAGS} cmd/main.go