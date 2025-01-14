BIN_DIR ?= $(PWD)/bin
BUILD_DIR ?= $(PWD)/build
OUT_DIR ?= $(PWD)/out
GOLANGCILINT ?= $(BIN_DIR)/golangci-lint
GOLANGCILINT_VERSION ?= 1.61.0
CLI ?= $(BUILD_DIR)/geo

# OPEN_WEATHER_MAP_API_KEY ?= not-set

SRC_FILES=$(shell find . -type f -name '*.go')

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

$(OUT_DIR):
	mkdir -p $(OUT_DIR)

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

$(GOLANGCILINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCILINT_VERSION)

.PHONY: golangci-lint
golangci-lint: $(GOLANGCILINT)

.PHONY: dependencies
dependencies: golangci-lint

.PHONY: lint
lint: $(GOLANGCILINT) $(SRC_FILES)
	$(GOLANGCILINT) run ./...

$(CLI): $(SRC_FILES)
	go build -o $(CLI) ./main.go


test: lint ## Test runs all go tests. Deliberately runs every test, no caching or source file checks.
	go test ./... -count=1

e2e: clean-fixtures test ## runs test with openweathermap API calls enabled. Recommend committing any changes to the "cassettes" directory after running

all: dependencies $(CLI)

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(BUILD_DIR)
	rm -rf $(OUT_DIR)

.PHONY: clean-fixtures
clean-fixtures: ## Remove VCR files for enemy/e2e tests
	find . -iname "*.yaml" | grep "cassettes\/[a-zA-Z0-9-]*\.yaml" | xargs -n1 rm