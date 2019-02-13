.PHONY: all
all: goenable namespace
	:

out:
	mkdir out

.PHONY: goenable
goenable: out
	go build -v -o out/goenable.so -buildmode=c-shared .

.PHONY: namespace
namespace: out
	go build -v -o out/namespace -buildmode=plugin ./namespace

.PHONY: clean
clean:
	rm -rf out
