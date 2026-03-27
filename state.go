package sdk

import "github.com/vmihailenco/msgpack/v5"

// StateEnvelope wraps any skill state with a schema version.
type StateEnvelope struct {
	SchemaVersion int    `msgpack:"schema_version"`
	Data          []byte `msgpack:"data"`
}

// MarshalState encodes skill state into a versioned envelope.
func MarshalState(schemaVersion int, state interface{}) ([]byte, error) {
	data, err := msgpack.Marshal(state)
	if err != nil {
		return nil, err
	}
	env := StateEnvelope{SchemaVersion: schemaVersion, Data: data}
	return msgpack.Marshal(env)
}

// UnmarshalState decodes the versioned envelope and populates dst.
// Returns the schema version so callers can run migrations if needed.
func UnmarshalState(raw []byte, dst interface{}) (schemaVersion int, err error) {
	if len(raw) == 0 {
		return 0, nil
	}
	var env StateEnvelope
	if err = msgpack.Unmarshal(raw, &env); err != nil {
		return 0, err
	}
	if err = msgpack.Unmarshal(env.Data, dst); err != nil {
		return 0, err
	}
	return env.SchemaVersion, nil
}
