package sdk

// Context is passed to every skill handler. Provides read access to invocation
// metadata and write access to response annotations.
type Context struct {
	req *InvokeRequest
	res *InvokeResponse
}

// NewContext creates a Context bound to the given request and response.
func NewContext(req *InvokeRequest, res *InvokeResponse) *Context {
	return &Context{req: req, res: res}
}

// RequestID returns the unique ID of this invocation.
func (c *Context) RequestID() string { return c.req.RequestID }

// TraceID returns the distributed trace ID.
func (c *Context) TraceID() string { return c.req.TraceID }

// SessionID returns the user session ID.
func (c *Context) SessionID() string { return c.req.SessionID }

// LLMSessionID returns the LLM session ID.
func (c *Context) LLMSessionID() string { return c.req.LLMSessionID }

// SkillID returns the ID of the skill being invoked.
func (c *Context) SkillID() string { return c.req.SkillID }

// CallerID returns the skill ID of the caller when invoked via call_module.
func (c *Context) CallerID() string { return c.req.CallerID }

// Method returns the method name being invoked.
func (c *Context) Method() string { return c.req.Method }

// SetLLMContext sets additional text appended to the tool result for the LLM.
func (c *Context) SetLLMContext(text string) {
	c.res.LLMContext = text
}

// AppendLLMContext appends text to any existing LLMContext, separated by "\n".
func (c *Context) AppendLLMContext(text string) {
	if c.res.LLMContext == "" {
		c.res.LLMContext = text
	} else {
		c.res.LLMContext += "\n" + text
	}
}
