//go:build wasm

package sdk

//go:wasmimport env tcp_connect
func hostTcpConnect(ptr uint32, length uint32) uint64

//go:wasmimport env tcp_set_callback
func hostTcpSetCallback(ptr uint32, length uint32) uint64

//go:wasmimport env tcp_write
func hostTcpWrite(ptr uint32, length uint32) uint64

//go:wasmimport env tcp_close
func hostTcpClose(ptr uint32, length uint32) uint64

//go:wasmimport env tcp_request
func hostTcpRequest(ptr uint32, length uint32) uint64
