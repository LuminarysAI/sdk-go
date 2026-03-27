//go:build !wasm

package sdk

func hostShellExec(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
