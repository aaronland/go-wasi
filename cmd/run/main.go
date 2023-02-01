package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aaronland/go-wasi"
)

func main() {

	ctx := context.Background()

	path_wasm := flag.String("wasm", "", "...")

	flag.Parse()

	args := flag.Args()

	wasm_r, err := os.Open(*path_wasm)

	if err != nil {
		log.Fatalf("Failed to open %s for reading, %v", *path_wasm, err)
	}

	result, err := wasi.Run(ctx, wasm_r, args...)

	if err != nil {
		log.Fatalf("Failed to run wasi binary, %v", err)
	}

	fmt.Println(string(result))
}
