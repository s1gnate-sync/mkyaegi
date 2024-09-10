cflags ?= -l -B
ldflags ?= -w -s

cmd:
	CGO_ENABLED=0 go build -a -trimpath -gcflags=all="$(cflags)" -ldflags="$(ldflags)" cmd/yaegi.go

lib:
	GOBIN=$(PWD)/tools test -e tools/extract || go install github.com/traefik/yaegi/internal/cmd/extract@latest
	PATH="$(PWD)/tools:$(PATH)" go generate -x lib/go1.23-generate.go

ref:
	tools/gen-reference

all: lib cmd ref	

.PHONY: all cmd lib ref
