cflags ?= -l -B
ldflags ?= -w -s

default:
	@echo "Usage of make:" >&2
	@echo -e "	make debug\n\t\tbuild for debug" >&2
	@echo -e "	make release\n\t\tbuild for release" >&2
	@echo -e "	make library\n\t\tgenerate library exports" >&2
	@echo -e "	make reference\n\t\tgenerate package reference" >&2

debug:
	@echo "[debug] building yaegi..." >&2
	@go build -C cmd/run -o $(PWD)/yaegi

release:
	@echo "[release] building yaegi..." >&2
	@(CGO_ENABLED=0 go build -C cmd/run -o $(PWD)/yaegi -a -trimpath -gcflags=all="$(cflags)" -ldflags="$(ldflags)")

tools/extract:
	@! test -e tools/extract
	@(GOBIN=$(PWD)/tools go install github.com/traefik/yaegi/internal/cmd/extract@latest)

library:
	@make tools/extract --quiet
	@echo "[library] generating exports..." >&2
	@rm -f lib/go1_23_*
	@(PATH="$(PWD)/tools:$(PATH)" go generate -x lib/go1.23-generate.go)

reference:
	@echo "[reference] generating references..." >&2
	@tools/gen-reference

.PHONY: default release debug library reference tools/extract
