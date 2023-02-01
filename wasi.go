package wasi

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Run compiles the WASI binary in 'wasm_r' and then invokes it with 'args'. The output of that function
// invocation is returned as a byte array. For reasons I don't understand it's not possible to wrap and
// reuse the underlying tetratelabs/wazero runtime and compiled code in a struct so, as of this writing,
// it's necessary to compile everything from scratch every time this method is invoked.
func Run(ctx context.Context, wasm_r io.Reader, args ...string) ([]byte, error) {

	runtime := wazero.NewRuntime(ctx)
	defer runtime.Close(ctx)

	wasi_snapshot_preview1.MustInstantiate(ctx, runtime)

	wasm_body, err := io.ReadAll(wasm_r)

	if err != nil {
		return nil, fmt.Errorf("Failed to read wasm, %w", err)
	}

	code, err := runtime.CompileModule(ctx, wasm_body)

	if err != nil {
		return nil, fmt.Errorf("Failed to compile module, %w", err)
	}

	var buf bytes.Buffer
	buf_wr := bufio.NewWriter(&buf)

	config := wazero.NewModuleConfig().WithStdout(buf_wr).WithStderr(os.Stderr)

	runtime_args := []string{
		"wasi",
	}

	runtime_args = append(runtime_args, args...)

	_, err = runtime.InstantiateModule(ctx, code, config.WithArgs(runtime_args...))

	if err != nil {
		return nil, fmt.Errorf("Failed to instantiate module, %w", err)
	}

	buf_wr.Flush()
	
	return bytes.TrimSpace(buf.Bytes()), nil
}
