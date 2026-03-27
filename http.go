package sdk

import (
	"github.com/vmihailenco/msgpack/v5"
)

// Header is an ordered HTTP header name/value pair.
type Header struct {
	Name  string `msgpack:"name"`
	Value string `msgpack:"value"`
}

// Cookie represents an HTTP cookie.
type Cookie struct {
	Name     string `msgpack:"name"`
	Value    string `msgpack:"value"`
	Domain   string `msgpack:"domain,omitempty"`
	Path     string `msgpack:"path,omitempty"`
	Expires  int64  `msgpack:"expires,omitempty"` // Unix timestamp; 0 = session
	Secure   bool   `msgpack:"secure,omitempty"`
	HTTPOnly bool   `msgpack:"httponly,omitempty"`
}

// HttpResponse is the response returned by all HTTP calls.
type HttpResponse struct {
	Status    int      `msgpack:"status"`
	Headers   []Header `msgpack:"headers,omitempty"`
	Cookies   []Cookie `msgpack:"cookies,omitempty"`
	Body      []byte   `msgpack:"body,omitempty"`
	Truncated bool     `msgpack:"truncated,omitempty"`
}

// HttpGet performs a GET request. Requires http.enabled.
//
// timeoutMs=0 uses 30 s default. maxBytes=0 uses the manifest limit (or 1 MiB default).
func HttpGet(url string, timeoutMs int64, maxBytes int) (HttpResponse, error) {
	req := map[string]interface{}{
		"url":        url,
		"timeout_ms": timeoutMs,
		"max_bytes":  maxBytes,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostHttpGet(ptr, ln))
	return unmarshalHttpResponse(raw, "http_get")
}

// HttpPost performs a POST request. Requires http.enabled.
//
// contentType defaults to "application/octet-stream" if empty.
func HttpPost(url string, body []byte, contentType string, timeoutMs int64, maxBytes int) (HttpResponse, error) {
	req := map[string]interface{}{
		"url":          url,
		"body":         body,
		"content_type": contentType,
		"timeout_ms":   timeoutMs,
		"max_bytes":    maxBytes,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostHttpPost(ptr, ln))
	return unmarshalHttpResponse(raw, "http_post")
}

// HttpRequestOptions configures a fully customised HTTP request.
type HttpRequestOptions struct {
	Method          string // defaults to "GET"
	URL             string
	Headers         []Header // applied in listed order
	Cookies         []Cookie
	Body            []byte
	TimeoutMs       int64
	MaxBytes        int // 0 = manifest limit (1 MiB default)
	FollowRedirects bool
	UseJar          bool // enable the persistent per-skill cookie jar
}

// HttpRequest performs a fully customised HTTP request with ordered headers,
// persistent cookie jar, and full control over method, redirects and body.
// Requires http.enabled.
func HttpRequest(opts HttpRequestOptions) (HttpResponse, error) {
	req := map[string]interface{}{
		"method":           opts.Method,
		"url":              opts.URL,
		"headers":          opts.Headers,
		"cookies":          opts.Cookies,
		"body":             opts.Body,
		"timeout_ms":       opts.TimeoutMs,
		"max_bytes":        opts.MaxBytes,
		"follow_redirects": opts.FollowRedirects,
		"use_jar":          opts.UseJar,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostHttpRequest(ptr, ln))
	return unmarshalHttpResponse(raw, "http_request")
}

// HeadersFromJSON parses a JSON object string into an ordered Header slice.
// Preserves key insertion order from the source string.
func HeadersFromJSON(jsonStr string) []Header {
	s := jsonStr
	var headers []Header
	// Simple ordered parser — reads key-value pairs sequentially.
	i := 0
	for i < len(s) && s[i] != '{' {
		i++
	}
	if i >= len(s) {
		return headers
	}
	i++ // skip '{'
	for i < len(s) {
		// skip to opening '"'
		for i < len(s) && s[i] != '"' && s[i] != '}' {
			i++
		}
		if i >= len(s) || s[i] == '}' {
			break
		}
		i++ // skip '"'
		keyStart := i
		i = findQuoteEnd(s, i)
		key := s[keyStart:i]
		i++ // skip closing '"'
		// skip ':' and whitespace to opening '"'
		for i < len(s) && s[i] != '"' {
			i++
		}
		if i >= len(s) {
			break
		}
		i++ // skip '"'
		valStart := i
		i = findQuoteEnd(s, i)
		val := s[valStart:i]
		i++ // skip closing '"'
		headers = append(headers, Header{Name: key, Value: val})
	}
	return headers
}

func findQuoteEnd(s string, pos int) int {
	i := pos
	for i < len(s) {
		if s[i] == '"' {
			return i
		}
		if s[i] == '\\' {
			i += 2
			continue
		}
		i++
	}
	return i
}

func unmarshalHttpResponse(raw []byte, op string) (HttpResponse, error) {
	var resp struct {
		HttpResponse
		Error string `msgpack:"error"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return HttpResponse{}, &ABIError{op + ": unmarshal response: " + err.Error()}
	}
	if resp.Error != "" {
		return HttpResponse{}, &ABIError{resp.Error}
	}
	return resp.HttpResponse, nil
}
