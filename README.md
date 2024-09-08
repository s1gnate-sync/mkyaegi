custom packages
---------------

```bash
$ cat lib/go1.23-generate.go 
//go:build go1.23

package lib

//go:generate extract github.com/bitfield/script

```

`$ make generate`


building
--------

`$ make`

usage
-----

```
yaegi [OPTIONS] [SCRIPT [ARGS]]

Usage:
  -eval code
    	evaluate code before running script
  -no-env
    	start with empty environment
  -no-import-used
    	disable automatic import of used packages
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
	
```

