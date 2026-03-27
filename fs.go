package sdk

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

type FSRequest struct {
	Path    string `msgpack:"path"`
	Content []byte `msgpack:"content,omitempty"`
}

type FSMkdirRequest struct {
	Path string `msgpack:"path"`
}

type FSResponse struct {
	Content []byte `msgpack:"content,omitempty"`
	Error   string `msgpack:"error,omitempty"`
}

type FSLsRequest struct {
	Path string `msgpack:"path"`
	Long bool   `msgpack:"long"`
}

// DirEntry represents one entry returned by FsLs.
type DirEntry struct {
	// Name is the file or directory name.
	Name string `msgpack:"name"`
	// Size in bytes (0 for directories).
	Size int64 `msgpack:"size"`
	// IsDir is true for directories.
	IsDir bool `msgpack:"is_dir"`
	// ModTime is the Unix timestamp in seconds. Only populated when long=true.
	ModTime int64 `msgpack:"mod_time,omitempty"`
	// Mode is the permission bits (e.g. 493 = 0o755). Only populated when long=true.
	Mode uint32 `msgpack:"mode,omitempty"`
	// ModeStr is human-readable permissions, e.g. "drwxr-xr-x". Only populated when long=true.
	ModeStr string `msgpack:"mode_str,omitempty"`
}

type FSLsResponse struct {
	Entries []DirEntry `msgpack:"entries"`
	Error   string     `msgpack:"error,omitempty"`
}

type FSChmodRequest struct {
	Path      string `msgpack:"path"`
	Mode      uint32 `msgpack:"mode"`
	Recursive bool   `msgpack:"recursive"`
}

// FSReadLinesRequest configures FsReadLines.
type FSReadLinesRequest struct {
	// Path is the file path. Required.
	Path string `msgpack:"path"`
	// Offset is the 0-based start line.
	Offset int `msgpack:"offset"`
	// Limit is the max lines to return (0 = all).
	Limit int `msgpack:"limit"`
}

// TextFileContent is the result of FsReadLines.
type TextFileContent struct {
	Lines       []string `msgpack:"lines"`
	TotalLines  int      `msgpack:"total_lines"`
	Offset      int      `msgpack:"offset"`
	IsTruncated bool     `msgpack:"is_truncated"`
}

type textFileContentResponse struct {
	Lines       []string `msgpack:"lines"`
	TotalLines  int      `msgpack:"total_lines"`
	Offset      int      `msgpack:"offset"`
	IsTruncated bool     `msgpack:"is_truncated"`
	Error       string   `msgpack:"error,omitempty"`
}

// GrepOptions configures FsGrep.
type GrepOptions struct {
	// Pattern is a regex (or literal when Fixed=true).
	Pattern string `msgpack:"pattern"`
	// Path to search. Empty = all allowed directories.
	Path            string   `msgpack:"path"`
	Fixed           bool     `msgpack:"fixed"`
	CaseInsensitive bool     `msgpack:"case_insensitive"`
	MaxDepth        int      `msgpack:"max_depth"`
	MaxCount        int      `msgpack:"max_count"`
	Workers         int      `msgpack:"workers"`
	TypeFilter      string   `msgpack:"type_filter"`
	Include         []string `msgpack:"include"`
	Exclude         []string `msgpack:"exclude"`
	IgnoreDirs      []string `msgpack:"ignore_dirs"`
	WithLines       bool     `msgpack:"with_lines"`
	FilenameOnly    bool     `msgpack:"filename_only"`
}

// GrepFileMatch holds all matches in one file.
type GrepFileMatch struct {
	Path    string          `msgpack:"path"`
	Matches []GrepLineMatch `msgpack:"matches"`
}

// GrepLineMatch is one matching line.
type GrepLineMatch struct {
	LineNum int         `msgpack:"line_num"`
	Line    string      `msgpack:"line,omitempty"`
	Ranges  []GrepRange `msgpack:"ranges,omitempty"`
}

// GrepRange is a [Start, End) byte offset within a line.
type GrepRange struct {
	Start int `msgpack:"start"`
	End   int `msgpack:"end"`
}

type grepResponse struct {
	Matches []GrepFileMatch `msgpack:"matches"`
	Error   string          `msgpack:"error,omitempty"`
}

// GlobOptions configures FsGlob.
type GlobOptions struct {
	// Patterns is a list of glob patterns (union). Supports *, **, ?, [abc], {a,b}.
	Patterns []string `msgpack:"patterns"`
	// Path is the base directory. Empty = all allowed directories.
	Path        string   `msgpack:"path"`
	MatchHidden bool     `msgpack:"match_hidden"`
	IgnoreDirs  []string `msgpack:"ignore_dirs"`
	MaxDepth    int      `msgpack:"max_depth"`
	OnlyFiles   bool     `msgpack:"only_files"`
	OnlyDirs    bool     `msgpack:"only_dirs"`
}

// GlobEntry is one result from FsGlob.
type GlobEntry struct {
	Path  string `msgpack:"path"`
	IsDir bool   `msgpack:"is_dir"`
}

type globEntryResponse struct {
	Matches []GlobEntry `msgpack:"matches"`
	Error   string      `msgpack:"error,omitempty"`
}

