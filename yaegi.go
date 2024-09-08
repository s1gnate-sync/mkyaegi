package main

import (
	"flag"
	"fmt"
	"go/scanner"
	"io"
	"reflect"

	"github.com/chzyer/readline"
	"github.com/s1gnate-sync/mkyaegi/lib"
	"github.com/traefik/yaegi/interp"

	"log"
	"os"
	"strings"
)

const (
	prompt = "> "
)

func main() {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	noImportUsed := flags.Bool("no-import-used", false, "disable automatic import of used packages")
	keepInitMain := flags.Bool("keep-init-main", false, "do not blank init and main functions before entering repl")
	noEnv := flags.Bool("no-env", false, "start with empty environment")
	eval := flags.String("eval", "", "evaluate `code` before running script")
	replAfterScript := flags.Bool("repl", os.Getenv("YAEGI_REPL") != "", "enter repl after script")
	tags := flags.String("tags", "", "build tags")

	flags.Parse(os.Args[1:])

	env := []string{}
	if !*noEnv {
		env = os.Environ()
	}

	args := flags.Args()
	script := "-"
	if len(args) > 0 {
		script = args[0]
		args = args[1:]
	}

	instance, err := lib.Interp(interp.Options{
		BuildTags: strings.Split(*tags, ","),
		Args:      append([]string{script}, args...),
		Env:       env,
	}, !*noImportUsed)

	if err != nil {
		log.Fatal("init", "\n", err)
	}

	if *eval != "" {
		if _, err := instance.Eval(*eval); err != nil {
			log.Fatal(*eval, "\n", err)
		}
	}

	if script != "-" {
		if fi, err := os.Stat(script); err == nil && fi.Mode().IsRegular() {
			bytes, err := os.ReadFile(script)
			if err != nil {
				log.Fatal(script, "\n", err)
			}

			var evalErr error
			if str := string(bytes); strings.HasPrefix(str, "#!") {
				_, evalErr = instance.Eval(strings.Replace(str, "#!", "//", 1))
			} else {
				_, evalErr = instance.EvalPath(script)
			}

			if evalErr != nil {
				log.Fatal(script, "\n", evalErr)
			}

		} else {
			log.Fatal("unrecognized script argument")
		}
	}

	if script == "-" || *replAfterScript {
		if !*keepInitMain {
			instance.Eval("func main(){}\nfunc init(){}\n")
		}

		result, err := repl(instance, script)

		if err != nil {
			log.Fatal("repl", "\n", err)
		}

		if result.IsValid() {
			fmt.Println(result)
		}
	}
}

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
