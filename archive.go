package sdk

import (
	"github.com/vmihailenco/msgpack/v5"
)

// ArchivePackResult is the response from ArchivePack.
type ArchivePackResult struct {
	FilesCount int    `msgpack:"files_count"`
	Format     string `msgpack:"format"`
	Error      string `msgpack:"error,omitempty"`
}

// ArchivePack creates a tar.gz or zip archive from a directory.
func ArchivePack(source, output, format, exclude string) (ArchivePackResult, error) {
	req := map[string]string{
		"source": source, "output": output, "format": format,
		"exclude": exclude,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostArchivePack(ptr, ln))
	var resp ArchivePackResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return ArchivePackResult{}, err
	}
	if resp.Error != "" {
		return ArchivePackResult{}, &ABIError{resp.Error}
	}
	return resp, nil
}

// ArchiveUnpackResult is the response from ArchiveUnpack.
type ArchiveUnpackResult struct {
	FilesCount int    `msgpack:"files_count"`
	Error      string `msgpack:"error,omitempty"`
}

// ArchiveUnpack extracts a tar.gz or zip archive to a directory.
func ArchiveUnpack(archive, dest, format, exclude string, strip int) (ArchiveUnpackResult, error) {
	req := map[string]interface{}{
		"archive": archive, "dest": dest, "format": format,
		"exclude": exclude, "strip": strip,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostArchiveUnpack(ptr, ln))
	var resp ArchiveUnpackResult
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return ArchiveUnpackResult{}, err
	}
	if resp.Error != "" {
		return ArchiveUnpackResult{}, &ABIError{resp.Error}
	}
	return resp, nil
}

// ArchiveListEntry is one entry in an archive listing.
type ArchiveListEntry struct {
	Name  string `msgpack:"name"`
	Size  int64  `msgpack:"size"`
	IsDir bool   `msgpack:"is_dir"`
}

// ArchiveList lists the contents of a tar.gz or zip archive.
func ArchiveList(archivePath, format, exclude string) ([]ArchiveListEntry, error) {
	req := map[string]string{
		"archive": archivePath, "format": format,
		"exclude": exclude,
	}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostArchiveList(ptr, ln))
	var resp struct {
		Entries []ArchiveListEntry `msgpack:"entries"`
		Error   string             `msgpack:"error,omitempty"`
	}
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return nil, err
	}
	if resp.Error != "" {
		return nil, &ABIError{resp.Error}
	}
	return resp.Entries, nil
}
