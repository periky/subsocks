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
	go build -o bin/subsocks

.PHONY: docker-build
docker-build:
	docker build -t periky/subsocks:$(Tag) --build-arg Version=$(Version) .

.PHOHY: docker-push
docker-push:
	docker push periky/subsocks:$(Tag)

CERT_DIR := $(shell pwd)/certs
.PHONY: openssl
openssl:
	rm -rf $(CERT_DIR) && mkdir $(CERT_DIR)
	@echo "根证书"
	@echo "生成Key..."
	openssl genrsa -out $(CERT_DIR)/ca.key 4096
	@echo "生成密钥..."
	openssl req -subj "/C=CN/ST=Shanghai/L=Shanghai/O=socks/OU=socks/CN=proxy.wangpu.xyz" \
    	-new -key $(CERT_DIR)/ca.key -out $(CERT_DIR)/ca.csr
	@echo "生成自签名证书..."
	openssl x509 -req -days 365 -in $(CERT_DIR)/ca.csr -signkey $(CERT_DIR)/ca.key \
		-out $(CERT_DIR)/ca.crt

	@echo "Server"
	@echo "生成服务端私钥..."
	openssl genrsa -out $(CERT_DIR)/server.key 4096
	@echo "生成CSR"
	openssl req -subj "/C=CN/ST=Shanghai/L=Shanghai/O=socks/OU=socks/CN=proxy.wangpu.xyz" \
		-new -key $(CERT_DIR)/server.key -out $(CERT_DIR)/server.csr
	@echo "生成服务端证书"
	openssl x509 -req -days 365 -CA $(CERT_DIR)/ca.crt -CAkey $(CERT_DIR)/ca.key -CAcreateserial \
		-in $(CERT_DIR)/server.csr -out $(CERT_DIR)/server.crt

	@echo "Client"
	@echo "生成客户端私钥..."
	openssl genrsa -out $(CERT_DIR)/client.key 4096
	@echo "生成CSR"
	openssl req -subj "/C=CN/ST=Shanghai/L=Shanghai/O=socks/OU=socks/CN=proxy.wangpu.xyz" \
		-new -key $(CERT_DIR)/client.key -out $(CERT_DIR)/client.csr
	@echo "生成客户端证书"
	openssl x509 -req -days 365 -CA $(CERT_DIR)/ca.crt -CAkey $(CERT_DIR)/ca.key -CAcreateserial \
		-in $(CERT_DIR)/client.csr -out $(CERT_DIR)/client.crt
