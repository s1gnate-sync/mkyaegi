custom packages
---------------

```bash
$ cat lib/go1.23-generate.go 
//go:build go1.23

package lib

//go:generate extract github.com/bitfield/script

```

`make generate`


building
--------

`make`

to reduce size install `upx` and run `make yaegi0`

usage
-----

```
yaegi [OPTIONS] [SCRIPT [ARGS]]

Usage of ./yaegi:
  -eval code
    	evaluate code before running script
  -keep-init-main
    	do not blank init and main functions before entering repl
  -no-env
    	start with empty environment
  -no-import-used
    	disable automatic import of used packages
  -repl
    	enter repl after script
  -tags string
    	build tags


Examples:

	# repl
	yaegi

	# execute script
	yaegi script.go

	# eval code and start repl
	yaegi -eval "println(os.Args)" - a b c

	# eval code and execute script
	yaegi -eval "println(1)" script.go a b c

	# enter repl after script 
	yaegi -repl script.go 

	# same, but for shebang 
	YAEGI_REPL=1 ./script.go

	# all 3
	yaegi -repl -eval "runs before script" script.go
	> enters repl after script
	
```

