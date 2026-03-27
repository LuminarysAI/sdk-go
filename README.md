# Luminarys Go SDK

Go SDK for building Luminarys WASM skills.

## Installation

```bash
go mod init my-company/my-skill
# Add to go.mod:
# require github.com/LuminarysAI/sdk-go v0.2.0
```

## Quick Start

```go title="skill.go"
package main

import sdk "github.com/LuminarysAI/sdk-go"

// @skill:id      com.my-company.my-skill
// @skill:name    "My Skill"
// @skill:version 1.0.0
// @skill:desc    "My first skill."

// @skill:method greet "Greet by name."
// @skill:param  name required "User name"
// @skill:result "Greeting text"
func Greet(ctx *sdk.Context, name string) (string, error) {
    return "Hello, " + name + "!", nil
}
```

Build:

```bash
lmsk genkey                            # once: create developer signing key
lmsk generate -lang go .              # generate main.go
GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o my-skill.wasm .
lmsk sign my-skill.wasm               # → com.my-company.my-skill.skill
```

Alternatively, you can use [TinyGo](https://tinygo.org/) for smaller binaries:

```bash
tinygo build -target=wasip1 -o my-skill.wasm .
lmsk sign my-skill.wasm
```

## Documentation

[luminarys.ai](https://luminarys.ai)

## Tools

Download `lmsk` from [releases](https://github.com/LuminarysAI/luminarys/releases).

## License

MIT
