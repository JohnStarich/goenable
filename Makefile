.PHONY: all
all: goenable plugins

.PHONY: plugins
plugins: out
	set -ex; \
	for d in $$(ls examples); do \
		go build -v -o out/"$$d" -buildmode=plugin ./examples/"$$d"; \
	done

out:
	mkdir out

.PHONY: clean
clean:
	rm -rf out

.PHONY: goenable
goenable: out
	go build -v -o out/goenable.so -buildmode=c-shared .
