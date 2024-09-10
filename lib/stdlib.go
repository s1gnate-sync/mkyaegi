package lib

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"github.com/traefik/yaegi/stdlib/syscall"
	"github.com/traefik/yaegi/stdlib/unrestricted"
	"github.com/traefik/yaegi/stdlib/unsafe"
)

var Symbols = interp.Exports{}

func init() {
	Symbols["lib/lib"] = map[string]reflect.Value{
		"Packages": reflect.ValueOf(Packages),
	}
}

func exportsList() []interp.Exports {
	return []interp.Exports{
		stdlib.Symbols,
		unrestricted.Symbols,
		interp.Symbols,
		syscall.Symbols,
		unsafe.Symbols,
		Symbols,
	}
}

func Load(instance *interp.Interpreter) error {
	for _, exports := range exportsList() {
		if err := instance.Use(exports); err != nil {
			return err
		}
	}

	return nil
}

func Packages() []string {
	packages := []string{}
	for _, exports := range exportsList() {
		for pkg := range exports {
			if strings.HasPrefix(pkg, "github.com/traefik/yaegi/stdlib") {
				continue
			}

			index := strings.LastIndex(pkg, "/")
			if index >= 0 {
				packages = append(packages, fmt.Sprintf("%s:%s", pkg[0:index], pkg[index+1:]))
			}
		}
	}
	return packages
}
