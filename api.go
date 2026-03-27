package sdk

import (
	"unsafe"

	"github.com/vmihailenco/msgpack/v5"
)

func ptrLen(b []byte) (uint32, uint32) {
	if len(b) == 0 {
		return 0, 0
	}
	return uint32(uintptr(unsafe.Pointer(&b[0]))), uint32(len(b))
}

func readHostResult(result uint64) []byte {
	ptr := uint32(result >> 32)
	length := uint32(result & 0xFFFFFFFF)
	if length == 0 {
		return nil
	}
	return unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), length)
}

func mustMarshal(v interface{}) []byte {
	b, err := msgpack.Marshal(v)
	if err != nil {
		panic("sdk: msgpack marshal: " + err.Error())
	}
	return b
}

// ABIError wraps an error string returned from an SDK call.
type ABIError struct {
	Msg string
}

func (e *ABIError) Error() string { return "abi error: " + e.Msg }
