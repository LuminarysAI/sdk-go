//go:build wasm

package sdk

//go:wasmimport env sys_info
func hostSysInfo(ptr, length uint32) uint64

//go:wasmimport env time_now
func hostTimeNow(ptr, length uint32) uint64

//go:wasmimport env disk_usage
func hostDiskUsage(ptr, length uint32) uint64

//go:wasmimport env env_get
func hostEnvGet(ptr, length uint32) uint64

//go:wasmimport env log_write
func hostLogWrite(ptr, length uint32) uint64
