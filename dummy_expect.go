// +build generate_numeric

package hit

type IExpectFloat32 interface{}
type IExpectFloat64 interface{}
type IExpectInt interface{}
type IExpectInt8 interface{}
type IExpectInt16 interface{}
type IExpectInt32 interface{}
type IExpectInt64 interface{}
type IExpectUint interface{}
type IExpectUint8 interface{}
type IExpectUint16 interface{}
type IExpectUint32 interface{}
type IExpectUint64 interface{}

func newExpectFloat32(...interface{}) IExpectFloat32 { return nil }
func newExpectFloat64(...interface{}) IExpectFloat64 { return nil }
func newExpectInt(...interface{}) IExpectInt         { return nil }
func newExpectInt8(...interface{}) IExpectInt8       { return nil }
func newExpectInt16(...interface{}) IExpectInt16     { return nil }
func newExpectInt32(...interface{}) IExpectInt32     { return nil }
func newExpectInt64(...interface{}) IExpectInt64     { return nil }
func newExpectUint(...interface{}) IExpectUint       { return nil }
func newExpectUint8(...interface{}) IExpectUint8     { return nil }
func newExpectUint16(...interface{}) IExpectUint16   { return nil }
func newExpectUint32(...interface{}) IExpectUint32   { return nil }
func newExpectUint64(...interface{}) IExpectUint64   { return nil }
