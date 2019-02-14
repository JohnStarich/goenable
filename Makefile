.PHONY: all
all: goenable plugins

.PHONY: plugins
plugins: namespace hello pow

out:
	mkdir out

.PHONY: clean
clean:
	rm -rf out

.PHONY: goenable
goenable: out
	go build -v -o out/goenable.so -buildmode=c-shared .

# Plugins

.PHONY: namespace
namespace: out
	go build -v -o out/namespace -buildmode=plugin ./namespace

.PHONY: hello
hello: out
	go build -v -o out/hello -buildmode=plugin ./hello

.PHONY: pow
pow: out
	go build -v -o out/pow -buildmode=plugin ./pow
