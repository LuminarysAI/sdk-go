//go:build !wasm

package sdk

func hostFsRead(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsWrite(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsCreate(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsDelete(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsMkdir(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsLs(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsChmod(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsReadLines(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsGrep(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsGlob(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsAllowedDirs(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFsCopy(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
