package lib

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
	"reflect"
)

var Symbols = map[string]map[string]reflect.Value{}

func Load(instance *interp.Interpreter) (err error) {
	if err = instance.Use(stdlib.Symbols); err != nil {
		return
	}

	if err = instance.Use(interp.Symbols); err != nil {
		return
	}

	if err = instance.Use(syscall.Symbols); err != nil {
		return
	}

	if err = instance.Use(unsafe.Symbols); err != nil {
		return
	}

	if err = instance.Use(unrestricted.Symbols); err != nil {
		return
	}

	if err = instance.Use(Symbols); err != nil {
		return
	}

	return nil
}
