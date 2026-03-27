//go:build !wasm

package sdk

func hostHistoryGet(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostPromptComplete(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
