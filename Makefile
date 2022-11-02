#ARCH
ARCH="`uname -s`"
LINUX="Linux"
Darwin="Darwin"
IP="127.0.0.1"
Tag="latest"
Version="1.0.0"

.PHONY: dependency
dependency:
	go mod tidy

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: build
build: dependency fmt vet
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/subsocks

.PHOHY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/subsocks-linux-amd64

.PHOHY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/subsocks-win-amd64.exe

.PHONY: docker-build
docker-build:
	docker build -t periky/subsocks:$(Tag) --build-arg Version=$(Version) .

.PHOHY: docker-push
docker-push:
	docker push periky/subsocks:$(Tag)

.PHONY: cert
cert:
	bash cert.sh $(IP)
