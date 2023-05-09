package wasi

import (
	"context"
	"os"
	"testing"
)

func TestRun(t *testing.T) {

	ctx := context.Background()

	path_wasm := "wasm/reverse.wasm"

	wasm_r, err := os.Open(path_wasm)

	if err != nil {
		t.Fatalf("Failed to open %s for reading, %v", path_wasm, err)
	}

	result, err := Run(ctx, wasm_r, "foo", "bar", "baz")

	if err != nil {
		t.Fatalf("Failed to run wasm binary, %v", err)
	}

	if string(result) != "baz bar foo" {
		t.Fatalf("Unexpected result from wasm binary, '%s'", string(result))
	}
}
