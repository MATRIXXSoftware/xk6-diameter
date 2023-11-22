GO = $(shell which go 2>/dev/null)
DOCKER = $(shell which docker 2>/dev/null)

.PHONY: all clean xk6 build generator

all: clean build generator

xk6:
	$(GO) install go.k6.io/xk6/cmd/xk6@latest

generator:
	$(GO) build -o bin/dict_generator cmd/dict_generator/*.go

build: xk6
	xk6 build v0.37.0 --with github.com/matrixxsoftware/xk6-diameter=. --output bin/k6

docker: build
	$(DOCKER) build --no-cache --build-arg K6_BINARY=bin/k6 -t ghcr.io/matrixxsoftware/xk6-diameter:latest .

clean:
	$(RM) bin/k6 bin/dict_generator
