module github.com/LuminarysAI/sdk-go

go 1.25

retract (
      v0.1.0      // Deprecated: use v0.2.0
      v0.2.0-rc2  // Pre-release: use v0.2.0
  )

require github.com/vmihailenco/msgpack/v5 v5.4.1

require github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
