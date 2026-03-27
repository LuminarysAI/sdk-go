package sdk

import (
	"unsafe"

	"github.com/vmihailenco/msgpack/v5"
)

// SDK version packed as uint32.
const (
	SDKVersionMajor = 0
	SDKVersionMinor = 2
	SDKVersionPatch = 0
	SDKVersion      = (SDKVersionMajor << 16) | (SDKVersionMinor << 8) | SDKVersionPatch
)

// Handler is the function signature every skill must implement.
type Handler func(req InvokeRequest) InvokeResponse

var registeredHandler Handler
var registeredMethods []MethodInfo

// extraHandlers holds handlers registered via RegisterMethods(...).Handle(...).
// Checked by CallHandler before falling through to the generated _skillDispatch.
var extraHandlers = map[string]Handler{}

// Register sets the skill's dispatch handler and appends method descriptors.
// Safe to call multiple times from different init() functions — each call
// appends its methods and, if h is non-nil, sets the handler.
//
// The generated skill_gen.go calls this once with _skillDispatch and the
// generated method list. You do not need to call it again unless you need
// to replace the top-level handler.
func Register(h Handler, methods ...MethodInfo) {
	if h != nil {
		registeredHandler = h
	}
	for _, m := range methods {
		registeredMethods = append(registeredMethods, m)
		if m.handlerFn != nil {
			extraHandlers[m.Name] = m.handlerFn
		}
	}
}

// RegisterMethods appends method descriptors and, if a method carries a
// Handle(...) handler, registers it for dispatch — without touching the
// generated _skillDispatch or _skillHandlers.
//
// This is the primary extension point for hand-written methods:
//
//	func init() {
//	    sdk.RegisterMethods(
//	        sdk.Method("debug_dump", "Dump internal state.").
//	            SetInternal().
//	            Handle(func(req sdk.InvokeRequest) sdk.InvokeResponse {
//	                // ...
//	            }),
//	    )
//	}
func RegisterMethods(methods ...MethodInfo) {
	Register(nil, methods...)
}

// MethodInfo describes a single exported method including its parameter schema.
type MethodInfo struct {
	Name        string      `msgpack:"name"`
	Description string      `msgpack:"description"`
	Params      []ParamInfo `msgpack:"params,omitempty"`
	// MCPHidden marks the method as hidden from MCP tools/list.
	// Hidden methods remain callable by other skills via call_module.
	// Set via @skill:internal annotation (codegen) or .SetInternal() builder.
	MCPHidden bool `msgpack:"mcp_hidden,omitempty"`
	// PrivateCallback marks the method as a private callback.
	// Private callbacks are ONLY callable by the skill itself (host enforces this).
	// Use for TCP/HTTP/WS read callbacks — they should never be invoked by
	// other skills. Set via @skill:callback annotation or .Callback() builder.
	PrivateCallback bool `msgpack:"private_callback,omitempty"`

	// handlerFn is the inline handler set via Handle().
	// Not serialised — stays inside the WASM module.
	handlerFn Handler `msgpack:"-"`
}

// Method starts building a MethodInfo with a fluent builder.
func Method(name, description string) MethodInfo {
	return MethodInfo{Name: name, Description: description}
}

// Param adds a parameter descriptor to the method.
func (m MethodInfo) Param(name string, p ParamInfo) MethodInfo {
	p.Name = name
	m.Params = append(m.Params, p)
	return m
}

// SetInternal marks the method as hidden from MCP tools/list.
// The method remains callable by other skills via call_module.
func (m MethodInfo) SetInternal() MethodInfo {
	m.MCPHidden = true
	return m
}

// Callback marks the method as a private callback.
// Private callbacks are hidden from MCP AND are only callable by the
// skill that owns them — other skills cannot invoke them via call_module.
//
// Use for all TCP/HTTP/WebSocket read callbacks:
//
//	sdk.Method("on_data", "Receives TCP data").Callback()
func (m MethodInfo) Callback() MethodInfo {
	m.MCPHidden = true
	m.PrivateCallback = true
	return m
}

// Handle attaches an inline handler function to this method descriptor.
// When the method is registered via RegisterMethods, the SDK dispatches
// incoming calls to this function — no changes to the generated files needed.
//
//	sdk.RegisterMethods(
//	    sdk.Method("debug_dump", "Dump internal state.").
//	        SetInternal().
//	        Handle(func(req sdk.InvokeRequest) sdk.InvokeResponse {
//	            payload, _ := msgpack.Marshal(map[string]int{"conns": len(activeConns)})
//	            return sdk.InvokeResponse{Payload: payload}
//	        }),
//	)
func (m MethodInfo) Handle(fn Handler) MethodInfo {
	m.handlerFn = fn
	return m
}

func TypeString() ParamInfo { return ParamInfo{TypeVal: "string"} }
func TypeInt() ParamInfo    { return ParamInfo{TypeVal: "integer"} }
func TypeNumber() ParamInfo { return ParamInfo{TypeVal: "number"} }
func TypeBool() ParamInfo   { return ParamInfo{TypeVal: "boolean"} }
func TypeObject() ParamInfo { return ParamInfo{TypeVal: "object"} }

// TypeArray creates an array parameter with the given item type:
//
//	sdk.TypeArray(sdk.TypeString())  →  array of strings
func TypeArray(items ParamInfo) ParamInfo {
	return ParamInfo{TypeVal: "array", Items: &items}
}

