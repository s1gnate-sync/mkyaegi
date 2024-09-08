cflags ?= -l -B
ldflags ?= -w -s

compile:
	CGO_ENABLED=0 go build -a -trimpath -gcflags=all="$(cflags)" -ldflags="$(ldflags)" yaegi.go

setup:
	go install github.com/traefik/yaegi/internal/cmd/extract

.PHONY: setup build
