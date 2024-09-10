package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/traefik/yaegi/interp"
)

func exec(script string, instance *interp.Interpreter) error {
	fi, err := os.Stat(script)

	if err != nil {
		return err
	}

	if !fi.Mode().IsRegular() {
		return fmt.Errorf("permission denied")
	}

	bytes, err := os.ReadFile(script)
	fatal(err, "read", script)

	var evalErr error
	if str := string(bytes); strings.HasPrefix(str, "#!") {
		_, evalErr = instance.Eval(strings.Replace(str, "#!", "//", 1))
	} else {
		_, evalErr = instance.EvalPath(script)
	}

	return evalErr
}