// ParamInfo describes one parameter of a method.
type ParamInfo struct {
	Name        string      `msgpack:"name"`
	Description string      `msgpack:"description,omitempty"`
	TypeVal     string      `msgpack:"type"`
	IsRequired  bool        `msgpack:"required,omitempty"`
	EnumVals    []string    `msgpack:"enum,omitempty"`
	Items       *ParamInfo  `msgpack:"items,omitempty"`
	DefaultVal  interface{} `msgpack:"default,omitempty"`
	ExampleVal  interface{} `msgpack:"example,omitempty"`
}

// Required marks the parameter as required.
func (p ParamInfo) Required() ParamInfo { p.IsRequired = true; return p }

// Desc sets the human-readable description.
func (p ParamInfo) Desc(d string) ParamInfo { p.Description = d; return p }

// Enum restricts accepted values to the given set.
func (p ParamInfo) Enum(values ...string) ParamInfo { p.EnumVals = values; return p }

// Default sets the default value shown in documentation.
func (p ParamInfo) Default(v interface{}) ParamInfo { p.DefaultVal = v; return p }

// Example sets an illustrative example value.
func (p ParamInfo) Example(v interface{}) ParamInfo { p.ExampleVal = v; return p }

// SkillDescriptor describes the skill's identity, methods, and requirements.
type SkillDescriptor struct {
	Version      int               `msgpack:"version"`
	SkillID      string            `msgpack:"skill_id,omitempty"`
	SkillName    string            `msgpack:"skill_name,omitempty"`
	SkillVersion string            `msgpack:"skill_version,omitempty"`
	Description  string            `msgpack:"description,omitempty"`
	Methods      []MethodInfo      `msgpack:"methods"`
	Requirements []RequirementInfo `msgpack:"requirements,omitempty"`
	// SDKVersionCode is the SDK version packed as uint32.
	SDKVersionCode uint32 `msgpack:"sdk_version,omitempty"`
}

// RequirementInfo declares one permission the skill expects from its manifest.
// Populated by codegen from @skill:require annotations.
type RequirementInfo struct {
	Kind    string `msgpack:"kind"`              // "fs", "http", "tcp", "ws", "shell", "llm", "env", "invoke"
	Pattern string `msgpack:"pattern,omitempty"` // glob, URL, command, env name, skill ID
	Mode    string `msgpack:"mode,omitempty"`    // "ro" or "rw" (fs only)
}

var registeredRequirements []RequirementInfo

// RegisterRequirements appends requirement descriptors for the skill.
// Called by generated init() code.
func RegisterRequirements(reqs ...RequirementInfo) {
	registeredRequirements = append(registeredRequirements, reqs...)
}

// SkillIdentity holds the skill's intrinsic metadata from annotations.
var skillIdentity struct {
	ID          string
	Name        string
	Version     string
	Description string
}

// SetSkillIdentity sets the skill's identity metadata.
// Called by generated init() code from @skill:id, @skill:name, etc.
func SetSkillIdentity(id, name, version, description string) {
	skillIdentity.ID = id
	skillIdentity.Name = name
	skillIdentity.Version = version
	skillIdentity.Description = description
}

// GetDescriptor returns the SkillDescriptor for the registered skill.
func GetDescriptor() SkillDescriptor {
	return SkillDescriptor{
		Version:        Version,
		SkillID:        skillIdentity.ID,
		SkillName:      skillIdentity.Name,
		SkillVersion:   skillIdentity.Version,
		Description:    skillIdentity.Description,
		Methods:        registeredMethods,
		Requirements:   registeredRequirements,
		SDKVersionCode: SDKVersion,
	}
}

// CallHandler dispatches an InvokeRequest to the registered handler.
// Extra handlers registered via RegisterMethods(...).Handle(...) take
// priority over the generated _skillDispatch.
func CallHandler(req InvokeRequest) InvokeResponse {
	// Extra handlers registered via RegisterMethods take priority.
	if h, ok := extraHandlers[req.Method]; ok {
		return h(req)
	}
	if registeredHandler == nil {
		return InvokeResponse{
			Version: Version,
			Error:   "no handler registered — call sdk.Register() in init()",
		}
	}
	return registeredHandler(req)
}

// WriteResult stores b in a GC-safe buffer and returns (ptr<<32 | len).
func WriteResult(b []byte) uint64 {
	keepAlive = b
	if len(b) == 0 {
		return 0
	}
	ptr := uint32(uintptr(unsafe.Pointer(&b[0])))
	return (uint64(ptr) << 32) | uint64(len(b))
}

var keepAlive []byte

// ShrinkResultBuf releases the persistent result buffer.
// Called by the generated skill_free when the host signals it has read the result.
func ShrinkResultBuf(ptr uint32) {
	if keepAlive == nil {
		return
	}
	bufPtr := uint32(uintptr(unsafe.Pointer(&keepAlive[0])))
	if ptr == bufPtr {
		keepAlive = nil
	}
}

// MarshalErrorResponse encodes an error-only InvokeResponse.
func MarshalErrorResponse(msg string) []byte {
	resp := InvokeResponse{Version: Version, Error: msg}
	b, _ := msgpack.Marshal(resp)
	return b
}
