package sdk

import (
	"github.com/vmihailenco/msgpack/v5"
)

// HistoryGet fetches dialogue history matching the given filter.
func HistoryGet(filter string) ([]HistoryMessage, error) {
	req := map[string]string{"filter": filter}
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostHistoryGet(ptr, ln))
	var msgs []HistoryMessage
	if err := msgpack.Unmarshal(raw, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

// PromptComplete sends a prompt to the LLM.
func PromptComplete(req PromptRequest) (PromptResponse, error) {
	b := mustMarshal(req)
	ptr, ln := ptrLen(b)
	raw := readHostResult(hostPromptComplete(ptr, ln))
	var resp PromptResponse
	if err := msgpack.Unmarshal(raw, &resp); err != nil {
		return PromptResponse{}, err
	}
	return resp, nil
}

// UnmarshalBatchResult decodes a batch callback payload into a BatchResult.
func UnmarshalBatchResult(payload []byte) (*BatchResult, error) {
	var r BatchResult
	if err := msgpack.Unmarshal(payload, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