// AllowedDir describes one directory the skill can access.
type AllowedDir struct {
	// Path is the absolute directory path.
	Path string `msgpack:"path"`
	// Mode is "ro" (read-only) or "rw" (read-write).
	Mode string `msgpack:"mode"`
}

// Permission bit constants.
const (
	PermOwnerRead  uint32 = 0400
	PermOwnerWrite uint32 = 0200
	PermOwnerExec  uint32 = 0100
	PermGroupRead  uint32 = 0040
	PermGroupWrite uint32 = 0020
	PermGroupExec  uint32 = 0010
	PermOtherRead  uint32 = 0004
	PermOtherWrite uint32 = 0002
	PermOtherExec  uint32 = 0001

	PermReadOnly   uint32 = 0444 // r--r--r--
	PermDefault    uint32 = 0644 // rw-r--r--
	PermDefaultDir uint32 = 0755 // rwxr-xr-x
	PermPrivate    uint32 = 0600 // rw-------
	PermPrivateDir uint32 = 0700 // rwx------
)

// FsRead reads a file and returns its contents as bytes. Requires fs.enabled.
func FsRead(path string) ([]byte, error) {
	req := FSRequest{Path: path}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsRead(ptr, ln))
	var resp FSResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, &ABIError{resp.Error}
	}
	return resp.Content, nil
}

// FsWrite writes content to a file, replacing it entirely. Requires fs.enabled.
func FsWrite(path string, content []byte) error {
	req := FSRequest{Path: path, Content: content}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsWrite(ptr, ln))
	var resp FSResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}

// FsCreate creates a new file. Fails if the file already exists. Requires fs.enabled.
func FsCreate(path string, content []byte) error {
	req := FSRequest{Path: path, Content: content}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsCreate(ptr, ln))
	var resp FSResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}

// FsDelete deletes a file or directory. Requires fs.enabled.
func FsDelete(path string) error {
	req := FSRequest{Path: path}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsDelete(ptr, ln))
	var resp FSResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}

// FsMkdir creates a directory (with parents, like mkdir -p). Requires fs.enabled.
//
// Supports brace expansion: "logs/{2024,2025}" creates two directories.
func FsMkdir(path string) error {
	req := FSMkdirRequest{Path: path}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsMkdir(ptr, ln))
	var resp FSResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}

// FsLs lists directory contents. Requires fs.enabled.
//
// When long is true, ModTime, Mode, and ModeStr are populated for each entry.
func FsLs(path string, long bool) ([]DirEntry, error) {
	req := FSLsRequest{Path: path, Long: long}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsLs(ptr, ln))
	var resp FSLsResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, &ABIError{resp.Error}
	}
	return resp.Entries, nil
}

// FsChmod changes file permissions. Requires fs.enabled.
//
// mode is the permission bits (e.g. 0o755 = 493). Use PermXxx constants.
func FsChmod(path string, mode uint32, recursive bool) error {
	req := FSChmodRequest{Path: path, Mode: mode, Recursive: recursive}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsChmod(ptr, ln))
	var resp FSResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return err
	}
	if resp.Error != "" {
		return &ABIError{resp.Error}
	}
	return nil
}

// FsReadLines reads a range of lines from a text file. Requires fs.enabled.
//
// Offset is 0-based. Limit 0 means all remaining lines.
func FsReadLines(req FSReadLinesRequest) (TextFileContent, error) {
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsReadLines(ptr, ln))
	var resp textFileContentResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return TextFileContent{}, err
	}
	if resp.Error != "" {
		return TextFileContent{}, errors.New(resp.Error)
	}
	return TextFileContent{
		Lines:       resp.Lines,
		TotalLines:  resp.TotalLines,
		Offset:      resp.Offset,
		IsTruncated: resp.IsTruncated,
	}, nil
}

// FsGrep searches file contents by regex pattern. Requires fs.enabled.
func FsGrep(opts GrepOptions) ([]GrepFileMatch, error) {
	b := mustMarshal(opts)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsGrep(ptr, ln))
	var resp grepResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, &ABIError{resp.Error}
	}
	return resp.Matches, nil
}

// FsGlob finds files matching glob patterns. Requires fs.enabled.
func FsGlob(opts GlobOptions) ([]GlobEntry, error) {
	b := mustMarshal(opts)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsGlob(ptr, ln))
	var resp globEntryResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, &ABIError{resp.Error}
	}
	return resp.Matches, nil
}

// FsAllowedDirs returns the directories this skill can access. Requires fs.enabled.
func FsAllowedDirs() ([]AllowedDir, error) {
	b := mustMarshal(map[string]any{})
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsAllowedDirs(ptr, ln))
	var resp struct {
		Dirs  []AllowedDir `msgpack:"dirs"`
		Error string       `msgpack:"error,omitempty"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, &ABIError{resp.Error}
	}
	return resp.Dirs, nil
}

// FsCopy copies a file. Both paths must be within allowed directories. Requires fs.enabled.
func FsCopy(source, dest string) error {
	req := map[string]string{"source": source, "dest": dest}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostFsCopy(ptr, ln))
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
