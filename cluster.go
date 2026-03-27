package sdk

import (
	"github.com/vmihailenco/msgpack/v5"
)

// ClusterNodeInfo describes a cluster node.
type ClusterNodeInfo struct {
	NodeID string   `msgpack:"node_id"`
	Role   string   `msgpack:"role"`
	Skills []string `msgpack:"skills"`
}

// ClusterNodeListResult is the response from ClusterNodeList.
type ClusterNodeListResult struct {
	CurrentNode string            `msgpack:"current_node"`
	Nodes       []ClusterNodeInfo `msgpack:"nodes"`
	Error       string            `msgpack:"error,omitempty"`
}

// ClusterNodeList returns the list of known cluster nodes.
func ClusterNodeList() (ClusterNodeListResult, error) {
	b := mustMarshal(map[string]string{})
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostClusterNodeList(ptr, ln))
	var resp ClusterNodeListResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return ClusterNodeListResult{}, err
	}
	if resp.Error != "" {
		return ClusterNodeListResult{}, &ABIError{resp.Error}
	}
	return resp, nil
}

// FileTransferSend sends a local file to a remote cluster node.
func FileTransferSend(targetNode, localPath, remotePath string) error {
	req := map[string]string{
		"target_node": targetNode,
		"local_path":  localPath,
		"remote_path": remotePath,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFileTransferSend(ptr, ln))
	var resp struct {
		Error string `msgpack:"error,omitempty"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}

// FileTransferRecv requests a file from a remote cluster node (pull mode).
func FileTransferRecv(sourceNode, remotePath, localPath string) error {
	req := map[string]string{
		"source_node": sourceNode,
		"remote_path": remotePath,
		"local_path":  localPath,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFileTransferRecv(ptr, ln))
	var resp struct {
		Error string `msgpack:"error,omitempty"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}
