TARGETS := linux/amd64,darwin/amd64
GO_VERSION := 1.11
BASH_VERSION := 5.0

.PHONY: all
all: goenable plugins

.PHONY: dist
dist: out
	@cd /tmp; go get github.com/karalabe/xgo  # avoid updating go.mod files
	@set -ex; \
		cd out; \
		CGO_ENABLED=1 \
		GO111MODULE=on \
		xgo \
			-go ${GO_VERSION} \
			-buildmode=c-shared \
			--targets="${TARGETS}" \
			--deps=http://ftpmirror.gnu.org/bash/bash-${BASH_VERSION}.tar.gz \
			github.com/johnstarich/goenable

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
