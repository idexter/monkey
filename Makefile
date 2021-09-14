.PHONY: build

all: clean build test run

build:
	go build -o bin/monkeyc main.go

test:
	go test ./...

clean:
	rm -Rf ./bin

run: build
	./bin/monkeyc
