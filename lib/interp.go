package lib

import (
	"go/build"

	"github.com/traefik/yaegi/interp"
)

func Interp(options interp.Options, importUsed bool) (*interp.Interpreter, error) {
	options.Unrestricted = true
	options.GoPath = build.Default.GOPATH

	instance := interp.New(options)
	if err := Load(instance); err != nil {
		return nil, err
	}

	if importUsed {
		instance.ImportUsed()
	}

	return instance, nil
}
