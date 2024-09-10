package main

import (
	"fmt"
	"go/scanner"
	"io"
	"reflect"

	"github.com/chzyer/readline"
	"github.com/traefik/yaegi/interp"

	"os"
	"strings"
)

const (
	prompt = "> "
)

func repl(instance *interp.Interpreter, name string) (reflect.Value, error) {
	var value reflect.Value

	reader, err := readline.New(prompt)
	if err != nil {
		return value, err
	}

	reader.CaptureExitSignal()
	defer reader.Close()

	lines := ""
	for {
		line, err := reader.Readline()
		if err == io.EOF {
			break
		}

		if err == readline.ErrInterrupt {
			lines = ""
			reader.SetPrompt(prompt)
			continue
		}

		lines += strings.TrimSpace(line) + "\n"

		value, err = instance.Eval(lines)
		if err != nil {
			switch e := err.(type) {
			case scanner.ErrorList:
				if len(e) > 0 && ignoreScannerError(e[0], lines) {
					reader.SetPrompt("")
					continue
				}
				fmt.Fprintln(os.Stderr, strings.TrimPrefix(e[0].Error(), name+":"))
			case interp.Panic:
				fmt.Fprintln(os.Stderr, e.Value)
				fmt.Fprintln(os.Stderr, string(e.Stack))
			default:
				fmt.Fprintln(os.Stderr, err)
			}
		}

		lines = ""

		if value.IsValid() {
			fmt.Fprintln(os.Stdout, value)
		}

		reader.SetPrompt(prompt)
	}

	return value, err
}

func ignoreScannerError(e *scanner.Error, s string) bool {
	msg := e.Msg
	if strings.HasSuffix(msg, "found 'EOF'") {
		return true
	}
	if msg == "raw string literal not terminated" {
		return true
	}
	if strings.HasPrefix(msg, "expected operand, found '}'") && !strings.HasSuffix(s, "}") {
		return true
	}
	return false
}
