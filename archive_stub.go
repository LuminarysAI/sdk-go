//go:build !wasm

package sdk

func hostArchivePack(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostArchiveUnpack(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostArchiveList(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
