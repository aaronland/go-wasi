package wasi

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

// Run compiles the WASI binary in 'wasm_r' and then invokes it with 'args'. The output of that function
// invocation is returned as a byte array. For reasons I don't understand it's not possible to wrap and
// reuse the underlying tetratelabs/wazero runtime and compiled code in a struct so, as of this writing,
// it's necessary to compile everything from scratch every time this method is invoked.
// invokeWasmModule invokes the given WASM module (given as a file path),
// setting its env vars according to env. Returns the module's stdout.
func Run(ctx context.Context, wasmReader io.Reader, args ...string) ([]byte, error) {

	modname := "MODNAME"

	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)
	wasi_snapshot_preview1.MustInstantiate(ctx, r)

	// Instantiate the wasm runtime, setting up exported functions from the host
	// that the wasm module can use for logging purposes.
	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().
		WithFunc(func(v uint32) {
			log.Printf("[%v]: %v", modname, v)
		}).
		Export("log_i32").
		NewFunctionBuilder().
		WithFunc(func(ctx context.Context, mod api.Module, ptr uint32, len uint32) {
			// Read the string from the module's exported memory.
			if bytes, ok := mod.Memory().Read(ptr, len); ok {
				log.Printf("[%v]: %v", modname, string(bytes))
			} else {
				log.Printf("[%v]: log_string: unable to read wasm memory", modname)
			}
		}).
		Export("log_string").
		Instantiate(ctx)
	if err != nil {
		return nil, err
	}

	// Set up stdout redirection and env vars for the module.

	var stdoutBuf bytes.Buffer

	config := wazero.NewModuleConfig().WithStdout(&stdoutBuf)

	runtime_args := []string{
		"wasi",
	}

	runtime_args = append(runtime_args, args...)

	config = config.WithArgs(runtime_args...)

	wasmObj, err := io.ReadAll(wasmReader)

	if err != nil {
		return nil, err
	}

	// Instantiate the module. This invokes the _start function by default.
	_, err = r.InstantiateWithConfig(ctx, wasmObj, config)

	if err != nil {
		return nil, err
	}

	return bytes.TrimSpace(stdoutBuf.Bytes()), nil
}

func runV1(ctx context.Context, wasm_r io.Reader, args ...string) ([]byte, error) {

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
