package sdk

import "github.com/vmihailenco/msgpack/v5"

// ErrorKind classifies a TCP connection error.
type ErrorKind string

const (
	ErrorKindNone    ErrorKind = ""
	ErrorKindEOF     ErrorKind = "eof"
	ErrorKindReset   ErrorKind = "reset"
	ErrorKindTimeout ErrorKind = "timeout"
	ErrorKindTLS     ErrorKind = "tls"
	ErrorKindIO      ErrorKind = "io"
)

// ConnEvent is the payload delivered to the skill's read callback.
type ConnEvent struct {
	ConnID    string    `msgpack:"conn_id"`
	Data      []byte    `msgpack:"data,omitempty"`
	ErrorKind ErrorKind `msgpack:"error_kind,omitempty"`
	ErrorMsg  string    `msgpack:"error_msg,omitempty"`
}

// UnmarshalConnEvent deserialises a ConnEvent from raw bytes.
func UnmarshalConnEvent(payload []byte) (ConnEvent, error) {
	var evt ConnEvent
	if err := msgpack.Unmarshal(payload, &evt); err != nil {
		return ConnEvent{}, err
	}
	return evt, nil
}

// TcpConnectOptions configures a persistent TCP connection. Requires tcp.enabled.
type TcpConnectOptions struct {
	// Host:port to connect to (required).
	Addr string `msgpack:"addr"`
	// Callback method name called with ConnEvent on reads. "" = drain silently.
	Callback string `msgpack:"callback"`
	// TLS enables TLS encryption.
	TLS bool `msgpack:"tls"`
	// Insecure skips TLS certificate verification (dev only).
	Insecure bool `msgpack:"insecure"`
	// ServerName overrides the TLS SNI hostname.
	ServerName string `msgpack:"server_name,omitempty"`
	// TimeoutMs is the dial timeout in milliseconds. 0 = 30s default.
	TimeoutMs int64 `msgpack:"timeout_ms"`
}

// TcpConnect dials a TCP connection (plain or TLS). Requires tcp.enabled.
func TcpConnect(opts TcpConnectOptions) (string, error) {
	b := mustMarshal(opts)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostTcpConnect(ptr, ln))
	return extractConnID(raw, "tcp_connect")
}

// TcpSetCallback updates the read callback for an existing connection.
func TcpSetCallback(connID, callback string) error {
	req := map[string]interface{}{
		"conn_id":  connID,
		"callback": callback,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostTcpSetCallback(ptr, ln))
	return extractError(raw, "tcp_set_callback")
}

// TcpWrite sends data over an existing connection.
func TcpWrite(connID string, data []byte) error {
	req := map[string]interface{}{
		"conn_id": connID,
		"data":    data,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostTcpWrite(ptr, ln))
	return extractError(raw, "tcp_write")
}

// TcpClose closes the connection. Idempotent.
func TcpClose(connID string) error {
	req := map[string]interface{}{"conn_id": connID}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostTcpClose(ptr, ln))
	return extractError(raw, "tcp_close")
}

// TcpRequestOptions configures a synchronous TCP request. Requires tcp.enabled.
type TcpRequestOptions struct {
	// Host:port to connect to (required).
	Addr string `msgpack:"addr"`
	// Data is the payload to send (required).
	Data []byte `msgpack:"data"`
	// TLS enables TLS for the connection.
	TLS bool `msgpack:"tls"`
	// Insecure skips TLS certificate verification (dev only).
	Insecure bool `msgpack:"insecure"`
	// TimeoutMs is the total timeout for connect+send+read. 0 = 30s.
	TimeoutMs int `msgpack:"timeout_ms"`
	// MaxBytes limits the response size. 0 = 1MB.
	MaxBytes int `msgpack:"max_bytes"`
}

// TcpRequestResult holds the response from TcpRequest.
type TcpRequestResult struct {
	Data  []byte `msgpack:"data"`
	Error string `msgpack:"error,omitempty"`
}

// TcpRequest performs a synchronous TCP request: connect, send, read, close. Requires tcp.enabled.
func TcpRequest(opts TcpRequestOptions) (TcpRequestResult, error) {
	b := mustMarshal(opts)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostTcpRequest(ptr, ln))
	var resp TcpRequestResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return TcpRequestResult{}, &ABIError{"tcp_request: unmarshal: " + err.Error()}
	}
	if resp.Error != "" {
		return TcpRequestResult{}, &ABIError{resp.Error}
	}
	return resp, nil
}

func extractConnID(raw []byte, op string) (string, error) {
	var resp struct {
		ConnID string `msgpack:"conn_id"`
		Error  string `msgpack:"error"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return "", &ABIError{op + ": unmarshal response: " + err.Error()}
	}
	if resp.Error != "" {
		return "", &ABIError{resp.Error}
	}
	return resp.ConnID, nil
}

func extractError(raw []byte, op string) error {
	var resp struct {
		Error string `msgpack:"error"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return &ABIError{op + ": unmarshal response: " + err.Error()}
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}
