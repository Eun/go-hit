package hit

import (
	"io"
	"net/http"
	"net/url"
)

func boolSliceToInterfaceSlice(params []bool) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func callbackSliceToInterfaceSlice(params []Callback) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func float32SliceToInterfaceSlice(params []float32) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func float64SliceToInterfaceSlice(params []float64) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func intSliceToInterfaceSlice(params []int) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func int8SliceToInterfaceSlice(params []int8) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func int16SliceToInterfaceSlice(params []int16) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func int32SliceToInterfaceSlice(params []int32) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func int64SliceToInterfaceSlice(params []int64) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func uintSliceToInterfaceSlice(params []uint) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func uint8SliceToInterfaceSlice(params []uint8) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func uint16SliceToInterfaceSlice(params []uint16) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func uint32SliceToInterfaceSlice(params []uint32) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func uint64SliceToInterfaceSlice(params []uint64) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func uint8SliceSliceToInterfaceSlice(params [][]uint8) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func readerSliceToInterfaceSlice(params []io.Reader) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

func stringSliceToInterfaceSlice(params []string) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

//nolint:deadcode,unused // keep this in for completion
func uRLSliceToInterfaceSlice(params []*url.URL) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

//nolint:deadcode,unused // keep this in for completion
func requestSliceToInterfaceSlice(params []*http.Request) []interface{} {
	iface := make([]interface{}, len(params))
	for i, v := range params {
		iface[i] = v
	}
	return iface
}

//nolint:deadcode,unused // keep this in for completion
func getLastInterfaceArgument(params []interface{}) (interface{}, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return nil, false
}

func getLastStringArgument(params []string) (string, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return "", false
}

//nolint:deadcode,unused // keep this in for completion
func getLastIntArgument(params []int) (int, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return 0, false
}

//nolint:deadcode,unused // keep this in for completion
func getLastUint8SliceArgument(params [][]uint8) ([]uint8, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return nil, false
}

//nolint:deadcode,unused // keep this in for completion
func getLastByteSliceArgument(params [][]byte) ([]byte, bool) {
	if i := len(params); i > 0 {
		return params[i-1], true
	}
	return nil, false
}
