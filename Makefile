.PHONY: build clean deploy

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/plan-api ./

test:
	go test ./...

clean:
	rm -rf ./bin



