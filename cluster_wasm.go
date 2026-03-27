//go:build wasm

package sdk

//go:wasmimport env cluster_node_list
func hostClusterNodeList(ptr, length uint32) uint64

//go:wasmimport env file_transfer_send
func hostFileTransferSend(ptr, length uint32) uint64

//go:wasmimport env file_transfer_recv
func hostFileTransferRecv(ptr, length uint32) uint64
