//go:build wasm

package sdk

//go:wasmimport env history_get
func hostHistoryGet(ptr, length uint32) uint64

//go:wasmimport env prompt_complete
func hostPromptComplete(ptr, length uint32) uint64
