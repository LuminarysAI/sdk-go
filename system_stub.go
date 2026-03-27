//go:build !wasm

package sdk

func hostSysInfo(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostTimeNow(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostDiskUsage(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostEnvGet(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostLogWrite(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
