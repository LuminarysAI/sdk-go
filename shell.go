package sdk

import (
	"github.com/vmihailenco/msgpack/v5"
)

// ShellExecRequest is the request for ShellExec. Requires shell.enabled.
type ShellExecRequest struct {
	// Command is the shell command to execute (required).
	Command string `msgpack:"command"`
	// Workdir is the working directory. Defaults to sandbox root.
	Workdir string `msgpack:"workdir"`
	// TimeoutMs is the command timeout in milliseconds. 0 = 30s default.
	TimeoutMs int `msgpack:"timeout_ms"`
	// Tail returns only the last N lines of output. 0 = all lines.
	Tail int `msgpack:"tail"`
	// Grep filters output lines by regex pattern.
	Grep string `msgpack:"grep"`
	// AsDaemon starts the process in background and returns immediately with PID.
	AsDaemon bool `msgpack:"as_daemon"`
	// LogFile is the output file for daemon mode.
	LogFile string `msgpack:"log_file,omitempty"`
}

// ShellExecResult is the response from ShellExec.
type ShellExecResult struct {
	Output   string `msgpack:"output"`
	ExitCode int    `msgpack:"exit_code"`
	Pid      int    `msgpack:"pid,omitempty"`
	LogFile  string `msgpack:"log_file,omitempty"`
	Error    string `msgpack:"error,omitempty"`
}

// ShellExec executes a shell command and returns combined output and exit code. Requires shell.enabled.
func ShellExec(req ShellExecRequest) (ShellExecResult, error) {
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostShellExec(ptr, ln))
	var resp ShellExecResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return ShellExecResult{}, err
	}
	if resp.Error != "" {
		return ShellExecResult{}, &ABIError{resp.Error}
	}
	return resp, nil
}
