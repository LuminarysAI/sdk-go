//go:build wasm

package sdk

//go:wasmimport env fs_read
func hostFsRead(ptr, length uint32) uint64

//go:wasmimport env fs_write
func hostFsWrite(ptr, length uint32) uint64

//go:wasmimport env fs_create
func hostFsCreate(ptr, length uint32) uint64

//go:wasmimport env fs_delete
func hostFsDelete(ptr, length uint32) uint64

//go:wasmimport env fs_mkdir
func hostFsMkdir(ptr, length uint32) uint64

//go:wasmimport env fs_ls
func hostFsLs(ptr, length uint32) uint64

//go:wasmimport env fs_chmod
func hostFsChmod(ptr, length uint32) uint64

//go:wasmimport env fs_read_lines
func hostFsReadLines(ptr, length uint32) uint64

//go:wasmimport env fs_grep
func hostFsGrep(ptr, length uint32) uint64

//go:wasmimport env fs_glob
func hostFsGlob(ptr, length uint32) uint64

//go:wasmimport env fs_allowed_dirs
func hostFsAllowedDirs(ptr, length uint32) uint64

//go:wasmimport env fs_copy
func hostFsCopy(ptr, length uint32) uint64
