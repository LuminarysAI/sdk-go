//go:build wasm

package sdk

//go:wasmimport env ws_connect
func hostWsConnect(ptr, length uint32) uint64

//go:wasmimport env ws_send
func hostWsSend(ptr, length uint32) uint64

//go:wasmimport env ws_close
func hostWsClose(ptr, length uint32) uint64
