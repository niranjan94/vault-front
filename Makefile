GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=vault-front

all: test build

build:
	cd ui && yarn build && cd ..
	rice embed-go
	$(GOBUILD) -o $(BINARY_NAME) -v -ldflags="-s -w"

test:
	bash scripts/test.sh

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
	go mod download
	mkdir -p .dev
	curl https://releases.hashicorp.com/vault/0.11.4/vault_0.11.4_linux_amd64.zip -o .dev/vault.zip
	unzip -d .dev .dev/vault.zip && rm -f .dev/vault.zip
	cd ui && yarn && yarn build
	cd $(GOPATH) && $(GOGET) -v github.com/oxequa/realize && $(GOGET) -v github.com/GeertJohan/go.rice/rice

