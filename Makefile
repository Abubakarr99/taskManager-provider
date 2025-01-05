HOSTNAME=dantata.com
NAMESPACE=aboudev
NAME=taskmanager
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
GOARCH  := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)

default: install

build:
	go build -o ${BINARY}

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${GOOS}_${GOARCH}