INTEGRATION  := rabbitmq
BINARY_NAME   = nri-$(INTEGRATION)
GOFLAGS		  = -mod=readonly
GOLANGCI_LINT = github.com/golangci/golangci-lint/cmd/golangci-lint
GOCOV         = github.com/axw/gocov/gocov
GOCOV_XML	  = github.com/AlekSi/gocov-xml

all: build

build: clean validate test compile

build-container:
	docker build -t nri-rabbitmq .

clean:
	@echo "=== $(INTEGRATION) === [ clean ]: Removing binaries and coverage file..."
	@rm -rfv bin coverage.xml

format:
	sh scripts/format.sh

validate:
	@printf "=== $(INTEGRATION) === [ validate ]: running golangci-lint & semgrep... "
	go run -mod=mod github.com/golangci/golangci-lint/cmd/golangci-lint run --verbose
	docker run --rm -v "${PWD}:/code:ro" --workdir /code returntocorp/semgrep -c /code/.semgrep.yml


bin/$(BINARY_NAME):
	@echo "=== $(INTEGRATION) === [ compile ]: building $(BINARY_NAME)..."
	@go build -v -o bin/$(BINARY_NAME) cmd/nri-postgresql/main.go

compile: bin/$(BINARY_NAME)

test:
	@echo "=== $(INTEGRATION) === [ test ]: running unit tests..."
	@go run $(GOFLAGS) $(GOCOV) test ./... | go run $(GOFLAGS) $(GOCOV_XML) > coverage.xml


integration-test:
	@echo "=== $(INTEGRATION) === [ test ]: running integration tests..."
	@go test -v -tags=integration ./tests/. || (ret=$$?; docker-compose -f tests/docker-compose.yml down && exit $$ret)
	@docker-compose -f tests/docker-compose.yml down


# Include thematic Makefiles
include $(CURDIR)/build/ci.mk
include $(CURDIR)/build/release.mk

.PHONY: all build clean tools tools-update deps validate compile test check-version
