cflags ?= -l -B
ldflags ?= -w -s

yaegi:
	CGO_ENABLED=0 go build -a -trimpath -gcflags=all="$(cflags)" -ldflags="$(ldflags)" yaegi.go

generate: extract lib/go1.23-generate.go
	go generate lib/go1.23-generate.go

extract:
	GOBIN=$(PWD) go install github.com/traefik/yaegi/internal/cmd/extract@latest

.PHONY: generate
