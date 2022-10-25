#! /bin/bash

ip=$1
if [ -z "$ip" ]
then
    echo -e "\033[31m请输入要签名的IP\033[0m"
    exit 127
fi

CERT_DIR=$(pwd)/certs
rm -rf ${CERT_DIR} && mkdir ${CERT_DIR}

# 生成自签名根证书
### 创建根证书私钥
openssl genrsa -out ${CERT_DIR}/root.key 4096
### 创建根证书请求文件
openssl req -subj "/C=CN/ST=Shanghai/L=Shanghai/O=socks/OU=socks/CN=root" -new \
    -out ${CERT_DIR}/root.csr -key ${CERT_DIR}/root.key
### 创建根证书
openssl x509 -req -in ${CERT_DIR}/root.csr -out ${CERT_DIR}/root.crt \
    -signkey ${CERT_DIR}/root.key -CAcreateserial -days 365


cat <<EOF > ${CERT_DIR}/mb.conf
[mb]
basicConstraints = critical, CA:false
subjectAltName = IP:${ip},IP:127.0.0.1
extendedKeyUsage=serverAuth,clientAuth
EOF

# 生成自签名服务端证书
### 生成服务端证书私钥
openssl genrsa -out ${CERT_DIR}/server.key 4096
### 生成服务证书请求文件
openssl req -subj "/C=cn/ST=Shanghai/L=Shanghai/O=subsocks/OU=subsocks" \
    -addext "basicConstraints = critical, CA:false, pathlen:0" \
    -addext "subjectAltName = IP:${ip},IP:127.0.0.1" \
    -addext "extendedKeyUsage=serverAuth,clientAuth" -new -out ${CERT_DIR}/server.csr \
    -key ${CERT_DIR}/server.key
### 生成服务端公钥证书
openssl x509 -req -in ${CERT_DIR}/server.csr -out ${CERT_DIR}/server.crt \
    -signkey ${CERT_DIR}/server.key -CA ${CERT_DIR}/root.crt -CAkey ${CERT_DIR}/root.key \
    -CAcreateserial -days 365 -extfile ${CERT_DIR}/mb.conf -extensions mb

# 生成自签名客户端证书
### 生成客户端证书私钥
openssl genrsa -out ${CERT_DIR}/client.key 4096
### 生成客户端证书请求文件
openssl req -subj "/C=cn/ST=Shanghai/L=Shanghai/O=subsocks/OU=subsocks" \
    -addext "basicConstraints = critical, CA:false, pathlen:0" \
    -addext "subjectAltName = IP:${ip},IP:127.0.0.1" \
    -addext "extendedKeyUsage=serverAuth,clientAuth" \
    -new -out ${CERT_DIR}/client.csr -key ${CERT_DIR}/client.key
### 生成客户端公钥证书
openssl x509 -req -in ${CERT_DIR}/client.csr -out ${CERT_DIR}/client.crt \
    -signkey ${CERT_DIR}/client.key -CA ${CERT_DIR}/root.crt -CAkey ${CERT_DIR}/root.key \
    -CAcreateserial -days 365 -extfile ${CERT_DIR}/mb.conf -extensions mb
