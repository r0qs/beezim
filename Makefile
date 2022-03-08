
LDFLAGS ?= -s -w
SRC_ROOT = .
CLI_DIR = $(SRC_ROOT)/cli
BIN_DIR = $(SRC_ROOT)/bin
BIN_CLI = $(BIN_DIR)/beezim-cli

# Go parameters
GOCMD ?= go
GOBUILD = $(GOCMD) build -v
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test -v -failfast
GORACE = GORACE="halt_on_error=1"
GOTEST_RACE = $(GORACE) $(GOTEST) -v -race

ifndef $(GOPATH)
GOPATH=$(shell echo $(shell $(GOCMD) env GOPATH) | sed -E "s;(.*):.*;\1;")
export GOPATH
endif

all: build

build: bin
	@echo "+ building beezim source"
	$(GOBUILD) -o $(BIN_CLI) $(CLI_DIR)

bin:
	@mkdir $@

install:
	@echo "+ installing beezim-cli"
	$(GOBUILD) -trimpath -ldflags "$(LDFLAGS)" -o $(BIN_CLI) $(CLI_DIR)
ifneq ("$(GOPATH)","")
	@mv $(BIN_CLI) $(GOPATH)/$(BIN_CLI)
endif

.PHONY: test
test:
	@echo "+ executing tests"
	$(GOCLEAN) -testcache && $(GOTEST) $(SRC_ROOT)/...

.PHONY: racetest
racetest:
	@echo "+ building tests using Race Detector"
	$(GOTEST_RACE) $(SRC_ROOT)/...

.PHONY: clean
clean:
	@echo "+ cleaning"
	$(GOCLEAN) -i -cache ./...
	@rm -r $(BIN_DIR)

.PHONY: codecheck
codecheck: fmt lint vet

.PHONY: fmt
fmt:
	@echo "+ go fmt"
	$(GOCMD) fmt $(SRC_ROOT)/...

.PHONY: lint
lint:
	@echo "+ go lint"
	@golint -min_confidence=0.1 $(SRC_ROOT)/...

.PHONY: vet
vet:
	@echo "+ go vet"
	$(GOCMD) vet $(SRC_ROOT)/...