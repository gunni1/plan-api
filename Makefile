.PHONY: all build docker docker-build prepare test build-api clean

TAG=gunni1/plan-api:master
DOCKERFILE=Dockerfile

all: prepare test build

build: build-api

prepare:
	go mod tidy
	go fmt ./...
	go vet ./...
	golint ./...

test:
	go test ./...

build-api:
	go build -o bin/runserver ./api/

docker: prepare test docker-build

docker-build:
	docker build -t $(TAG) -f $(DOCKERFILE) .

clean:
rm bin/*