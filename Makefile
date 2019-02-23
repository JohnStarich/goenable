GO111MODULE := on
CGO_ENABLED := 1
export
TARGETS := darwin/amd64,linux/amd64
GO_VERSION := 1.11.5
BASH_VERSION := 5.0

.PHONY: all
all: goenable plugins

.PHONY: dist
dist: out
	cd /tmp; GO111MODULE=auto go get -u github.com/johnstarich/xgo  # avoid updating go.mod files
	xgo \
		--buildmode=c-shared \
		--deps="http://ftpmirror.gnu.org/bash/bash-${BASH_VERSION}.tar.gz" \
		--depsargs="--disable-nls" \
		--dest=out \
		--go="${GO_VERSION}" \
		--image="johnstarich/xgo:1.11-nano" \
		--targets="${TARGETS}" \
		.
	mv out/github.com/johnstarich/* out/ && rm -rf out/github.com
	go run ./cmd/rename_binaries.go ./out

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
