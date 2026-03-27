//go:build !wasm

package sdk

func hostClusterNodeList(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFileTransferSend(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}

func hostFileTransferRecv(ptr, length uint32) uint64 {
	panic("wasmimport: called outside WASM runtime")
}
