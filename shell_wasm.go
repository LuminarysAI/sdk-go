//go:build wasm

package sdk

//go:wasmimport env shell_exec
func hostShellExec(ptr, length uint32) uint64
