GO = $(shell which go 2>/dev/null)

.PHONY: all clean xk6 build generator

all: clean build generator

xk6:
	go install go.k6.io/xk6/cmd/xk6@latest

generator:
	$(GO) build -o bin/dict_generator cmd/dict_generator/*.go

build: xk6
	xk6 build v0.37.0 --with github.com/matrixxsoftware/xk6-diameter=. --output bin/k6

clean:
	$(RM) bin/k6 bin/dict_generator
