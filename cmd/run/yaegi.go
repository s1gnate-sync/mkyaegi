package main

import (
	"flag"
	"fmt"

	"github.com/s1gnate-sync/mkyaegi/lib"
	"github.com/traefik/yaegi/interp"

	"os"
	"strings"
)

func parseFlags() ([]string, bool, bool, bool, string, bool, string, string) {
	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	noImportUsed := flags.Bool("no-import-used", false, "disable automatic import of used packages")
	keepInitMain := flags.Bool("keep-init-main", false, "do not blank init and main functions before entering repl")
	noEnv := flags.Bool("no-env", false, "start with empty environment")
	eval := flags.String("eval", "", "evaluate `code` before running script")
	enterRepl := flags.Bool("repl", os.Getenv("YAEGI_REPL") != "", "enter repl after script, environment YAEGI_REPL can be also used")
	tags := flags.String("tags", "", "build tags")

	flags.Parse(os.Args[1:])

	args := flag.Args()
	script := "-"
	if len(args) > 0 {
		script = args[0]
		args = args[1:]
	}

	return args, *noImportUsed, *keepInitMain, *noEnv, *eval, *enterRepl, *tags, script
}

func main() {
	args, noImportUsed, keepInitMain, noEnv, eval, enterRepl, tags, script := parseFlags()

	env := []string{}
	if !noEnv {
		env = os.Environ()
	}

	instance, err := lib.Interp(interp.Options{
		BuildTags: strings.Split(tags, ","),
		Args:      append([]string{script}, args...),
		Env:       env,
	}, !noImportUsed)

	fatal(err, "interp")

	if eval != "" {
		_, err := instance.Eval(eval)
		fatal(err, "eval")
	}

	if script != "-" {
		if err := exec(script, instance); err != nil {
			fatal(err, "exec", script)
		}
	}

	if enterRepl || (script == "-" && eval == "") {
		if !keepInitMain {
			instance.Eval("func main(){}\nfunc init(){}\n")
		}

		result, err := repl(instance, script)
		fatal(err, "repl")

		if result.IsValid() {
			fmt.Println(result)
		}
	}
}

func fatal(err error, tag ...string) {
	if err == nil {
		return
	}

	fmt.Fprintf(os.Stderr, "%s failed: %s\n", strings.Join(tag, " "), err)
	os.Exit(100)
}
