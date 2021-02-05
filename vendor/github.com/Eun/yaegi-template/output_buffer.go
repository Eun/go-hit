package yaegi_template

import (
	"bytes"

	"go.uber.org/atomic"
)

type outputBuffer struct {
	buf           *bytes.Buffer
	discardWrites *atomic.Bool
	size          uint64
}

func newOutputBuffer(discardWrites bool) *outputBuffer {
	return &outputBuffer{
		buf:           bytes.NewBuffer(nil),
		discardWrites: atomic.NewBool(discardWrites),
		size:          0,
	}
}

func (ob *outputBuffer) Write(p []byte) (int, error) {
	if ob.discardWrites.Load() {
		return len(p), nil
	}
	n, err := ob.buf.Write(p)
	if n > 0 {
		ob.size += uint64(n)
	}
	return n, err
}

func (ob *outputBuffer) Reset() {
	ob.buf.Reset()
	ob.size = 0
}

func (ob *outputBuffer) Bytes() []byte {
	return ob.buf.Bytes()
}

func (ob *outputBuffer) DiscardWrites(v bool) {
	ob.discardWrites.Store(v)
}

func (ob *outputBuffer) Length() uint64 {
	return ob.size
}
