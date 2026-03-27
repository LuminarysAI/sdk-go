//go:build wasm

package sdk

//go:wasmimport env http_get
func hostHttpGet(ptr, length uint32) uint64

//go:wasmimport env http_post
func hostHttpPost(ptr, length uint32) uint64

//go:wasmimport env http_request
func hostHttpRequest(ptr, length uint32) uint64
