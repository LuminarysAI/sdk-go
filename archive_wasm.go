//go:build wasm

package sdk

//go:wasmimport env archive_pack
func hostArchivePack(ptr, length uint32) uint64

//go:wasmimport env archive_unpack
func hostArchiveUnpack(ptr, length uint32) uint64

//go:wasmimport env archive_list
func hostArchiveList(ptr, length uint32) uint64
