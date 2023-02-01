compile:
	tinygo build -no-debug -o wasm/reverse.wasm -target wasi ./cmd/reverse/main.go
