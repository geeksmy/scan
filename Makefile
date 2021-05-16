# workdir info
PACKAGE=scan
PREFIX=$(shell pwd)
CMD_PACKAGE=${PACKAGE}
OUTPUT_DIR=${PREFIX}/bin
OUTPUT_FILE=${OUTPUT_DIR}/scan
COMMIT_ID=$(shell git rev-parse --short HEAD)
VERSION=$(shell git describe --tags || echo "v0.0.0")
VERSION_IMPORT_PATH=github.com/lneoe/go-help-libs/version
BUILD_TIME=$(shell date '+%Y-%m-%dT%H:%M:%S%Z')
VCS_BRANCH=$(shell git symbolic-ref --short -q HEAD)

# which golint
GOLINT=$(shell which golangci-lint || echo '')

# build args
BUILD_ARGS := \
    -ldflags "-s -w -X $(VERSION_IMPORT_PATH).appName=$(PACKAGE) \
    -X $(VERSION_IMPORT_PATH).version=$(VERSION) \
    -X $(VERSION_IMPORT_PATH).revision=$(COMMIT_ID) \
    -X $(VERSION_IMPORT_PATH).branch=$(VCS_BRANCH) \
    -X $(VERSION_IMPORT_PATH).buildDate=$(BUILD_TIME)"
EXTRA_BUILD_ARGS=

.PONY: lint test
default: lint test build

lint:
	@echo "+ $@"
	@$(if $(GOLINT), , \
		$(error Please install golint: `go get -u github.com/golangci/golangci-lint/cmd/golangci-lint`))
	golangci-lint run --deadline=10m --disable-all -E errcheck ./...

test:
	@echo "+ test"
	go test -cover $(EXTRA_BUILD_ARGS) ./...

build:
	@echo "+ build"
	go build $(BUILD_ARGS) $(EXTRA_BUILD_ARGS) -o ${OUTPUT_FILE} $(CMD_PACKAGE)

setup:
	mkdir -p bin/linux
	mkdir -p bin/osx
	mkdir -p bin/windows

clean:
	@echo "+ $@"
	@rm -r "${OUTPUT_DIR}"