//go:build !wasm

package sdk

// Stubs for non-WASM builds (IDE, go vet, unit tests).
// These functions are never called outside a WASM runtime.

func hostTcpConnect(ptr uint32, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostTcpSetCallback(ptr uint32, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostTcpWrite(ptr uint32, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostTcpClose(ptr uint32, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostTcpRequest(ptr uint32, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
