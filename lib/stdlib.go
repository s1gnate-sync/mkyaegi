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
		"Exports":  reflect.ValueOf(Exports),
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
			parts := strings.Split(pkg, "/")
			packages = append(packages, strings.Join(parts[0:len(parts)-1], "/"))
		}
	}
	return packages
}

func Exports(query string) ([]string, error) {
	for _, exports := range exportsList() {
		for pkg, symbols := range exports {
			parts := strings.Split(pkg, "/")
			if strings.Join(parts[0:len(parts)-1], "/") == query || query == parts[len(parts)-1] {
				names := []string{}

				for name := range symbols {
					names = append(names, name)
				}
				return names, nil
			}
		}
	}
	return nil, fmt.Errorf("'%s' not found", query)
}
