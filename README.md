## go-wasi

Opinionated Go package for running WASI binaries in Go (using the `tetratelabs/wazero` package).

## Motivation

This package is meant to be a simple wrapper to invoke a WASM (WASI) binary with zero or more string arguments returning the output of that binary as a `[]byte` instance.

It is as much to help me understand the boundaries of what is or isn't possible with Go and WASM (WASI) right now. Under the hood it relies on the [tetratelabs/wazero](https://github.com/tetratelabs/wazero) package for doing all the heavy lifting.

In order to compile WASI binaries derived from Go code you will need to install and build [tinygo](https://tinygo.org) and you should consider that there is still quite a lot of the standard Go programming language that tinygo [does not support yet](https://tinygo.org/docs/reference/lang-support/).

## Example

```
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/aaronland/go-wasi"
)

func main() {

	ctx := context.Background()

	path_wasm := flag.String("wasm", "", "...")
	
	flag.Parse()
	args := flag.Args()

	wasm_r, _ := os.Open(*path_wasm)
	result, _ := wasi.Run(ctx, wasm_r, args...)

	fmt.Println(string(result))
}
```

_Error handling omitted for the sake of brevity._

### reverse

```
$> go run cmd/reverse/main.go foo bar baz
baz bar foo
```

```
$> tinygo build -no-debug -o wasm/reverse.wasm -target wasi ./cmd/reverse/main.go

$> wasmtime wasm/reverse.wasm foo bar baz
baz bar foo
```

```
$> go run cmd/run/main.go -wasm ./wasm/reverse.wasm foo bar baz
baz bar foo
```

## See also

* https://github.com/tetratelabs/wazero
* https://tinygo.org