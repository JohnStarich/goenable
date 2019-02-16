TARGETS := darwin/amd64,linux/amd64
GO_VERSION := 1.11
BASH_VERSION := 5.0
# Set default remote and branch, but allow env var overrides:
#   DIST_REMOTE, DIST_BRANCH
DIST_REMOTE := $(shell [[ -n "$${DIST_REMOTE}" ]] && echo "--remote=$${DIST_REMOTE}")
DIST_BRANCH := $(shell echo "$${DIST_BRANCH:-$${TRAVIS_TAG:-master}}")

.PHONY: all
all: goenable plugins

.PHONY: dist
dist: out
	cd /tmp; go get github.com/karalabe/xgo  # avoid updating go.mod files
	@set -ex; \
		CGO_ENABLED=1 \
		GO111MODULE=on \
		xgo \
			${DIST_REMOTE} \
			--branch=${DIST_BRANCH} \
			--buildmode=c-shared \
			--deps="http://ftpmirror.gnu.org/bash/bash-${BASH_VERSION}.tar.gz" \
			--depsargs="--disable-nls" \
			--dest=out \
			--go="${GO_VERSION}" \
			--image="johnstarich/xgo:1.11-nano" \
			--targets="${TARGETS}" \
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
