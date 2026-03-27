package sdk

import "github.com/vmihailenco/msgpack/v5"

const (
	WsMessageText   = "text"
	WsMessageBinary = "binary"
	WsMessageClose  = "close"
)

// WsEvent is the payload delivered to the skill's WebSocket callback.
type WsEvent struct {
	ConnID      string    `msgpack:"conn_id"`
	Data        []byte    `msgpack:"data,omitempty"`
	MessageType string    `msgpack:"message_type,omitempty"` // "text"|"binary"|"close"
	CloseCode   int       `msgpack:"close_code,omitempty"`
	CloseText   string    `msgpack:"close_text,omitempty"`
	ErrorKind   ErrorKind `msgpack:"error_kind,omitempty"`
	ErrorMsg    string    `msgpack:"error_msg,omitempty"`
}

// UnmarshalWsEvent deserialises a WsEvent from raw bytes.
func UnmarshalWsEvent(payload []byte) (WsEvent, error) {
	var evt WsEvent
	if err := msgpack.Unmarshal(payload, &evt); err != nil {
		return WsEvent{}, err
	}
	return evt, nil
}

// WsConnect dials a WebSocket connection. Requires http.enabled and http.allow_websocket.
func WsConnect(url string, headers []Header, timeoutMs int64, callback string, insecure bool) (string, error) {
	req := map[string]interface{}{
		"url":        url,
		"headers":    headers,
		"timeout_ms": timeoutMs,
		"callback":   callback,
		"insecure":   insecure,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostWsConnect(ptr, ln))

	var resp struct {
		ConnID string `msgpack:"conn_id"`
		Error  string `msgpack:"error"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return "", &ABIError{"ws_connect: unmarshal: " + err.Error()}
	}
	if resp.Error != "" {
		return "", &ABIError{resp.Error}
	}
	return resp.ConnID, nil
}

// WsSend sends a message over an existing WebSocket connection.
func WsSend(connID string, data []byte, messageType string) error {
	req := map[string]interface{}{
		"conn_id":      connID,
		"data":         data,
		"message_type": messageType,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostWsSend(ptr, ln))
	return extractError(raw, "ws_send")
}

// WsClose sends a Close frame and removes the connection.
func WsClose(connID string, code int, reason string) error {
	req := map[string]interface{}{
		"conn_id": connID,
		"code":    code,
		"reason":  reason,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostWsClose(ptr, ln))
	return extractError(raw, "ws_close")
}
