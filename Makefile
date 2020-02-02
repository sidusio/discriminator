all: lint build

build:
	go build ./cmd/discriminator

build-docker:
	docker build .

lint:
	golangci-lint run --fast

clean:
	rm ./discriminator