local-vault:
	vault server -dev

setup-test-data:
	./helpers/setup.sh

.PHONY: install-osx
install-osx:
	cp ./bin/vaul7y/vaul7y /usr/local/bin/vaul7y

.PHONY: dev
dev: ## Build for the current development version
	@echo "==> Building Vaul7y..."
	@mkdir -p ./bin
	@CGO_ENABLED=0 go build -o ./bin/vaul7y ./cmd/vaul7y
	@rm -f $(GOPATH)/bin/vaul7y
	@cp ./bin/vaul7y/vaul7y $(GOPATH)/bin/vaul7y
	@echo "==> Done"

.PHONY: build
build:
	go build -o bin/vaul7y ./cmd/vaul7y

.PHONY: run
run:
	./bin/vaul7y

.PHONY: test
test:
	go test ./...