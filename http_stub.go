//go:build !wasm

package sdk

// Stubs for non-WASM builds (IDE, go vet, unit tests).
// These functions are never called outside a WASM runtime.

func hostHttpGet(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostHttpPost(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostHttpRequest(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
