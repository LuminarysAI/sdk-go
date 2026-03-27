package sdk

// NewEmitEvent creates a command that publishes an event to the MessageBus.
func NewEmitEvent(topic string, payload interface{}) Command {
	return Command{
		Type: CmdEmitEvent,
		Payload: map[string]interface{}{
			"topic":   topic,
			"payload": payload,
		},
	}
}

// NewCallModule creates a command that invokes another skill.
func NewCallModule(skillID, method string, payload interface{}, callback, callCtx string) Command {
	return Command{
		Type: CmdCallModule,
		Payload: map[string]interface{}{
			"skill_id": skillID,
			"method":   method,
			"payload":  payload,
			"callback": callback,
			"call_ctx": callCtx,
		},
	}
}

// NewStoreKV creates a command to store a value in the Shared KV (L3).
func NewStoreKV(key string, value interface{}) Command {
	return Command{
		Type: CmdStoreKV,
		Payload: map[string]interface{}{
			"key":   key,
			"value": value,
		},
	}
}

// NewLoadKV creates a command to load a value from the Shared KV (L3).
func NewLoadKV(key, callback string) Command {
	return Command{
		Type: CmdLoadKV,
		Payload: map[string]interface{}{
			"key":      key,
			"callback": callback,
		},
	}
}

// NewSchedule creates a command to schedule a delayed method invocation.
func NewSchedule(method string, delayMs int64, payload interface{}) Command {
	return Command{
		Type: CmdSchedule,
		Payload: map[string]interface{}{
			"method":    method,
			"delay_ms":  delayMs,
			"payload":   payload,
		},
	}
}
