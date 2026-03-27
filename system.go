package sdk

import (
	"github.com/vmihailenco/msgpack/v5"
)

// SysInfoResult holds host OS and hardware information.
type SysInfoResult struct {
	OS       string `msgpack:"os"`
	Arch     string `msgpack:"arch"`
	Hostname string `msgpack:"hostname"`
	NumCPU   int    `msgpack:"num_cpu"`
}

// SysInfo returns information about the host OS and hardware.
func SysInfo() (SysInfoResult, error) {
	b := mustMarshal(map[string]any{})
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostSysInfo(ptr, ln))
	var resp SysInfoResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return SysInfoResult{}, err
	}
	return resp, nil
}

// TimeNowResult holds the current host time.
type TimeNowResult struct {
	Unix      int64  `msgpack:"unix"`
	UnixNano  int64  `msgpack:"unix_nano"`
	RFC3339   string `msgpack:"rfc3339"`
	Timezone  string `msgpack:"timezone"`
	UTCOffset int    `msgpack:"utc_offset"`
}

// TimeNow returns the current host time.
func TimeNow() (TimeNowResult, error) {
	b := mustMarshal(map[string]any{})
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostTimeNow(ptr, ln))
	var resp TimeNowResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return TimeNowResult{}, err
	}
	return resp, nil
}

// DiskUsageResult holds disk space information.
type DiskUsageResult struct {
	TotalBytes int64   `msgpack:"total_bytes"`
	FreeBytes  int64   `msgpack:"free_bytes"`
	UsedBytes  int64   `msgpack:"used_bytes"`
	UsedPct    float64 `msgpack:"used_pct"`
	Error      string  `msgpack:"error,omitempty"`
}

// DiskUsage returns disk space information for the given path. Requires fs.enabled.
func DiskUsage(path string) (DiskUsageResult, error) {
	b := mustMarshal(map[string]string{"path": path})
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostDiskUsage(ptr, ln))
	var resp DiskUsageResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return DiskUsageResult{}, err
	}
	if resp.Error != "" {
		return DiskUsageResult{}, &ABIError{resp.Error}
	}
	return resp, nil
}

// GetEnv reads a named value declared in the skill's manifest under env.
func GetEnv(key string) string {
	req := map[string]string{"key": key}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostEnvGet(ptr, ln))
	var resp struct {
		Value string `msgpack:"value"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return ""
	}
	return resp.Value
}

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

// Log writes a structured log message.
func Log(level, message string, fields map[string]any) {
	req := map[string]any{
		"level":   level,
		"message": message,
		"fields":  fields,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	readHostResult(hostLogWrite(ptr, ln))
}

// LogDebug is a shorthand for Log("debug", ...).
func LogDebug(message string, fields map[string]any) { Log("debug", message, fields) }

// LogInfo is a shorthand for Log("info", ...).
func LogInfo(message string, fields map[string]any) { Log("info", message, fields) }

// LogWarn is a shorthand for Log("warn", ...).
func LogWarn(message string, fields map[string]any) { Log("warn", message, fields) }

// LogError is a shorthand for Log("error", ...).
func LogError(message string, fields map[string]any) { Log("error", message, fields) }
