XK6 = $(shell which xk6 2>/dev/null)

BIN := bin/k6

all: clean xk6-diameter

clean:
	$(RM) $(BIN)

xk6-diameter:
	$(XK6) build v0.37.0 --with github.com/matrixxsoftware/xk6-diameter=. --output $(BIN)
