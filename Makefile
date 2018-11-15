GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=vault-front

MAKEPID:= $(shell echo $$PPID)
WORKING_DIR:= $(shell pwd)

all: test build

build:
	cd ui && yarn build && cd ..
	rice embed-go
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags="-s -w"

test:
	sh scripts/test.sh

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f rice-box.go

run: build
	./$(BINARY_NAME)

dev:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	$(GOGET) -v github.com/oxequa/realize
	$(GOGET) -v github.com/GeertJohan/go.rice/rice
	dep ensure -v
	mkdir -p .dev
	curl https://releases.hashicorp.com/vault/0.11.4/vault_0.11.4_linux_amd64.zip -o .dev/vault.zip
	unzip -d .dev .dev/vault.zip && rm -f .dev/vault.zip
	cd ui && yarn && yarn build
